package main

import "fmt"

//we require a way to systematically generate new labels
//we do this by using the number of labels that have been used (labels_used)
//then incrementing labels_used
//this results in labels being assign sequentially
//this variable will be passed by reference so that any function generating code can use it
var labels_used byte
//every named function has a label indicating where it starts
//these are kept in the label table
var label_table map[byte]byte
//these are global variables because the type signature of represent_funtion() would be ungodly otherwise
func represent() {
	labels_used = 0
	label_table = make(map[byte]byte)
	L: for {
		//each iteration of this loop processes one line of code
		switch t_to_r.get() {
		case identifier:
			//first the function is given a label
			fmt.Printf("L%d:\n", labels_used)
			label_table[t_to_r.get()] = labels_used
			labels_used++
			//next we accumulate the instructions in the function definition
			ar := t_to_r.get()
			args := make([]byte, 0)
			for i := byte(0); i < ar; i++ {
				fmt.Printf("LOAD R%d [STACK+%d]\n", i+1, ar-i-1)
				args = append(args, i+1)
			}
			fmt.Printf("ADD STACK %d\n", ar)
			represent_function(0, ar, args)
			fmt.Printf("RET\n")
		case end:
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
		fmt.Printf("MOV R%d %d\n", target, t_to_r.get())
	case suc:
		fmt.Printf("ADD R%d R%d 1\n", target, args[0])
	case proj:
		fmt.Printf("MOV R%d R%d\n", target, args[t_to_r.get()])
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
		fmt.Printf("MOV R%d 0\n", target)
		//min requires a loop
		//for which we allocate new labels
		//allocate 2 labels, l and l+1
		l := labels_used
		labels_used+=2
		fmt.Printf("L%d: ", l)
		represent_function(protected+1, protected+1, append([]byte{target}, args...))
		fmt.Printf("CMP R%d 0\n", protected+1)
		fmt.Printf("BEQ L%d\n", l+1)
		fmt.Printf("ADD R%d 1\n", target)
		fmt.Printf("B L%d\n", l)
		fmt.Printf("L%d:\n", l+1)		
	case rec:
		represent_function(protected+1, protected+1, args[1:])
		fmt.Printf("MOV R%d 0\n", protected+2)
		//allocate 2 new labels
		l := labels_used
		labels_used += 2
		fmt.Printf("L%d:\n", l)
		fmt.Printf("MOV R%d R%d\n", target, protected+1)
		fmt.Printf("CMP R%d R%d\n", protected+2, args[0])
		fmt.Printf("BEQ L%d\n", l+1)
		represent_function(protected+1, protected+2, append([]byte{protected+2, target}, args...))
		fmt.Printf("ADD R%d 1\n", protected+2)
		fmt.Printf("B L%d\n", l)
		fmt.Printf("L%d:\n", l+1)
	case identifier:
		fmt.Printf("SUB STACK %d\n", protected+1+byte(len(args)))
		for i := byte(0); i <= protected; i++ {
			if i < target {
				fmt.Printf("STR [STACK+%d] R%d\n", protected+byte(len(args))-i, i)
			} else if i > target {
				fmt.Printf("STR [STACK+%d] R%d\n", protected+1+byte(len(args))-i, i)				
			}
		}
		fmt.Printf("STR [STACK+%d] PC\n", len(args))
		for i, r := range args {
			fmt.Printf("STR [STACK+%d] R%d\n", len(args)-1-i, r)
		}
		fmt.Printf("B L%d\n", label_table[t_to_r.get()])
		if target != 0 {
			fmt.Printf("MOV R%d R0\n", target)
		}
		for i := byte(0); i <= protected; i++ {
			if i < target {
				fmt.Printf("LOAD R%d [STACK+%d]\n", i, protected-i)
			} else if i > target {
				fmt.Printf("LOAD R%d [STACK+%d]\n", i, protected+1-i)
			}
		}
		fmt.Printf("ADD STACK %d\n", protected+1)
	}
}
