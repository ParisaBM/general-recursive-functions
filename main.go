//this program takes a filename from the command line
//it opens it and outputs it to stdout

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//this program will sometimes be run as its executable, and sometimes with go run
	//hence the filename is assumed to be the last arguement
	//this feature will be removed in the future
	//it is beyond the scope of this project to support multi-file codebases
	filename := os.Args[len(os.Args)-1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
