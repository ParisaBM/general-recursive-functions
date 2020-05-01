package main

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
			panic("arity error")
		}
		return ar0
	} else if ar0.known && !ar1.known {
		if ar0.arity < ar1.arity {
			panic("arity error")
		}
		return ar0
	} else if !ar0.known && ar1.known {
		if ar1.arity < ar0.arity {
			panic("arity error")
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
			panic("arity error")
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

var arity_table map[byte]Arity
func semantic() {
	 arity_table = make(map[byte]Arity)
L:
	for {
		switch p_to_t.get() {
		case identifier:
			t_to_r.put(identifier)
			id := p_to_t.get()
			t_to_r.put(id)
			t_to_r.delimit_buffering() //begin
			ar := sem_function()
			t_to_r.delimit_buffering() //end
			t_to_r.put(byte(ar.arity))
			t_to_r.put_buffer()
			arity_table[id] = Arity{ar.arity, true}
		case end:
			break L
		}
	}
	_, ok := arity_table[0]
	if !ok {
		panic("no main function")
	}
}

//reccursively semantically analyzes a function, returns the arity
func sem_function() Arity {
	switch p_to_t.get() {
	case constant:
		t_to_r.put(constant)
		t_to_r.put(p_to_t.get())
		return Arity{0, false}
	case suc:
		t_to_r.put(suc)
		return Arity{1, true}
	case proj:
		t_to_r.put(proj)
		//recieving proj n m
		n := p_to_t.get()
		m := p_to_t.get()
		t_to_r.put(n) //m is not needed by the next phase
		return Arity{int8(m), true}
	case comp:
		t_to_r.put(comp)
		//n is the arity of the top level function
		//for example in comp(F0, F1, ...Fn), n is the arity of F0
		n := p_to_t.get()
		t_to_r.put(n)
		//the arity of the whole thing is the intersection of F1, ... Fn
		ar := Arity{0, false}		
		for i := byte(0); i < n; i++ {
			ar = merge(ar, sem_function())
		}
		//the result of this merge isn't used in any way
		//its just to make sure its correct
		merge(Arity{int8(n), true}, sem_function())
		return ar
	case min:
		t_to_r.put(min)
		return add(sem_function(), -1)
	case rec:
		t_to_r.put(rec)
		return merge(add(sem_function(), 1), add(sem_function(), -1))
	case identifier:
		t_to_r.put(identifier)
		id := p_to_t.get()
		t_to_r.put(id)
		return arity_table[id]
	default:
		panic("unable to match")
	}
}
