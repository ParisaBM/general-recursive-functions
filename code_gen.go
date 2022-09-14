package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
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
	represent(f, param_registers)
	for {
		if t_to_c.get() == end {
			break
		}
	}
	fmt.Println(m)
	c_to_e.put(0)
}

func represent(f *ir.Func, params []value.Value) {
	b := f.NewBlock("")
	b.NewRet(params[0])
	// 	switch t_to_c.get() {
	// 	case:
	// 	}
}
