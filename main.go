//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
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
	//representation tokens
	mov         = iota
        add         = iota
        sub         = iota
	inc         = iota
	cmp         = iota
        branch      = iota
	beq         = iota
	label       = iota
	ret         = iota
	register    = iota
	stack       = iota
	stack_offset = iota
	str         = iota
	load        = iota
)

//these are the channels the phases use to communicate with one another
//f=file, s=scanner, p=parser, t=semantic, c=code_gen
var f_to_s Stream
var s_to_p Stream
var p_to_t Stream
var t_to_r Stream
var r_to_c Stream

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
		//as there are 4 phases with output, n must be 0, 1, 2, or 3
		if n < 0 || n > 3 {
			panic("invalid debug phase")
		}
	}
	//next we initialize the channels
	f_to_s = from_file(file)
	s_to_p = new_stream()
	p_to_t = new_stream()
	t_to_r = new_stream()
	r_to_c = new_stream()
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
	if n == 3 {
		go represent()
	}
	debug(n)
}
