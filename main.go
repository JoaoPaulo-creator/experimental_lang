package main

import (
	"experimental/lexer"
	"experimental/parser"
	"fmt"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	data, _ := os.ReadFile("./test.fin")

	in := lexer.Tokenize(string(data))
	fmt.Printf("\n\n")

	litter.Dump(in)
	parser := parser.NewBTParser(in)
	prog, err := parser.ParseProgram()

	if err != nil {
		panic(err.Error())
	}

	litter.Dump(prog)
}
