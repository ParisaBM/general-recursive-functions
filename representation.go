package main

//we require a way to systematically generate new labels
//we do this by using the number of labels that have been used (labels_used)
//then incrementing labels_used
//this results in labels being assign sequentially
//this variable will be passed by reference so that any function generating code can use it
var labels_used byte
//every named function has a label indicating where it starts
//these are kept in the label table
var label_table map[byte]byte
//each register may or may not be on the stack
//when a function call is made, it has to know which need to be saved, and which do not
//reg_location keeps track of where each register is from the top of the stack, if it is on the stack
//every time a register is updated in any way it gets deleted from reg_locations
var reg_locations map[byte]byte
//reg_usage tracks which registers a particular function modifies
//in particular it maps a function id, the set of registers it uses
//a map onto the unit type is effectively a set
//https://en.wikipedia.org/wiki/Unit_type
var reg_usage map[byte]map[byte]struct{}
//stack_usage is how much space on the stack a particular function uses
var stack_usage map[byte]byte
//current_id is the id of the function that is being represented
var current_id byte
//these are global variables because the type signature of represent_funtion() would be ungodly otherwise
func represent() {
	labels_used = 0
	label_table = make(map[byte]byte)
	reg_usage = make(map[byte]map[byte]struct{})
	stack_usage = make(map[byte]byte)
	L: for {
		//each iteration of this loop processes one line of code
		switch t_to_r.get() {
		case identifier:
			//first the function is given a label
			r_to_c.put(label)
			r_to_c.put(labels_used)
			current_id = t_to_r.get()
			label_table[current_id] = labels_used
			labels_used++
			//next we accumulate the instructions in the function definition
			ar := t_to_r.get()
			args := make([]byte, 0)
			//reg_locations is reset
			reg_locations = make(map[byte]byte)
			for i := byte(0); i < ar; i++ {
				args = append(args, i+1)
			}
			reg_usage[current_id] = make(map[byte]struct{})
			stack_usage[current_id] = 0
			r_to_c.delimit_buffering() //begin
			represent_function(0, ar, args)
			r_to_c.delimit_buffering() //end
			if stack_usage[current_id] != 0 {
				r_to_c.put(add)
				r_to_c.put(stack)
				r_to_c.put(stack)
				r_to_c.put(constant)
				r_to_c.put(stack_usage[current_id])
			}
			r_to_c.put_buffer()
			if stack_usage[current_id] != 0 {
				r_to_c.put(sub)
				r_to_c.put(stack)
				r_to_c.put(stack)
				r_to_c.put(constant)
				r_to_c.put(stack_usage[current_id])
			}
			r_to_c.put(ret)
		case end:
			r_to_c.put(end)
			break L
		}
	}
	r_to_c.put(0)
}

