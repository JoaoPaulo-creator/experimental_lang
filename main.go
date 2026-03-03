package main

import (
	"experimental/lexer"
	"experimental/parser"
	"log"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	data, _ := os.ReadFile("./test.fin")

	in := lexer.Tokenize(string(data))

	parser := parser.NewParser(in)
	program, err := parser.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}

	litter.Dump(program)
}
