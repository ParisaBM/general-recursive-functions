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

var func_list []*ir.Func

func code_gen(output_file_name string) {
	m := ir.NewModule()
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	atoi := m.NewFunc("atoi", types.I32)
	atoi.Sig.Variadic = true
	format_string := m.NewGlobalDef("fstr", irutil.NewCString("%d\n"))
	for {
		//t_to_c should consist of identifier with their function definitions then an end token
		//anything else means something's gone wrong
		switch t_to_c.get() {
		case identifier:
			id := t_to_c.get()
			name_list_mutex.Lock()
			function_name := name_list[id]
			name_list_mutex.Unlock()
			param_count := t_to_c.get()
			param_registers := make([]value.Value, 0)
			var f *ir.Func
			if function_name == "main" {
				argc := ir.NewParam("argc", types.I32)
				argv := ir.NewParam("argv", types.NewPointer(types.NewPointer(types.I8)))
				f = m.NewFunc("main", types.Void, argc, argv)
				entry := f.NewBlock("")
				for i := 0; i < int(param_count); i += 1 {
					param_string_pointer := entry.NewGetElementPtr(types.NewPointer(types.I8), argv, constant.NewInt(types.I32, int64(i+1)))
					param_string := entry.NewLoad(types.NewPointer(types.I8), param_string_pointer)
					param_registers = append(param_registers, entry.NewCall(atoi, param_string))
				}
				ret := represent(f, &entry, param_registers)
				format_string_pointer := entry.NewGetElementPtr(types.NewArray(4, types.I8), format_string, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				//  := entry.NewGetElementPtr(types.NewPointer(types.I8), format_string, constant.NewInt(types.I32, 0))
				entry.NewCall(printf, format_string_pointer, ret)
				entry.NewRet(nil)
			} else {
				params := make([]*ir.Param, 0)
				for i := byte(0); i < param_count; i += 1 {
					params = append(params, ir.NewParam("p"+fmt.Sprint(i), types.I32))
				}
				f = m.NewFunc(function_name, types.I32, params...)
				//the params variable uses a list of ir.Param values
				//this type is only used for function parameter values
				//however every value saved in a register implements the value.Value interface
				//thus we convert to this interface to obtain values we can treat generically as registers
				param_registers := make([]value.Value, 0)
				for _, r := range params {
					param_registers = append(param_registers, r)
				}
				entry := f.NewBlock("")
				ret := represent(f, &entry, param_registers)
				entry.NewRet(ret)
			}
			func_list = append(func_list, f)
		case end:
			file, err := os.Create(output_file_name)
			if err != nil {
				println(err)
			}
			file.WriteString(fmt.Sprint(m))
			c_to_e.put(0)
			return
		default:
			panic("invalid token in t_to_c")
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
	switch t_to_c.get() {
	case const_t:
		return constant.NewInt(types.I32, int64(t_to_c.get()))
	case proj:
		return params[t_to_c.get()]
	case suc:
		new_value := (*b).NewAdd(params[0], constant.NewInt(types.I32, 1))
		return new_value
	case comp:
		//suppose we have a composition like f(g(x), h(x))
		//g and h will appear first in the stream
		//the values of g(x) and h(x) we call our top level parameters
		//once they've been computed we pass them along to f
		top_level_params := make([]value.Value, 0)
		num_top_level := t_to_c.get()
		for i := byte(0); i < num_top_level; i += 1 {
			top_level_params = append(top_level_params, represent(f, b, params))
		}
		return represent(f, b, top_level_params)
	case rec:
		//to handle primitive recursion we create 3 blocks: loop_validate, loop_body, and loop_exit
		//first we initialize current_value, and counter to the base case function and 0 respectively
		//then we can branch to loop_validate
		//loop_validate checks the counter to see if enough iterations of the loop have been perfomed already
		//then it branches to either loop_body or loop_exit
		//loop_body updates current_value and increments counter
		//the it branches back to loop_validate

		//our blocks
		loop_validate := f.NewBlock("")
		loop_body := f.NewBlock("")
		loop_exit := f.NewBlock("")

		//setup code i.e. b code
		initial_value := represent(f, b, params[1:])
		(*b).NewBr(loop_validate)

		//loop_validate code
		current_value := loop_validate.NewPhi(ir.NewIncoming(initial_value, *b))
		counter := loop_validate.NewPhi(ir.NewIncoming(constant.NewInt(types.I32, 0), *b))
		stop_condition := loop_validate.NewICmp(enum.IPredULT, counter, params[0])
		loop_validate.NewCondBr(stop_condition, loop_body, loop_exit)

		//loop_body code
		next_value := represent(f, &loop_body, append([]value.Value{counter, current_value}, params[1:]...))
		next_counter := loop_body.NewAdd(counter, constant.NewInt(types.I32, 1))
		loop_body.NewBr(loop_validate)

		//here we can do those phi instruction from earlier
		current_value.Incs = append(current_value.Incs, ir.NewIncoming(next_value, loop_body))
		counter.Incs = append(counter.Incs, ir.NewIncoming(next_counter, loop_body))

		//loop_exit is where we continue from
		*b = loop_exit

		return current_value
	case min:
		//this case is extremely similar to rec
		//here we increment counter until result is 0

		//our blocks
		loop_body := f.NewBlock("")
		loop_exit := f.NewBlock("")

		(*b).NewBr(loop_body)

		//loop_body code
		counter := loop_body.NewPhi(ir.NewIncoming(constant.NewInt(types.I32, 0), *b))
		result := represent(f, &loop_body, append([]value.Value{counter}, params...))
		next_counter := loop_body.NewAdd(counter, constant.NewInt(types.I32, 1))
		exit_condition := loop_body.NewICmp(enum.IPredEQ, result, constant.NewInt(types.I32, 0))
		loop_body.NewCondBr(exit_condition, loop_exit, loop_body)

		counter.Incs = append(counter.Incs, ir.NewIncoming(next_counter, loop_body))

		*b = loop_exit

		return counter
	case identifier:
		id := t_to_c.get()
		return (*b).NewCall(func_list[id], params...)
	default:
		panic("invalid token in t_to_c")
	}
	return params[0]
}