//represent_fucntion recursively generates psuedo-assembly for a function
//the first arguement, target is the register where the result should be stored
//protected, is the index of the last register that is protected
//for example if protected=3, R0, R1, R2, and R3 cannot be used except for
//whichever of those registers is the target
//args is the in-order sequence of registers where the arguements are
//occasionally we will have to peek at the input stream, then put the value back
//if this is the case, the value will be put back into b
//b is none if nothing was put back this way
func represent_function(target byte, protected byte, args []byte) {
	switch t_to_r.get() {
	case constant:
		r_to_c.put(mov)
		r_to_c.put(register)
		r_to_c.put(target)
		r_to_c.put(constant)
		r_to_c.put(t_to_r.get())
	case suc:
		r_to_c.put(add)
		r_to_c.put(register)
		r_to_c.put(target)
		r_to_c.put(register)
		r_to_c.put(args[0])
		r_to_c.put(constant)
		r_to_c.put(1)
	case proj:
		r_to_c.put(mov)
		r_to_c.put(register)
		r_to_c.put(target)
		r_to_c.put(register)
		r_to_c.put(args[t_to_r.get()])
	case comp:
		n := t_to_r.get()
		//new_args are the arguements of the top-level function
		new_args := make([]byte, 0)
		for i := byte(0); i < n; i++ {
			if t_to_r.get() == proj {
				new_args = append(new_args, args[t_to_r.get()])
			} else {
				t_to_r.undo()
				protected++
				new_args = append(new_args, protected)
				represent_function(protected, protected, args)
			}
		}
		represent_function(target, protected, new_args)
	case min:
		//counter is the variable we tick up to until we find the minimum
		var counter byte
		if contains(args, target) {
			protected++
			counter = protected
			reg_usage[current_id][protected] = struct{}{}
		} else {
			counter = target
		}
		r_to_c.put(mov)
		r_to_c.put(register)
		r_to_c.put(target)
		r_to_c.put(constant)
		r_to_c.put(0)
		delete(reg_locations, target)
		//min requires a loop
		//for which we allocate new labels
		//allocate 2 labels, l and l+1
		l := labels_used
		labels_used+=2
		r_to_c.put(label)
		r_to_c.put(l)
		represent_function(protected+1, protected+1, append([]byte{counter}, args...))
		r_to_c.put(cmp)
		r_to_c.put(register)
		r_to_c.put(protected+1)
		r_to_c.put(constant)
		r_to_c.put(0)
		r_to_c.put(beq)
		r_to_c.put(l+1)
		r_to_c.put(inc)
		r_to_c.put(register)
		r_to_c.put(target)
		delete(reg_locations, target)
		r_to_c.put(branch)
		r_to_c.put(l)
		r_to_c.put(label)
		r_to_c.put(l+1)
		if counter != target {
			r_to_c.put(mov)
			r_to_c.put(register)
			r_to_c.put(target)
			r_to_c.put(register)
			r_to_c.put(counter)
		}
	case rec:
		//sub_target is where the base case and recursive cases should put their results
		var sub_target byte
		if contains(args, target) {
			protected++
			sub_target = protected
		} else {
			sub_target = target
		}
		represent_function(protected, sub_target, args[1:])
		r_to_c.put(mov)
		r_to_c.put(register)
		r_to_c.put(protected+1)
		r_to_c.put(constant)
		r_to_c.put(0)
		delete(reg_locations, protected+1)
		//allocate 2 new labels
		l := labels_used
		labels_used += 2
		r_to_c.put(label)
		r_to_c.put(l)
		if target != sub_target {
			r_to_c.put(mov)
			r_to_c.put(register)
			r_to_c.put(target)
			r_to_c.put(register)
			r_to_c.put(sub_target)
			delete(reg_locations, target)
		}
		r_to_c.put(cmp)
		r_to_c.put(register)
		r_to_c.put(protected+1)
		r_to_c.put(register)
		r_to_c.put(args[0])
		r_to_c.put(beq)
		r_to_c.put(l+1)
		represent_function(sub_target, protected+1, append([]byte{protected+1, target}, args...))
		r_to_c.put(inc)
		r_to_c.put(register)
		r_to_c.put(protected+1)
		delete(reg_locations, protected+1)
		r_to_c.put(branch)
		r_to_c.put(l)
		r_to_c.put(label)
		r_to_c.put(l+1)
		reg_usage[current_id][protected+1] = struct{}{}
	case identifier:
		callee_id := t_to_r.get()
		//moves all the registers to the stack that need to be protected
		for i := byte(0); i <= protected; i++ {
			_, saved := reg_locations[i]
			_, needs_saving := reg_usage[callee_id][i]
			if !saved && needs_saving {
				r_to_c.put(str)
				r_to_c.put(stack_offset)
				r_to_c.put(stack_usage[current_id])
				r_to_c.put(register)
				r_to_c.put(i)
				reg_locations[i] = stack_usage[current_id]
				stack_usage[current_id]++
			}
		}
		//the arguements to the function are loaded
		for i, reg := range args {
			if reg != byte(i)+1 {
				r_to_c.put(load)
				r_to_c.put(register)
				r_to_c.put(byte(i)+1)
				r_to_c.put(stack_offset)
				r_to_c.put(reg_locations[reg])
			}
		}
		//the external function is invoked
		r_to_c.put(branch)
		r_to_c.put(callee_id)
		//its result is saved in target (if applicable)
		if target != 0 {
			r_to_c.put(mov)
			r_to_c.put(register)
			r_to_c.put(target)
			r_to_c.put(register)
			r_to_c.put(0)
		}
		//restore the state
		for reg, loc := range reg_locations {
			if reg != target && reg <= protected {
				r_to_c.put(load)
				r_to_c.put(register)
				r_to_c.put(reg)
				r_to_c.put(stack_offset)
				r_to_c.put(loc)
			}
		}
	}
	delete(reg_locations, target)
	reg_usage[current_id][target] = struct{}{}
}

//checks if b is in l
func contains(l []byte, b byte) bool {
	for _, v := range l {
		if b == v {
			return true
		}
	}
	return false
}
