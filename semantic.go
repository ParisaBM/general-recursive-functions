package main

//On Hold
func semantic(){}

/*import "fmt"

//the semantic phase must check that:
//there is a main function
//every function is defined before it is used
//no function is defined twice
//for every proj n m, we have n<m
//the arity rules are enforced
//every function/expression has either a known arity, or unknown arity
//if it unknown, there will be a (possibly 0) lower bound
//as such every arity can be thought of a set containing either 1, or infinite possibilies
//to merge 2 arities is to take the intersection of those possibilities
//2 arities are comptible if their intersection is non-empty
//an error is generated if incompatible arities must be merged
//these rules are given by:

type Arity struct {
	//Arity expresses what is know about the type of a function
	//arity is the number that is either the lower bound or exact number
	arity int8
	known bool
}

func max(x, y int8) int8 {
	//used to merge two unknown Arities
	if x > y {
		return x
	}
	return y
}

func merge(ar0, ar1 Arity) Arity {
	//sometimes calls to merge don't actually use the returned value
	//instead it's just to check compatibility
	if ar0.known && ar1.known {
		if ar0.arity != ar1.arity {
			fmt.Println("arity error")
		}
		return ar0
	} else if ar0.known && !ar1.known {
		if ar0.arity < ar1.arity {
			fmt.Println("arity error")
		}
		return ar0
	} else if !ar0.known && ar1.known {
		if ar1.arity < ar0.arity {
			fmt.Println("arity error")
		}
		return ar1
	} else {
		return Arity{max(ar0.arity, ar1.arity), false}
	}
}

func add(ar Arity, n int8) Arity {
	//adds n to an Arity, handles edge cases
	//n might be negative
	ar.arity += n
	if ar.arity < 0 {
		if ar.known {
			//this means the function has certainly negative arity
			fmt.Println("arity error")
		} else {
			ar.arity = 0
		}
	}
	return ar
}

//the arity of suc is 1
//the arity of a constant function is atleast 0
//this is because the programmer does not need to specify the arity of constants
//the arity of proj n m is m
//the arity of min(F) is the arity of F-1, this means F can't be known to be niladic
//the arity of comp(F0, F1, ...Fn) is the merge of F1 through Fn, where F0 must be compatible with n
//the arity of rec(F, G) is arity of F+1, the arity of F+2 must be compatible with the arity of G

func semantic() {
	//during the evaluation of a function's definition, all the information is kept
	//this evaluation uses the arity_stack
	arity_stack := make([]Arity, 0)
	//once it is finished evaluating, if the arity is still unknow it is assumed to be the minimum
	//the arity_table maps each identifier number to its arity, which requires ONLY A BYTE
	arity_table := make(map[int8]int8)
	//id is the numeric identifier of the function being defined
	//it is -1 before an assignment is made
	var id int8 = -1
	//the previous phases were all able to emit output as it was generated
	//the semantic phase must add arity annotations to the start of each line
	//thus the output must be buffered until then
	buffer := make([]int8, 0)
L:
	for {
		switch <-p_to_t {
		case identifier:
			//how to handle an identifier depends whether we're in an expression or not
			//we're in an expression if the stack is non-empty
			if id == -1 {
				id = <- p_to_t
				//we check id hasn't already been defined
				_, ok := arity_table[id]
				if ok {
					fmt.Println("double definition")
				}
				t_to_r <- identifier
				t_to_r <- id
			} else {
				n, ok := arity_table[<-p_to_t]
				if !ok {
					fmt.Println("unknown identifier")
				}
				arity_stack = append(arity_stack, Arity{n, true})
			}
		case end:
			break L
		case constant:
			arity_stack = append(arity_stack, Arity{0, false})
			buffer = append(buffer, constant, <- p_to_t)
		case suc:
			arity_stack = append(arity_stack, Arity{1, true})
			buffer = append(buffer, suc)
		case proj:
			n := <-p_to_t
			m := <-p_to_t
			if n >= m {
				fmt.Println("bad arity")
			}
			arity_stack = append(arity_stack, Arity{m, true})
			//m is not needed in the next phase
			buffer = append(buffer, proj, n)
		case comp:
			n := <-p_to_t
			for i := int8(0); i < n-1; i++ {
				//this loops merges the top 2 items n-1 times
				arity_stack = append(arity_stack[:len(arity_stack)-2],
					merge(arity_stack[len(arity_stack)-1], arity_stack[len(arity_stack)-2]))
			}
			//then we make sure the second from top arity is compatible with n, then delete it
			merge(arity_stack[len(arity_stack)-2], Arity{n, true})
			arity_stack = append(arity_stack[:len(arity_stack)-2], arity_stack[len(arity_stack)-1])
			//note how comp, but neither of the other postfix operators are kept in the output stream
			buffer = append(buffer, comp)
		case min:
			arity_stack[len(arity_stack)-1] = add(arity_stack[len(arity_stack)-1], -1)
		case rec:
			//this part is rather tricky, it must be exactly one more than the first operand, and one less than the second
			arity_stack = append(arity_stack[:len(arity_stack)-2],
				merge(add(arity_stack[len(arity_stack)-2], 1), add(arity_stack[len(arity_stack)-1], -1)))
		case equals:
			//pop the arity_stack and assign it to id
			arity_table[id] = arity_stack[len(arity_stack)-1].arity
			t_to_r <- arity_table[id]
			for i := buffer {
				t_to_r <- buffer[i]
			}
			arity_stack = make([]int8, 0)
			buffer = make([]int8, 0)
			//reset the id variable
			id = -1
		//these tokens aren't meant for the semantic analyzer
		//they get just get passed along to the representation phase
		case prefix_comp:
			buffer = append(buffer, prefix_comp)
		case prefix_min:
			buffer = append(buffer, prefix_min)
		case prefix_rec:
			buffer = append(buffer, prefix_rec)
		}
	}
	_, ok := arity_table[0]
	if !ok {
		fmt.Println("no main function")
	}
	for n, m := range arity_table {
		fmt.Println(n, m)
	}
	t_to_r <- 0
}*/
