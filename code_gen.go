package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func code_gen() {
	u32 := types.I32
	m := ir.NewModule()
	t_to_c.get()
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
	entry := f.NewBlock("entry")
	ret := represent(f, entry, param_registers)
	entry.NewRet(ret)
	for {
		if t_to_c.get() == end {
			break
		}
	}
	fmt.Println(m)
	c_to_e.put(0)
}

func represent(f *ir.Func, b *ir.Block, params []value.Value) value.Value {
	// b := f.NewBlock("")
	// b.NewRet(params[0])
	switch t_to_c.get() {
	case const_t:
		return constant.NewInt(types.I32, int64(t_to_c.get()))
	case proj:
		return params[t_to_c.get()]
	case suc:
		new_value := b.NewAdd(params[0], constant.NewInt(types.I32, 1))
		return new_value
	case comp:
		top_level_params := make([]value.Value, 0)
		num_top_level := t_to_c.get()
		for i := byte(0); i < num_top_level; i += 1 {
			top_level_params = append(top_level_params, represent(f, b, params))
		}
		return represent(f, b, top_level_params)
	case rec:
		current_value := represent(f, b, params[1:])
		iteration := constant.NewInt(types.I32, 0)
		loop := f.NewBlock("rec")
	default:
		panic("not a valid token in t_to_c")
	}
	return params[0]
}
