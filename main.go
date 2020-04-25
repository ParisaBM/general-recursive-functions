//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
	"io/ioutil"
	"os"
)

//these are the signals the phases use to communicate with one another
const (
	constant    = iota
	suc         = iota
	proj        = iota
	comp        = iota
	min         = iota
	rec         = iota
	identifier  = iota
	equals      = iota
	open_paren  = iota
	close_paren = iota
	comma       = iota
	newline     = iota
	err         = iota
	end         = iota
)

//these are the channels the phases use to communicate with one another
//s=scanner, p=parser, t=semantic, c=code_gen
var s_to_p chan byte
var p_to_t chan byte
var t_to_c chan byte

func main() {
	//this program will sometimes be run as its executable, and sometimes with go run
	//hence the filename is assumed to be the last arguement
	//this feature will be removed in the future
	//it is beyond the scope of this project to support multi-file codebases
	filename := os.Args[len(os.Args)-1]
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	//next we initialize the channels
	s_to_p = make(chan byte)
	p_to_t = make(chan byte)
	//then we begin the phases
	go scan(string(file))
	go parse()
	go semantic()
	<- t_to_c
}
