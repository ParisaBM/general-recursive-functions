package main

import "fmt"

func parse() {
	//the syntax of this language is highly regular
	//each line consists of code contains a definition
	//each iteration of this loop consumes one definition or a blank line, or the end of the file
L:
	for {
		switch s_to_p.get() {
		case identifier:
			//the general case is: identifer id equals function newline
			p_to_t.put(identifier)
			p_to_t.put(s_to_p.get())
			expect(equals)
			function()
			expect(newline)
			p_to_t.put(newline)
		case newline:
			p_to_t.put(newline)
		case end:
			p_to_t.put(end)
			break L
		default:
			//the errors that just say "error" are placeholders
			fmt.Println("error")
		}
	}
}

func function() {
	//consumes one function from the input stream
	//a function can be one of the primitives: suc, constant, proj
	//one of the operators: comp, min, rec
	//an identifer or a function enclosed in unnecessary brackets
	switch s_to_p.get() {
	//suc and constant are unchanged
	case suc:
		p_to_t.put(suc)
	case const_t:
		p_to_t.put(const_t)
		p_to_t.put(s_to_p.get())
	//the constant tags after proj are deleted from the stream
	//this is because they are redundant
	case proj:
		p_to_t.put(proj)
		expect(const_t)
		p_to_t.put(s_to_p.get())
		expect(const_t)
		p_to_t.put(s_to_p.get())
	//comp expects a comma seperated list of atleast one function in brackets
	case comp:
		p_to_t.put(comp)
		p_to_t.delimit_buffering() //begin
		expect(open_paren)
		p_to_t.delimit_buffering() //begin
		function()                 //this is the first function
		p_to_t.delimit_buffering() //end
		//a composition is emitted with its arity so the sematic analyzer can deduce
		//how many functions are being composed together
		var arity byte = 0
	L:
		for { //this loop consumes all addition functions
			switch s_to_p.get() {
			//if there's a comma there's an addition function
			case comma:
				function()
				arity++
			//if there's a close_paren, we've reached the end of the argument list
			case close_paren:
				break L
			default:
				fmt.Println("error, bad end of comp")
			}
		}
		p_to_t.put_buffer()
		p_to_t.delimit_buffering() //end
		p_to_t.put(arity)
		p_to_t.put_buffer()
	//min expects a function enclosed in brackets
	case min:
		p_to_t.put(min)
		expect(open_paren)
		function()
		expect(close_paren)
	//rec expects 2 funtions enclosed in brackets
	case rec:
		p_to_t.put(rec)
		expect(open_paren)
		function()
		expect(comma)
		function()
		expect(close_paren)
	//identifiers are unchanged
	case identifier:
		p_to_t.put(identifier)
		p_to_t.put(s_to_p.get())
	//this is the case of useless brackets
	case open_paren:
		function()
		expect(close_paren)
	}
}

func expect(b byte) {
	//expect consumes a single value from the input stream
	//making sure it's value is the one given as the function parameter
	if s_to_p.get() != b {
		panic("expectation")
	}
}
