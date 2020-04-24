package main

func parse(s_to_p chan byte, p_to_t chan byte) {
	//stack := make([]byte)
	for {
		/*c := <- s_to_p
		if c == identifier {
			p_to_t <- <- s_to_p
			
		} else if c != newline {
			fmt.Println("error: unexpected symbol")
			p_to_t <- err
			return
		}*/
		c := <- s_to_p
		if c==end {
			break
		}
	}
	p_to_t <- end
}
