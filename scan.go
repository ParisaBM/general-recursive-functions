package main

import "strconv"

func isAlphabetic(c byte) bool {
	// returns true if the input is a letter and false otherwise
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// scan() is the function that instantiates the scanner
// it recieves its input from fToS and send output to sToP
// the scanner consists of a loop that finds one token each iteration
// it invokes auxiliary functions for more complex token
func scan() {
	// the name table is used in screening, it maps each user defined function name to a unique integer
	nameTable := make(map[string]byte)
	for {
		b := fToS.get()
		if b == '/' {
			scanComment()
		} else if isDigit(b) {
			fToS.undo()
			scanConstant()
		} else if isAlphabetic(b) || b == '_' {
			fToS.undo()
			scanIdentifier(nameTable)
		} else if b == '(' {
			sToP.put(openParen)
		} else if b == ')' {
			sToP.put(closeParen)
		} else if b == '=' {
			sToP.put(equals)
		} else if b == '\n' {
			sToP.put(newline)
		} else if b == ',' {
			sToP.put(comma)
		} else if b == '\x00' {
			// \0 signifies the end
			sToP.put(end)
			break
		} else if b != ' ' {
			panic("unexpected symbol")
		}
	}
}

// scanComment expects the second slash, then consumes everything up to,
// but not include the next newline
func scanComment() {
	if fToS.get() != '/' {
		panic("expected slash")
	}
	for fToS.get() != '\n' {
	}
	fToS.undo()
}

// scanConstant consumes a sequence of digits, then emits the resulting number
func scanConstant() {
	buffer := ""
	for {
		b := fToS.get()
		if isDigit(b) {
			buffer += string(b)
		} else {
			value, _ := strconv.Atoi(buffer)
			sToP.put(constT)
			sToP.put(byte(value))
			fToS.undo()
			break
		}
	}
}

// scanIdentifier consumes a sequence of letters digits and underscores
// if the result is a keyword, that keyword's token is emmited
// otherwise the name is mapped to a number using the name table and that is emmited
func scanIdentifier(nameTable map[string]byte) {
	buffer := ""
	for {
		b := fToS.get()
		if isDigit(b) || isAlphabetic(b) || b == '_' {
			buffer += string(b)
		} else {
			if buffer == "suc" {
				sToP.put(suc)
			} else if buffer == "proj" {
				sToP.put(proj)
			} else if buffer == "rec" {
				sToP.put(rec)
			} else if buffer == "comp" {
				sToP.put(comp)
			} else if buffer == "min" {
				sToP.put(min)
			} else {
				// if it is not a keyword, we handle it using the nameTable
				sToP.put(identifier)
				n, ok := nameTable[buffer]
				// ok is false if this is the first time the name is seen
				if !ok {
					// we assign the size of the table to the new identifier
					// this will always result in unique values
					n = byte(len(nameTable))
					nameTable[buffer] = n
					nameListMutex.Lock()
					nameList = append(nameList, buffer)
					nameListMutex.Unlock()
				}
				sToP.put(n)
			}
			fToS.undo()
			break
		}
	}
}
