//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

//these are the signals the phases use to communicate with one another
const (
	//scanner and parser tokens
	constant    = iota
	suc         = iota
	proj        = iota
	comp        = iota
	min         = iota
	rec         = iota
	identifier  = iota
	equals      = iota
	newline     = iota
	err         = iota
	end         = iota
	//scanner_tokens
	open_paren  = iota
	close_paren = iota
	comma       = iota
	//parser tokens
	prefix_comp = iota
	prefix_min  = iota
	prefix_rec  = iota
	//representation token
	none        = iota
)

//these are the channels the phases use to communicate with one another
//s=scanner, p=parser, t=semantic, c=code_gen
var s_to_p chan int8
var p_to_t chan int8
var t_to_r chan int8

func main() {
	//it is beyond the scope of this project to support multi-file codebases
	//there is an optional second arguement to debug a particular phase
	//if an integer is given as the fifth arguement, its respective phase
	//0=scan, 1=parse etc. doesn't output to the next phase and is instead
	//output to stdout with annotations
	filename := os.Args[1]
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	n := 2 //n is the phase the is being debbuged (if any)
	if len(os.Args) == 3 {
		n, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		//as there are 2 phases with output, n must be 0, or 1
		if n < 0 || n > 1 {
			panic("invalid debug phase")
		}
	}
	//next we initialize the channels
	s_to_p = make(chan int8)
	p_to_t = make(chan int8)
	t_to_r = make(chan int8)
	//then we begin the phases
	if n >= 0 {
		go scan(string(file))
	}
	if n >= 1 {
		go parse()
	}
	if n >= 2 {
		go semantic()
	}
	if n != 2 {
		debug(n)
	} else {
		<- t_to_r
	}
}

func debug(n int) {
	//debug outputs the result of the nth phase
	//ch is whichever channel it's supposed to be listening to
	var ch chan int8
	switch n {
	case 0:
		ch = s_to_p
	case 1:
		ch = p_to_t
	}
L:
	for {
		//there's a bit of code to handle compound tokens, but otherwise it just displays each token
		switch <-ch {
		case constant:
			fmt.Println("constant")
			fmt.Println(<-ch)
		case suc:
			fmt.Println("suc")
		case proj:
			fmt.Println("proj")
			//in the scanner output there are constant markers that are removed by the parser
			if n == 1 {
				fmt.Println(<-ch)
				fmt.Println(<-ch)
			}
		case comp:
			fmt.Println("comp")
			//in the parser, comp is followed by its arity
			if n == 1 {
				fmt.Println(<-ch)
			}
		case min:
			fmt.Println("min")
		case rec:
			fmt.Println("rec")
		case identifier:
			fmt.Println("identifier")
			fmt.Println(<-ch)
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
		}
	}
}
