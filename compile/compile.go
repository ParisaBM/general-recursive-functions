package compile

import (
	"os"
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
// f=file, s=scanner, p=parser, t=semantic, c=code_gen
var fToS Stream
var sToP Stream
var pToT Stream
var tToC Stream

// the name list puts all the keyword identifiers in order
// it is the inverse to the nameTable defined in the scanner
// the scanner adds names to the nameList which will be used again in the code generation phase
var nameList []string
var nameListMutex sync.Mutex

// com
func Compile(inputFile, outputFile *os.File, debugStage int) {
	// next we initialize the channels
	fToS = fromFile(inputFile) // closes the file
	sToP = newStream()
	pToT = newStream()
	tToC = newStream()
	// then we initialize the nameList
	nameList = make([]string, 0)
	// then we begin the phases
	if debugStage >= 0 {
		go scan()
	}
	if debugStage >= 1 {
		go parse()
	}
	if debugStage >= 2 {
		go semantic()
	}
	if debugStage == 3 {
		codeGen(outputFile)
	} else {
		debug(debugStage)
	}
}
