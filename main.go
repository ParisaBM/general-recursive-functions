//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
//	"os"
	"fmt"
//	"io/ioutil"
)

/*func is_alphabetic(c byte) bool {
	//returns true if the input is a letter and false otherwise
	return ('A' < c && c < 'Z') || ('a' < c && c < 'z')
}

func scan(s []byte) {
	//reads the inputs one byte at a time and emits what it finds
	//currently it only finds comments, and ignores their contents

	//slashes appear in pairs
	//expected_slash is true if a slash was just read
	expect_slash := false
	//comment is true if the scanner is within a comment
	comment := false
	buffer := ""
	identifier := false
	constant := false
	for i := range s {
		//first non-newline characters withing a comment are ignored
		if comment {
			if s[i] == '\n' {
				comment = false
			}
		//next the beginning of comments is handled
		} else if expect_slash {
			if s[i] == '/' {
				fmt.Println("comment")
				comment = true
				expect_slash = false
			} else {
				fmt.Println("error: bad comment")
				return
			}
		}
		} else if s[i] == '/' {
			expect_slash = true
		//next special symbols are found
		} else if s[i] == '(' {
			fmt.Println("open bracket")
		} else if s[i] == ')' {
			fmt.Println("close bracket")
		} else if s[i] == '=' {
			fmt.Println("equals")
		} else if s[i] == ',' {
			fmt.Println("comma")
		//any other symbols are a syntax error
		} else if s[i] != ' ' || s[i] != '\n' {
			fmt.Println("invalid symbol ", s[i])
			return
		}
	}
}*/

func main() {
	//this program will sometimes be run as its executable, and sometimes with go run
	//hence the filename is assumed to be the last arguement
	//this feature will be removed in the future
	//it is beyond the scope of this project to support multi-file codebases
	s := "abcdefghijklmnopqrstuvwxyz"
	for i := range s {
		fmt.Println(i)
		i+=1
	}

	//filename := os.Args[len(os.Args)-1]
	//file, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	panic(err)
	//}
	//scan(file)
}
