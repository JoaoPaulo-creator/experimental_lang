package main

import (
	"experimental/lexer"
	"experimental/parser"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	data, _ := os.ReadFile("./test.fin")

	in := lexer.Tokenize(string(data))
	parser := parser.Parse(in)
	litter.Dump(parser)
}
