package main

//import "fmt"
func parse() {} //On Hold
/*func parse() {
	//the syntax of this language is highly regular
	//each line consists of code contains a definition
	//each iteration of this loop consumes one definition or a blank line, or the end of the file
L:
	for {
		switch <-s_to_p {
		case identifier:
			//the general case is: identifer id equals function newline
			p_to_t <- identifier
			p_to_t <- <-s_to_p
			expect(equals)
			function()
			expect(newline)
			p_to_t <- equals
			p_to_t <- newline
		case newline:
			p_to_t <- newline
		case end:
			break L
		default:
			//the errors that just say "error" are placeholders
			fmt.Println("error")
		}
	}
	p_to_t <- end
}

func function() {
	//consumes one function from the input stream
	//a function can be one of the primitives: suc, constant, proj
	//one of the operators: comp, min, rec
	//an identifer or a function enclosed in unnecessary brackets
	switch <-s_to_p {
	//suc and constant are unchanged
	case suc:
		p_to_t <- suc
	case constant:
		p_to_t <- constant
		p_to_t <- <-s_to_p
	//the constant tags after proj are deleted from the stream
	//this is because they are redundant
	case proj:
		p_to_t <- proj
		expect(constant)
		p_to_t <- <-s_to_p
		expect(constant)
		p_to_t <- <-s_to_p
	//comp expects a comma seperated list of atleast one function in brackets
	case comp:
		p_to_t <- prefix_comp
		expect(open_paren)
		function() //this is the first function
		//a composition is emitted with its arity so the sematic analyzer can deduce
		//how many functions are being composed together
		var arity int8 = 0
	L:
		for { //this loop consumes all addition functions
			switch <-s_to_p {
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
		p_to_t <- comp
		p_to_t <- arity
	//min expects a function enclosed in brackets
	case min:
		p_to_t <- prefix_min
		expect(open_paren)
		function()
		expect(close_paren)
		p_to_t <- min
	//rec expects 2 funtions enclosed in brackets
	case rec:
		p_to_t <- prefix_rec
		expect(open_paren)
		function()
		expect(comma)
		function()
		expect(close_paren)
		p_to_t <- rec
	//identifiers are unchanged
	case identifier:
		p_to_t <- identifier
		p_to_t <- <-s_to_p
	//this is the case of useless brackets
	case open_paren:
		function()
		expect(close_paren)
	}
}

func expect(b int8) {
	//expect consumes a single value from the input stream
	//making sure it's value is the one given as the function parameter
	if <-s_to_p != b {
		fmt.Println("error")
	}
}*/
