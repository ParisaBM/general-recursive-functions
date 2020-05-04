package main

import "fmt"

func debug(n int) {
	//debug outputs the result of the nth phase
	//stream is whichever streameam it's supposed to be listening to
	var stream Stream
	switch n {
	case 0:
		stream = s_to_p
	case 1:
		stream = p_to_t
	case 2:
		stream = t_to_r
	case 3:
		stream = r_to_c
	}
	//next_id is used to determine if a particular identifier is a declaration
	//it is used when debugging the semantic phase
	next_id := byte(1)
	//in the representation phase different operations have different numbers of operands
	//operands_left is how many more need to be consumed before emitting a newline
	operands_left := 0
L:
	for {
		//there's a bit of code to handle compound tokens, but otherwise it just displays each token
		switch stream.get() {
		case constant:
			if n == 3 {
				fmt.Printf("%d ", stream.get())
				operands_left--
			} else {
				fmt.Println("constant")
				fmt.Println(stream.get())
			}
		case suc:
			fmt.Println("suc")
		case proj:
			fmt.Println("proj")
			//in the scanner output there are constant markers that are removed by the parser
			if n == 1 || n == 2 {
				fmt.Println(stream.get())
			}
			if n == 1 {
				fmt.Println(stream.get())
			}
		case comp:
			fmt.Println("comp")
			//in the parser, comp is followed by its arity
			if n == 1 || n == 2 {
				fmt.Println(stream.get())
			}
		case min:
			fmt.Println("min")
		case rec:
			fmt.Println("rec")
		case identifier:
			fmt.Println("identifier")
			id := stream.get()
			fmt.Println(id)
			if n == 2 && (id == next_id || id == 0) {
				fmt.Println(stream.get())
				next_id++
			}
		case equals:
			fmt.Println("equals")
		case open_paren:
			fmt.Println("open_paren")
		case close_paren:
			fmt.Println("close_paren")
		case comma:
			fmt.Println("comma")
		case newline:
			//the output is grouped by line
			fmt.Println("newline\n")
		case err:
			fmt.Println("err")
		case end:
			fmt.Println("end")
			break L
		case mov:
			fmt.Print("MOV ")
			operands_left = 2
		case add:
			fmt.Print("ADD ")
			operands_left = 3
		case sub:
			fmt.Print("SUB ")
			operands_left = 3
		case cmp:
			fmt.Print("CMP ")
			operands_left = 2
		case load:
			fmt.Print("LOAD ")
			operands_left = 2
		case str:
			fmt.Print("STR ")
			operands_left = 2
		case inc:
			fmt.Print("INC ")
			operands_left = 1
		case branch:
			fmt.Printf("B L%d", stream.get())
		case beq:
			fmt.Printf("BEQ L%d", stream.get())
		case label:
			fmt.Printf("L%d :", stream.get())
		case ret:
			fmt.Print("RET")
		case register:
			fmt.Printf("R%d ", stream.get())
			operands_left--
		case stack:
			fmt.Print("STACK ")
			operands_left--
		case stack_offset:
			fmt.Printf("[STACK+%d] ", stream.get())
			operands_left--
		}
		if n == 3 && operands_left == 0 {
			fmt.Println()
		}
	}
}
