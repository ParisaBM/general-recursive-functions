// this program takes a filename from the command line
// it opens it and outputs it to stdout

package main

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

// these are the signals the phases use to communicate with one another
const (
	constT     = iota // T short for token to seperate it from the go keyword
	suc        = iota
	proj       = iota
	comp       = iota
	min        = iota
	rec        = iota
	identifier = iota
	definition = iota
	equals     = iota
	newline    = iota
	err        = iota
	end        = iota
	// scanner only tokens
	openParen  = iota
	closeParen = iota
	comma      = iota
)

// these are the channels the phases use to communicate with one another
// f=file, s=scanner, p=parser, t=semantic, c=code_gen, e=exit
var fToS Stream
var sToP Stream
var pToT Stream
var tToC Stream
var cToE Stream

// the name list puts all the keyword identifiers in order
// it is the inverse to the nameTable defined in the scanner
// the scanner adds names to the nameList which will be used again in the code generation phase
var nameList []string
var nameListMutex sync.Mutex

func main() {
	// it is beyond the scope of this project to support multi-file codebases
	// there is an optional second arguement to debug a particular phase
	// if an integer is given as the fifth arguement, its respective phase
	// 0=scan, 1=parse etc. doesn't output to the next phase and is instead
	// output to stdout with annotations
	inputFileName := os.Args[1]
	file, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	// the output filename replaces the .grf with .ll
	// if the file extension isn't .grf we just add .ll
	name, extension, ok := strings.Cut(inputFileName, ".")
	var outputFileName string
	if ok && extension == "grf" {
		outputFileName = name + ".ll"
	} else {
		println("File name is supposed to end with .grf, not that this matters.")
		outputFileName = inputFileName + ".ll"
	}
	n := 3 // n is the phase the is being debbuged (if any)
	if len(os.Args) == 3 {
		n, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		// as there are 4 phases with output, n must be 0, 1, 2, or 3
		if n < 0 || n > 3 {
			panic("invalid debug phase")
		}
	}
	// next we initialize the channels
	fToS = fromFile(file) // closes the file
	sToP = newStream()
	pToT = newStream()
	tToC = newStream()
	cToE = newStream()
	// then we initialize the nameList
	nameList = make([]string, 0)
	// then we begin the phases
	if n >= 0 {
		go scan()
	}
	if n >= 1 {
		go parse()
	}
	if n >= 2 {
		go semantic()
	}
	if n == 3 {
		go codeGen(outputFileName)
		cToE.get()
	}
	if n != 3 {
		debug(n)
	}
}
