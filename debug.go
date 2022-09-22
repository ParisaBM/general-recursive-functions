package main

func debug(n int) {
	// debug outputs the result of the nth phase
	// stream is whichever streameam it's supposed to be listening to
	var stream Stream
	switch n {
	case 0:
		stream = sToP
	case 1:
		stream = pToT
	case 2:
		stream = tToC
	}
L:
	for {
		// there's a bit of code to handle compound tokens, but otherwise it just displays each token
		switch stream.get() {
		case constT:
			println("constant")
			println(stream.get())
		case suc:
			println("suc")
		case proj:
			println("proj")
			// in the scanner output there are constant markers that are removed by the parser
			if n == 1 || n == 2 {
				println(stream.get())
			}
			if n == 1 {
				println(stream.get())
			}
		case comp:
			println("comp")
			// in the parser and semantic analyzer, comp is followed by its arity
			if n == 1 || n == 2 {
				println(stream.get())
			}
		case min:
			println("min")
		case rec:
			println("rec")
		case identifier:
			println("identifier")
			println(stream.get())
		case definition:
			println("define")
			println(stream.get())
			// in the semantic analyzer, a definition is followed by its arity
			if n == 2 {
				println(stream.get())
			}
		case equals:
			println("equals")
		case openParen:
			println("openParen")
		case closeParen:
			println("closeParen")
		case comma:
			println("comma")
		case newline:
			// the output is grouped by line
			println("newline\n")
		case err:
			println("err")
		case end:
			println("end")
			break L
		}
	}
}
