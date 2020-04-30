package main

//import "fmt"
import "strconv"

func is_alphabetic(c byte) bool {
	//returns true if the input is a letter and false otherwise
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func is_digit(c byte) bool {
	return '0' <= c && c <= '9'
}

//scan() is the function that instantiates the scanner
//it recieves its input from f_to_s and send output to s_to_p
//the scanner consists of a loop that finds one token each iteration
//it invokes auxiliary functions for more complex token
func scan() {
	//the name table is used in screening, it maps each user defined function name to a unique integer
	name_table := make(map[string]byte)
	//main is always mapped to 0, this fact will be used in subsequent phases
	name_table["main"] = 0
	for {
		b := f_to_s.get()
		if b == '/' {
			scan_comment()
		} else if is_digit(b) {
			f_to_s.undo()
			scan_constant()
		} else if is_alphabetic(b) || b == '_' {
			f_to_s.undo()
			scan_identifier(name_table)
		} else if b == '(' {
			s_to_p.put(open_paren)
		} else if b == ')' {
			s_to_p.put(close_paren)
		} else if b == '=' {
			s_to_p.put(equals)
		} else if b == '\n' {
			s_to_p.put(newline)
		} else if b == ',' {
			s_to_p.put(comma)
		} else if b == '\x00' {
			//\0 signifies the end
			s_to_p.put(end)
			break
		} else if b != ' ' {
			panic("unexpected symbol")
		}
	}
}

//scan_comment expects the second slash, then consumes everything up to,
//but not include the next newline
func scan_comment() {
	if f_to_s.get() != '/' {
		panic("expected slash")
	}
	for f_to_s.get() != '\n' {
	}
	f_to_s.undo()
}

//scan_constant consumes a sequence of digits, then emits the resulting number
func scan_constant() {
	buffer := ""
	for {
		b := f_to_s.get()
		if is_digit(b) {
			buffer += string(b)
		} else {
			value, _ := strconv.Atoi(buffer)
			s_to_p.put(constant)
			s_to_p.put(byte(value))
			f_to_s.undo()
			break
		}
	}
}

//scan_identifier consumes a sequence of letters digits and underscores
//if the result is a keyword, that keyword's token is emmited
//otherwise the name is mapped to a number using the name table and that is emmited
func scan_identifier(name_table map[string]byte) {
	buffer := ""
	for {
		//fmt.Println(f_to_s.buf_in_use)
		b := f_to_s.get()
		if is_digit(b) || is_alphabetic(b) || b == '_' {
			buffer += string(b)
		} else {
			if buffer == "suc" {
				s_to_p.put(suc)
			} else if buffer == "proj" {
				s_to_p.put(proj)
			} else if buffer == "rec" {
				s_to_p.put(rec)
			} else if buffer == "comp" {
				s_to_p.put(comp)
			} else if buffer == "min" {
				s_to_p.put(min)
			} else {
				//fmt.Println(buffer)
				//if it is not a keyword, we handle it using the name_table
				s_to_p.put(identifier)
				n, ok := name_table[buffer]
				//ok is false if this is the first time the name is seen
				if !ok {
					//we assign the size of the table to the new identifier
					//this will always result in unique values
					n = byte(len(name_table))
					name_table[buffer] = n
				}
				s_to_p.put(n)
			}
			f_to_s.undo()
			break
		}
	}
}
