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

func scan(s []byte) {
	//reads the inputs one byte at a time and emits what it finds

	//comment is true if the scanner is within a comment
	comment := false

	//slashes appear in pairs
	//expected_slash is true if a slash was just read
	expect_slash := false


	//if the scanner is within a constant, then constant is true, and const buffer saves the result
	constant := false
	var const_buffer byte = 0

	identifier := false
	id_buffer := []byte{}
	for i := 0; i < len(s); i++ {
		//anywhere there is an i-- this means the character is not consumed

		//this loop is divided into 2 parts
		//first all the cases where one of the boolean flags is set are handled\
		//otherwise the scanner is between tokens and it reads the character to determine how to proceed
		
		//here is the first section, where flags are checked
		if comment {
			if s[i] == '\n' {
				comment = false
				i--
			}
		} else if expect_slash {
			if s[i] == '/' {
				fmt.Println("comment")
				comment = true
				expect_slash = false
			} else {
				fmt.Println("error: bad comment")
				return
			}
		} else if constant {
			if is_digit(s[i]) {
				const_buffer = (const_buffer*10)+s[i]-'0'
			} else if is_alphabetic(s[i]) || s[i] == '_' {
				fmt.Println("bad identifier")
				return
			} else {
				fmt.Println("found constant:")
				fmt.Println(const_buffer)
				const_buffer = 0
				constant = false
				i--
			}
		} else if identifier {
			if is_alphabetic(s[i]) || is_digit(s[i]) || s[i] == '_' {
				id_buffer = append(id_buffer, s[i])
			} else {
				fmt.Println("found identifier:")
				fmt.Println(id_buffer)
				id_buffer = []byte{}
				identifier = false
				i--
			}
		//here is the second section, where different symbols are handled
		} else if is_alphabetic(s[i]) || s[i]=='_' {
			id_buffer = append(id_buffer, s[i])
			identifier = true
		} else if is_digit(s[i]) {
			const_buffer = s[i]-'0'
			constant = true
		} else if s[i] == '/' {
			expect_slash = true
		} else if s[i] == '(' {
			fmt.Println("open bracket")
		} else if s[i] == ')' {
			fmt.Println("close bracket")
		} else if s[i] == '=' {
			fmt.Println("equals")
		} else if s[i] == ',' {
			fmt.Println("comma")
		} else if s[i] == '\n' {
			fmt.Println("line break")
		} else if s[i] != ' ' {
			fmt.Println("invalid symbol ", s[i])
			return
		}
	}
}

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
	scan(file)
}
