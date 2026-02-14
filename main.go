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
	litter.Dump(in)
	fmt.Printf("\n\n")
	parser := parser.Parse(in)
	litter.Dump(parser)
}
