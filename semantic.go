package main

// the semantic phase must check that:
// there is a main function
// every function is defined before it is used
// no function is defined twice
// for every proj n m, we have n<m
// the arity rules are enforced
// every function/expression has either a known arity, or unknown arity
// if it unknown, there will be a (possibly 0) lower bound
// as such every arity can be thought of a set containing either 1, or infinite possibilies
// to merge 2 arities is to take the intersection of those possibilities
// 2 arities are comptible if their intersection is non-empty
// an error is generated if incompatible arities must be merged
// these rules are given by:

type Arity struct {
	// Arity expresses what is know about the type of a function
	// arity is the number that is either the lower bound or exact number
	arity int8
	known bool
}

func max(x, y int8) int8 {
	// used to merge two unknown Arities
	if x > y {
		return x
	}
	return y
}

func merge(ar0, ar1 Arity) Arity {
	// sometimes calls to merge don't actually use the returned value
	// instead it's just to check compatibility
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

func addArity(ar Arity, n int8) Arity {
	// adds n to an Arity, handles edge cases
	// n might be negative
	ar.arity += n
	if ar.arity < 0 {
		if ar.known {
			// this means the function has certainly negative arity
			panic("arity error")
		} else {
			ar.arity = 0
		}
	}
	return ar
}

// the arity of suc is 1
// the arity of a constant function is atleast 0
// this is because the programmer does not need to specify the arity of constants
// the arity of proj n m is m
// the arity of min(F) is the arity of F-1, this means F can't be known to be niladic
// the arity of comp(F0, F1, ...Fn) is the merge of F1 through Fn, where F0 must be compatible with n
// the arity of rec(F, G) is arity of F+1, the arity of F+2 must be compatible with the arity of G

var arityTable map[byte]Arity

func semantic() {
	arityTable = make(map[byte]Arity)
L:
	for {
		switch pToT.get() {
		case identifier:
			tToC.put(identifier)
			id := pToT.get()
			tToC.put(id)
			tToC.delimitBuffering() // begin
			ar := semFunction()
			tToC.delimitBuffering() // end
			tToC.put(byte(ar.arity))
			tToC.putBuffer()
			arityTable[id] = Arity{ar.arity, true}
		case end:
			tToC.put(end)
			break L
		}
	}
	_, ok := arityTable[0]
	if !ok {
		panic("no main function")
	}
}

// reccursively semantically analyzes a function, returns the arity
func semFunction() Arity {
	switch pToT.get() {
	case constT:
		tToC.put(constT)
		tToC.put(pToT.get())
		return Arity{0, false}
	case suc:
		tToC.put(suc)
		return Arity{1, true}
	case proj:
		tToC.put(proj)
		// recieving proj n m
		n := pToT.get()
		m := pToT.get()
		tToC.put(n) // m is not needed by the next phase
		return Arity{int8(m), true}
	case comp:
		tToC.put(comp)
		// n is the arity of the top level function
		// for example in comp(F0, F1, ...Fn), n is the arity of F0
		n := pToT.get()
		tToC.put(n)
		// the arity of the whole thing is the intersection of F1, ... Fn
		ar := Arity{0, false}
		for i := byte(0); i < n; i++ {
			ar = merge(ar, semFunction())
		}
		// the result of this merge isn't used in any way
		// its just to make sure its correct
		merge(Arity{int8(n), true}, semFunction())
		return ar
	case min:
		tToC.put(min)
		return addArity(semFunction(), -1)
	case rec:
		tToC.put(rec)
		return merge(addArity(semFunction(), 1), addArity(semFunction(), -1))
	case identifier:
		tToC.put(identifier)
		id := pToT.get()
		tToC.put(id)
		return arityTable[id]
	default:
		panic("unable to match")
	}
}
