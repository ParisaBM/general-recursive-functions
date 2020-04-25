package main

import "fmt"

func parse() {
	for {
		switch <- s_to_p {
		case identifier:
			p_to_t <- <- s_to_p
			expect(equals)
			function()
			p_to_t <- equals
			expect(newline)
		case newline:
		case end:
			break
		default:
			fmt.Println("error")
		}
	}
	p_to_t <- end
}

func function() {
	switch <- s_to_p {
	case suc:
		p_to_t <- suc
	case constant:
		p_to_t <- constant
		p_to_t <- <- s_to_p
	case proj:
		p_to_t <- proj
		expect(constant)
		p_to_t <- <- s_to_p
		expect(constant)
		p_to_t <- <- s_to_p
	case comp:
		expect(open_paren)
		function()
		var arity byte = 0
		for {
			b := <- s_to_p
			switch b {
			case comma:
				function()
				arity++
			case close_paren:
				break
			default:
				fmt.Println("error, bad end of comp")
				fmt.Println(b)
			}
		}
		p_to_t <- comp
		p_to_t <- arity
	case min:
		expect(open_paren)
		function()
		expect(close_paren)
		p_to_t <- min
	case rec:
		expect(open_paren)
		function()
		expect(comma)
		function()
		expect(close_paren)
		p_to_t <- rec
	case identifier:
		p_to_t <- identifier
		p_to_t <- <- s_to_p
	/*case open_paren:
		function()
		expect(close_paren)*/
	}
}

func expect(b byte) {
	if <- s_to_p != b {
		fmt.Println("error")
	}
}
