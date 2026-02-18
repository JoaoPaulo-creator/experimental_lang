package main

import (
	"experimental/lexer"
	"experimental/parser"
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("./test.fin")

	in := lexer.Tokenize(string(data))
	fmt.Printf("\n\n")
	parser := parser.NewPEGParser(in)
	prog, ok := parser.Program()
	fmt.Println("PEG ok:", ok, "funs:", len(prog.Func))
	if ok {
		fmt.Println("Parsed function:", prog.Func[0].Name, "stmts:", len(prog.Func[0].Body))
	}
}
