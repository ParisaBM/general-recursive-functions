package main

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

var funcList []*ir.Func

func codeGen(outputFileName string) {
	m := ir.NewModule()
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	atoi := m.NewFunc("atoi", types.I32)
	atoi.Sig.Variadic = true
	formatString := m.NewGlobalDef("fstr", irutil.NewCString("%d\n"))
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
				argc := ir.NewParam("argc", types.I32)
				argv := ir.NewParam("argv", types.NewPointer(types.NewPointer(types.I8)))
				f = m.NewFunc("main", types.Void, argc, argv)
				entry := f.NewBlock("")
				for i := 0; i < int(paramCount); i += 1 {
					paramStringPointer := entry.NewGetElementPtr(types.NewPointer(types.I8), argv, constant.NewInt(types.I32, int64(i+1)))
					paramString := entry.NewLoad(types.NewPointer(types.I8), paramStringPointer)
					paramRegisters = append(paramRegisters, entry.NewCall(atoi, paramString))
				}
				ret := represent(f, &entry, paramRegisters)
				formatStringPointer := entry.NewGetElementPtr(types.NewArray(4, types.I8), formatString, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				entry.NewCall(printf, formatStringPointer, ret)
				entry.NewRet(nil)
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
				entry := f.NewBlock("")
				ret := represent(f, &entry, paramRegisters)
				entry.NewRet(ret)
			}
			funcList = append(funcList, f)
		case end:
			file, err := os.Create(outputFileName)
			if err != nil {
				println(err)
			}
			file.WriteString(fmt.Sprint(m))
			cToE.put(0)
			return
		default:
			panic("invalid token in tToC")
		}
	}
}

// represent recursively generates the llvm ir for a single function definition
// these instructions all get added to the first parameter f
// b is the current active block
// b is doubly indirect so that instances of represent can update it, and when they do other instances further down in
// the stack can add instructions to the block where it was left off
// params are the parameters to that sub function
// what gets returned is either a register or constant with the result of that function
func represent(f *ir.Func, b **ir.Block, params []value.Value) value.Value {
	switch tToC.get() {
	case constT:
		return constant.NewInt(types.I32, int64(tToC.get()))
	case proj:
		return params[tToC.get()]
	case suc:
		newValue := (*b).NewAdd(params[0], constant.NewInt(types.I32, 1))
		return newValue
	case comp:
		// suppose we have a composition like f(g(x), h(x))
		// g and h will appear first in the stream
		// the values of g(x) and h(x) we call our top level parameters
		// once they've been computed we pass them along to f
		topLevelParams := make([]value.Value, 0)
		numTopLevel := tToC.get()
		for i := byte(0); i < numTopLevel; i += 1 {
			topLevelParams = append(topLevelParams, represent(f, b, params))
		}
		return represent(f, b, topLevelParams)
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
		loopExit := f.NewBlock("")

		// setup code i.e. b code
		initialValue := represent(f, b, params[1:])
		(*b).NewBr(loopValidate)

		// loopValidate code
		currentValue := loopValidate.NewPhi(ir.NewIncoming(initialValue, *b))
		counter := loopValidate.NewPhi(ir.NewIncoming(constant.NewInt(types.I32, 0), *b))
		stopCondition := loopValidate.NewICmp(enum.IPredULT, counter, params[0])
		loopValidate.NewCondBr(stopCondition, loopBody, loopExit)

		// loopBody code
		nextValue := represent(f, &loopBody, append([]value.Value{counter, currentValue}, params[1:]...))
		nextCounter := loopBody.NewAdd(counter, constant.NewInt(types.I32, 1))
		loopBody.NewBr(loopValidate)

		// here we can do those phi instruction from earlier
		currentValue.Incs = append(currentValue.Incs, ir.NewIncoming(nextValue, loopBody))
		counter.Incs = append(counter.Incs, ir.NewIncoming(nextCounter, loopBody))

		// loopExit is where we continue from
		*b = loopExit

		return currentValue
	case min:
		// this case is extremely similar to rec
		// here we increment counter until result is 0

		// our blocks
		loopBody := f.NewBlock("")
		loopExit := f.NewBlock("")

		(*b).NewBr(loopBody)

		// loopBody code
		counter := loopBody.NewPhi(ir.NewIncoming(constant.NewInt(types.I32, 0), *b))
		result := represent(f, &loopBody, append([]value.Value{counter}, params...))
		nextCounter := loopBody.NewAdd(counter, constant.NewInt(types.I32, 1))
		exitCondition := loopBody.NewICmp(enum.IPredEQ, result, constant.NewInt(types.I32, 0))
		loopBody.NewCondBr(exitCondition, loopExit, loopBody)

		counter.Incs = append(counter.Incs, ir.NewIncoming(nextCounter, loopBody))

		*b = loopExit

		return counter
	case identifier:
		id := tToC.get()
		return (*b).NewCall(funcList[id], params...)
	default:
		panic("invalid token in tToC")
	}
}
