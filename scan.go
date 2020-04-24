package main

import "fmt"

func scan(s string, s_to_p chan byte) {
	//reads the inputs one byte at a time and emits what it finds
	//scanner flags are named as: X_f

	//comment is true if the scanner is within a comment
	comment_f := false

	//slashes appear in pairs
	//expected_slash is true if a slash was just read
	expect_slash_f := false


	//if the scanner is within a constant, then constant_fs is true, and const buffer saves the result
	constant_f := false
	var constant_buffer byte = 0

	//identifiers are treated likewise
	identifier_f := false
	id_buffer := ""
	
	//the name table is used in screening, it maps each user defined function name to a unique integer
	name_table := make(map[string]byte)
	name_table["main"]=0
	for i := 0; i < len(s); i++ {
		//anywhere there is an i-- this means the character is not consumed

		//this loop is divided into 2 parts
		//first all the cases where one of the boolean flags is set are handled\
		//otherwise the scanner is between tokens and it reads the character to determine how to proceed
		
		//here is the first section, where flags are checked
		if comment_f {
			if s[i] == '\n' {
				comment_f = false
				i--
			}
		} else if expect_slash_f {
			if s[i] == '/' {
				comment_f = true
				expect_slash_f = false
			} else {
				fmt.Println("error: bad comment")
				s_to_p <- err
				return
			}
		} else if constant_f {
			if is_digit(s[i]) {
				constant_buffer = (constant_buffer*10)+s[i]-'0'
			} else if is_alphabetic(s[i]) || s[i] == '_' {
				fmt.Println("error: bad identifier")
				s_to_p <- err
				return
			} else {
				s_to_p <- constant
				s_to_p <- constant_buffer
				//the buffer is reset (effectively, the content will be overwritten later if it needs to be)
				constant_f = false
				i--
			}
		} else if identifier_f {
			if is_alphabetic(s[i]) || is_digit(s[i]) || s[i] == '_' {
				id_buffer += string(s[i])
			} else {
				//fist the identifier is compared agains all the keywords
				//main for our purposes is a keyword
				if id_buffer=="suc" {
					s_to_p <- suc
				} else if id_buffer=="proj" {
					s_to_p <- proj
				} else if id_buffer=="rec" {
					s_to_p <- rec
				} else if id_buffer=="comp" {
					s_to_p <- comp
				} else if id_buffer=="min" {
					s_to_p <- min
				//if it is not a keyword, we handle it using the name_table
				} else {
					s_to_p <- identifier
					n, ok := name_table[id_buffer]
					//ok is false if this is the first time the name is seen
					if !ok {
						//we assign the size of the table to the new identifier
						//this will always result in unique values
						n = byte(len(name_table))
						name_table[id_buffer]=n
					}
					s_to_p <- n
				}
				identifier_f = false
				i--
			}
		//here is the second section, where different symbols are handled
		} else if is_alphabetic(s[i]) || s[i]=='_' {
			id_buffer = string(s[i])
			identifier_f = true
		} else if is_digit(s[i]) {
			constant_buffer = s[i]-'0'
			constant_f = true
		} else if s[i] == '/' {
			expect_slash_f = true
		} else if s[i] == '(' {
			s_to_p <- open_paren
		} else if s[i] == ')' {
			s_to_p <- close_paren
		} else if s[i] == '=' {
			s_to_p <- equals
		} else if s[i] == ',' {
			s_to_p <- comma
		} else if s[i] == '\n' {
			s_to_p <- newline
		} else if s[i] != ' ' {
			fmt.Println("error: invalid symbol ", s[i])
			s_to_p <- err
			return
		}
	}
	s_to_p <- end
}
