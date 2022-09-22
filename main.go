// this program takes a filename from the command line
// it opens it and outputs it to stdout

package main

import (
	"main/compile"
	"os"
	"strconv"
	"strings"
)

func main() {
	// it is beyond the scope of this project to support multi-file codebases
	// there is an optional second arguement to debug a particular phase
	// if an integer is given as the fifth arguement, its respective phase
	// 0=scan, 1=parse etc. doesn't output to the next phase and is instead
	// output to stdout with annotations
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
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
	debugStage := 3 // n is the phase the is being debbuged (if any)
	if len(os.Args) == 3 {
		debugStage, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		// as there are 4 phases with output, n must be 0, 1, 2, or 3
		if debugStage < 0 || debugStage > 3 {
			panic("invalid debug phase")
		}
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	compile.Compile(inputFile, outputFile, debugStage)
}
