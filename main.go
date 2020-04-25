//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func is_alphabetic(c byte) bool {
	//returns true if the input is a letter and false otherwise
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func is_digit(c byte) bool {
	return '0' <= c && c <= '9'
}

const (
	//scanner signals
	constant = iota
	suc = iota
	proj = iota
	comp = iota
	min = iota
	rec = iota
	identifier = iota
	equals = iota
	open_paren = iota
	close_paren = iota
	comma = iota
	newline = iota
	err = iota
	end = iota
)

var s_to_p chan byte
var p_to_t chan byte

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
	s_to_p = make(chan byte)
	p_to_t = make(chan byte)
	go scan(string(file))
	go parse()
	for {
		c := <- p_to_t
		fmt.Println(c)
		if c == end {
			break
		}
	}
}
