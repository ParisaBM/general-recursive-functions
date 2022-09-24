package compile

import (
	"fmt"
	"os"

	"github.com/llir/irutil"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// funcList maps function ids to their function objects in the module
var funcList []*ir.Func

// we'll use this as a more consice way to create constants
// n will typically be either 0 or 1
func CI(n byte) *constant.Int {
	return constant.NewInt(types.I32, int64(n))
}

// code gen uses the llir package to add all the instructions of a program to a module,
// then write the output to outputFile
func codeGen(outputFile *os.File) {
	m := ir.NewModule()
	// since input is from command line, and output is to stdio
	// we use some library functions to handle these
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	atoi := m.NewFunc("atoi", types.I32)
	atoi.Sig.Variadic = true
	formatString := m.NewGlobalDef("fstr", irutil.NewCString("%d\n"))
	funcList = make([]*ir.Func, 0)
	for {
		// tToC should consist of identifier with their function definitions then an end token
		// anything else means something's gone wrong
		switch tToC.get() {
		case definition:
			id := tToC.get()
			nameListMutex.Lock()
			functionName := nameList[id]
			nameListMutex.Unlock()
			paramCount := tToC.get()
			paramRegisters := make([]value.Value, 0)
			var f *ir.Func
			if functionName == "main" {
				// this is the overhead we add to main to get and parse the command line arguements
				argc := ir.NewParam("argc", types.I32)
				argv := ir.NewParam("argv", types.NewPointer(types.NewPointer(types.I8)))
				f = m.NewFunc("main", types.I32, argc, argv)
				b := f.NewBlock("")
				for i := byte(0); i < paramCount; i += 1 {
					paramStringPointer := b.NewGetElementPtr(types.NewPointer(types.I8), argv, CI(i+1))
					paramString := b.NewLoad(types.NewPointer(types.I8), paramStringPointer)
					paramRegisters = append(paramRegisters, b.NewCall(atoi, paramString))
				}
				ret := represent(f, paramRegisters)
				// unlike other functions that return their results, main prints it and returns 0
				b = f.Blocks[len(f.Blocks)-1]
				formatStringPointer := b.NewGetElementPtr(types.NewArray(4, types.I8), formatString, CI(0), CI(0))
				b.NewCall(printf, formatStringPointer, ret)
				b.NewRet(CI(0))
			} else {
				params := make([]*ir.Param, 0)
				for i := byte(0); i < paramCount; i += 1 {
					params = append(params, ir.NewParam("p"+fmt.Sprint(i), types.I32))
				}
				f = m.NewFunc(functionName, types.I32, params...)
				// the params variable uses a list of ir.Param values
				// this type is only used for function parameter values
				// however every value saved in a register implements the value.Value interface
				// thus we convert to this interface to obtain values we can treat generically as registers
				paramRegisters := make([]value.Value, 0)
				for _, r := range params {
					paramRegisters = append(paramRegisters, r)
				}
				f.NewBlock("")
				ret := represent(f, paramRegisters)
				b := f.Blocks[len(f.Blocks)-1]
				b.NewRet(ret)
			}
			funcList = append(funcList, f)
		case end:
			outputFile.WriteString(fmt.Sprint(m))
			return
		default:
			panic("invalid token in tToC")
		}
	}
}

// represent recursively generates the llvm ir for a single function definition
// these instructions all get added to the first parameter f
// what gets returned is either a register or constant with the result of that function
// all instructions get added starting from the last block in f i.e. the active block
func represent(f *ir.Func, params []value.Value) value.Value {
	b := f.Blocks[len(f.Blocks)-1] // this finds the active block
	switch tToC.get() {
	case constT:
		return CI(tToC.get())
	case proj:
		return params[tToC.get()]
	case suc:
		newValue := b.NewAdd(params[0], CI(1))
		return newValue
	case comp:
		// suppose we have a composition like f(g(x), h(x))
		// g and h will appear first in the stream
		// the values of g(x) and h(x) we call our top level parameters
		// once they've been computed we pass them along to f
		topLevelParams := make([]value.Value, 0)
		numTopLevel := tToC.get()
		for i := byte(0); i < numTopLevel; i += 1 {
			topLevelParams = append(topLevelParams, represent(f, params))
		}
		return represent(f, topLevelParams)
	case rec:
		// to handle primitive recursion we create 3 blocks: looValidate, loopBody, and loopExit
		// first we initialize currentValue, and counter to the base case function and 0 respectively
		// then we can branch to loopValidate
		// loopValidate checks the counter to see if enough iterations of the loop have been perfomed already
		// then it branches to either loopBody or loopExit
		// loopBody updates currentValue and increments counter
		// the it branches back to loopValidate

		// our blocks
		loopValidate := f.NewBlock("")
		loopBody := f.NewBlock("")
		// loop exit gets added at the end so loop body will be the active block

		// setup code i.e. b code
		initialValue := represent(f, params[1:])
		b.NewBr(loopValidate)

		// loopValidate code
		currentValue := loopValidate.NewPhi(ir.NewIncoming(initialValue, b))
		counter := loopValidate.NewPhi(ir.NewIncoming(CI(0), b))
		stopCondition := loopValidate.NewICmp(enum.IPredULT, counter, params[0])

		// loopBody code
		nextValue := represent(f, append([]value.Value{counter, currentValue}, params[1:]...))
		nextCounter := loopBody.NewAdd(counter, CI(1))
		loopBody.NewBr(loopValidate)

		// here we can do those phi instruction from earlier
		currentValue.Incs = append(currentValue.Incs, ir.NewIncoming(nextValue, loopBody))
		counter.Incs = append(counter.Incs, ir.NewIncoming(nextCounter, loopBody))

		// loopExit is where we continue from
		loopExit := f.NewBlock("")
		loopValidate.NewCondBr(stopCondition, loopBody, loopExit)

		return currentValue
	case min:
		// this case is extremely similar to rec
		// here we increment counter until result is 0

		loopBody := f.NewBlock("")
		b.NewBr(loopBody)

		// loopBody code
		counter := loopBody.NewPhi(ir.NewIncoming(CI(0), b))
		result := represent(f, append([]value.Value{counter}, params...))
		nextCounter := loopBody.NewAdd(counter, CI(1))
		exitCondition := loopBody.NewICmp(enum.IPredEQ, result, CI(0))
		loopExit := f.NewBlock("")
		loopBody.NewCondBr(exitCondition, loopExit, loopBody)

		counter.Incs = append(counter.Incs, ir.NewIncoming(nextCounter, loopBody))

		return counter
	case identifier:
		id := tToC.get()
		return b.NewCall(funcList[id], params...)
	default:
		panic("invalid token in tToC")
	}
}
