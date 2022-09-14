package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var func_list []*ir.Func

func code_gen() {
	u32 := types.I32
	m := ir.NewModule()
	for {
		//t_to_c should consist of identifier with their function definitions then an end token
		//anything else means something's gone wrong
		switch t_to_c.get() {
		case identifier:
			id := t_to_c.get()
			param_count := t_to_c.get()
			params := make([]*ir.Param, 0)
			for i := byte(0); i < param_count; i += 1 {
				params = append(params, ir.NewParam("p"+fmt.Sprint(i), u32))
			}
			name_list_mutex.Lock()
			f := m.NewFunc(name_list[id], u32, params...)
			name_list_mutex.Unlock()
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
			func_list = append(func_list, f)
		case end:
			fmt.Println(m)
			c_to_e.put(0)
			return
		default:
			panic("invalid token in t_to_c")
		}
	}
}

//represent recursively generates the llvm ir for a single function definition
//these instructions all get added to the first parameter f
//b is the current active block
//b is doubly indirect so that instances of represent can update it, and when they do other instances further down in
//the stack can add instructions to the block where it was left off
//params are the parameters to that sub function
//what gets returned is either a register or constant with the result of that function
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
		counter_pointer := (*b).NewAlloca(types.I32)
		(*b).NewStore(constant.NewInt(types.I32, 0), counter_pointer)
		current_value_pointer := (*b).NewAlloca(types.I32)
		current_value := represent(f, b, params[1:])
		(*b).NewStore(current_value, current_value_pointer)
		(*b).NewBr(loop_validate)

		//loop_validate code
		counter := loop_validate.NewLoad(types.I32, counter_pointer)
		stop_condition := loop_validate.NewICmp(enum.IPredULT, counter, params[0])
		loop_validate.NewCondBr(stop_condition, loop_body, loop_exit)

		//loop_body code
		current_value = loop_body.NewLoad(types.I32, current_value_pointer)
		current_value = represent(f, &loop_body, append([]value.Value{counter, current_value}, params[1:]...))
		if current_value == counter {
			//this is somewhat of a corner case
			//if these are same counter will get updated once more which will be an issue
			current_value = loop_body.NewAdd(counter, constant.NewInt(types.I32, 0))
		}
		loop_body.NewStore(current_value, current_value_pointer)
		//note that the old counter value is still in a register
		incremented_counter := loop_body.NewAdd(counter, constant.NewInt(types.I32, 1))
		loop_body.NewStore(incremented_counter, counter_pointer)
		loop_body.NewBr(loop_validate)

		//loop_exit is where we continue from
		*b = loop_exit

		return current_value
	case min:
		//this case is extremely similar to rec
		//here we increment counter until result is 0

		//our blocks
		loop_body := f.NewBlock("")
		loop_exit := f.NewBlock("")

		//setup code
		counter_pointer := (*b).NewAlloca(types.I32)
		//we initialize the counter to -1 so we can increment it at the start of the loop
		(*b).NewStore(constant.NewInt(types.I32, -1), counter_pointer)
		(*b).NewBr(loop_body)

		//loop_body code
		counter := loop_body.NewLoad(types.I32, counter_pointer)
		incremented_counter := loop_body.NewAdd(counter, constant.NewInt(types.I32, 1))
		loop_body.NewStore(incremented_counter, counter_pointer)
		result := represent(f, &loop_body, append([]value.Value{incremented_counter}, params...))
		exit_condition := loop_body.NewICmp(enum.IPredEQ, result, constant.NewInt(types.I32, 0))
		loop_body.NewCondBr(exit_condition, loop_exit, loop_body)

		*b = loop_exit

		return incremented_counter
	case identifier:
		id := t_to_c.get()
		return (*b).NewCall(func_list[id], params...)
	default:
		panic("invalid token in t_to_c")
	}
	return params[0]
}
