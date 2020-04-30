package main

//On Hold
/*import "fmt"

//we require a way to systematically generate new labels
//we do this by using the number of labels that have been used (labels_used)
//then incrementing labels_used
//this results in labels being assign sequentially
//this variable will be passed by reference so that any function generating code can use it
labels_used := int8(0)
//every named function has a label indicating where it starts
//these are kept in the label table
label_table := make(chan[int8]int8)
//these are global variables because the type signature of represent_funtion() would be ungodly otherwise
func represent() {
	L: for {
		//each iteration of this loop processes one line of code
		switch t_to_r {
		case end:
			break L
		case identifier:
			//first the function is given a label
			fmt.Println("L", labels_used)
			label_table[<- t_to_r] = labels_used
			labels_used++
			//next we accumulate the instructions in the function definition
			arity := <- t_to_r
			args := make([]int8, 0)
			for i := 1; i <= arity; i++ {
				fmt.Println("pop R", i)
				args = append(args, i)
				protected = append(protected, i)
			}
			represent_function(0, arity, args, none)
			}
		}
	}
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
func represent_function(target int8, protected int8, args []int8, b int8) {
	if b != none {
		b = <- t_to_r
	}
	switch b {
	case constant:
		fmt.Println("MOV R", target, " ", <- t_to_r)
	case suc:
		fmt.Pritnln("ADD R", target, " R", args[0], " 1")
	case proj:
		fmt.Println("MOV R", target, " R", args[<- t_to_r])
	case prefix_comp:

	case prefix_min:
		fmt.Println("MOV R", target, " 0")
		//min requires a loop
		//for which we allocate a new label
		l := labels_used
		labels_used++
		fmt.Println("L", l, ":")
		represent_function(protected+1, protected+1, append(target, args), none)
		fmt.Println("ADD R", target, " 1")
		fmt.Println("CMP R", protected+1, " 0")
		fmt.Println("BNE L", l)
		//the loop over-increments by 1
		fmt.Println("SUB R", target, " 1")
	case prefix_rec:
		represent_funtion(protected+1, protected+1, args[1:], none)
		fmt.Println("MOV R", protected+2, " 0")
		//allocate 2 labels, l and l+1
		l := labels_used
		labels_used += 2
		fmt.Println("L", l, ":")
		fmt.Println("MOV R", target, " R", protected+1)
		fmt.Println("CMP R", protected+2, " R", args[0])
		fmt.Println("BEQ L", l+1)
		represent_function(protected+1, protected+2, append(protected+2, target, args), none)
		fmt.Println("ADD R", protected+2, " 1")
		fmt.Println("B L", l)
		fmt.Println("L", l+1, ":")
	}
}*/
