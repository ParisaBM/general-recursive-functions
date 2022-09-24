package compile_test

import (
	"fmt"
	"main/compile"
	"os"
	"os/exec"
	"testing"
)

// TestCompile runs a series of tests to ensure each grf program can be compiled, ran, and gives the right outputs
// first it's compiled to llvm ir using the compile package we've written, then assembled with clang, then ran
// if any of these steps fail the test is a failure
// any files generated get cleaned up before the process finishes
func TestCompile(t *testing.T) {
	// each test contains:
	// the file name of the program WITHOUT the extension
	// a list of inputs for the program, each represented as a list of ints
	// a list of expected results, each as an integer
	tests := []struct {
		name    string
		inputs  [][]int
		outputs []int
	}{
		{"add1", [][]int{{0}, {1}, {6}, {7}}, []int{1, 2, 7, 8}},
		{"double", [][]int{{0}, {2}, {3}, {10}}, []int{0, 4, 6, 20}},
		{"sqrt", [][]int{{0}, {3}, {4}, {17}}, []int{0, 1, 2, 4}},
		{"lcm", [][]int{{1, 3}, {4, 6}, {5, 10}, {5, 6}}, []int{3, 12, 10, 30}},
	}
	for _, test := range tests {
		// this first part generates the executable as a.out
		inputFile, err := os.Open(test.name + ".grf")
		if err != nil {
			t.Fatal(err)
		}
		outputFile, err := os.Create(test.name + ".ll")
		if err != nil {
			t.Fatal(err)
		}
		compile.Compile(inputFile, outputFile, 3)
		err = exec.Command("clang", test.name+".ll").Run()
		if err != nil {
			t.Fatal(err)
		}
		for i := range test.inputs {
			// this inner loop runs the executable on each test case
			args := make([]string, 0)
			for _, x := range test.inputs[i] {
				args = append(args, fmt.Sprint(x))
			}
			out, err := exec.Command("./a.out", args...).Output()
			if err != nil {
				t.Fatal(err)
			}
			out = out[:len(out)-1]
			if string(out[:]) != fmt.Sprint(test.outputs[i]) {
				t.Fatalf("%s given %v returned %s, expected %d", test.name, args, string(out[:]), test.outputs[i])
			}
		}
		inputFile.Close()
		outputFile.Close()
		os.Remove(test.name + ".ll")
	}
	os.Remove("a.out")
}
