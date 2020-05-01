//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
	"fmt"
	"os"
	"strconv"
)

//these are the signals the phases use to communicate with one another
const (
	//scanner and parser tokens
	constant   = iota
	suc        = iota
	proj       = iota
	comp       = iota
	min        = iota
	rec        = iota
	identifier = iota
	equals     = iota
	newline    = iota
	err        = iota
	end        = iota
	//scanner_tokens
	open_paren  = iota
	close_paren = iota
	comma       = iota
	//representation token
	none = iota
)

//these are the channels the phases use to communicate with one another
//f=file, s=scanner, p=parser, t=semantic, c=code_gen
var f_to_s Stream
var s_to_p Stream
var p_to_t Stream
var t_to_r Stream

func main() {
	//it is beyond the scope of this project to support multi-file codebases
	//there is an optional second arguement to debug a particular phase
	//if an integer is given as the fifth arguement, its respective phase
	//0=scan, 1=parse etc. doesn't output to the next phase and is instead
	//output to stdout with annotations
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	n := 3 //n is the phase the is being debbuged (if any)
	if len(os.Args) == 3 {
		n, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		//as there are 3 phases with output, n must be 0, 1 or 2
		if n < 0 || n > 2 {
			panic("invalid debug phase")
		}
	}
	//next we initialize the channels
	f_to_s = from_file(file)
	s_to_p = new_stream()
	p_to_t = new_stream()
	t_to_r = new_stream()
	//then we begin the phases
	if n >= 0 {
		go scan()
	}
	if n >= 1 {
		go parse()
	}
	if n >= 2 {
		go semantic()
	}
	if n != 3 {
		debug(n)
	} else {
		t_to_r.get()
	}
}

func debug(n int) {
	//debug outputs the result of the nth phase
	//ch is whichever channel it's supposed to be listening to
	var ch Stream
	switch n {
	case 0:
		ch = s_to_p
	case 1:
		ch = p_to_t
	case 2:
		ch = t_to_r
	}
	//tracks whether at the start of a new line
	//only relevant to semantic phase
	start_of_line := true
L:
	for {
		//there's a bit of code to handle compound tokens, but otherwise it just displays each token
		switch ch.get() {
		case constant:
			fmt.Println("constant")
			fmt.Println(ch.get())
		case suc:
			fmt.Println("suc")
		case proj:
			fmt.Println("proj")
			//in the scanner output there are constant markers that are removed by the parser
			if n == 1 || n == 2 {
				fmt.Println(ch.get())
				fmt.Println(ch.get())
			}
		case comp:
			fmt.Println("comp")
			//in the parser, comp is followed by its arity
			if n == 1 || n == 2 {
				fmt.Println(ch.get())
			}
		case min:
			fmt.Println("min")
		case rec:
			fmt.Println("rec")
		case identifier:
			fmt.Println("identifier")
			fmt.Println(ch.get())
			if n == 2 && start_of_line {
				fmt.Println(ch.get())
				start_of_line = false
			}
		case equals:
			fmt.Println("equals")
			start_of_line = true
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
		}
	}
}
