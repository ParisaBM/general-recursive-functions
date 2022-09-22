package main

import "fmt"

func parse() {
	// the syntax of this language is highly regular
	// each line consists of code contains a definition
	// each iteration of this loop consumes one definition or a blank line, or the end of the file
	// since functions are meant to be on seperate lines, acceptingFunction keeps track of when a new line has started
	acceptingFunction := true
	for {
		switch sToP.get() {
		case identifier:
			// the general case is: identifer id equals function newline
			if !acceptingFunction {
				panic("multiple functions on the same line")
			}
			pToT.put(identifier)
			pToT.put(sToP.get())
			expect(equals)
			function()
			acceptingFunction = false
		case newline:
			pToT.put(newline)
			acceptingFunction = true
		case end:
			pToT.put(end)
			return
		default:
			// the errors that just say "error" are placeholders
			fmt.Println("error")
		}
	}
}

func function() {
	// consumes one function from the input stream
	// a function can be one of the primitives: suc, constant, proj
	// one of the operators: comp, min, rec
	// an identifer or a function enclosed in unnecessary brackets
	switch sToP.get() {
	// suc and constant are unchanged
	case suc:
		pToT.put(suc)
	case constT:
		pToT.put(constT)
		pToT.put(sToP.get())
	// the constant tags after proj are deleted from the stream
	// this is because they are redundant
	case proj:
		pToT.put(proj)
		expect(constT)
		pToT.put(sToP.get())
		expect(constT)
		pToT.put(sToP.get())
	// comp expects a comma seperated list of atleast one function in brackets
	case comp:
		pToT.put(comp)
		pToT.delimitBuffering() // begin
		expect(openParen)
		pToT.delimitBuffering() // begin
		function()              // this is the first function
		pToT.delimitBuffering() // end
		// a composition is emitted with its arity so the sematic analyzer can deduce
		// how many functions are being composed together
		var arity byte = 0
	L:
		for { // this loop consumes all addition functions
			switch sToP.get() {
			// if there's a comma there's an addition function
			case comma:
				function()
				arity++
			// if there's a closeParen, we've reached the end of the argument list
			case closeParen:
				break L
			default:
				fmt.Println("error, bad end of comp")
			}
		}
		pToT.putBuffer()
		pToT.delimitBuffering() // end
		pToT.put(arity)
		pToT.putBuffer()
	// min expects a function enclosed in brackets
	case min:
		pToT.put(min)
		expect(openParen)
		function()
		expect(closeParen)
	// rec expects 2 funtions enclosed in brackets
	case rec:
		pToT.put(rec)
		expect(openParen)
		function()
		expect(comma)
		function()
		expect(closeParen)
	// identifiers are unchanged
	case identifier:
		pToT.put(identifier)
		pToT.put(sToP.get())
	// this is the case of useless brackets
	case openParen:
		function()
		expect(closeParen)
	}
}

func expect(b byte) {
	// expect consumes a single value from the input stream
	// making sure it's value is the one given as the function parameter
	if sToP.get() != b {
		panic("expectation")
	}
}
