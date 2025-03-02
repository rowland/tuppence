package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rowland/tuppence/tup/tok"
	"github.com/spf13/pflag"
)

func main() {
	var input string
	var output string
	pflag.StringVarP(&input, "input", "i", "", "Input file")
	pflag.StringVarP(&output, "output", "o", "", "Output file")
	pflag.Parse()

	if input != "" {
		fmt.Printf("Input file: '%s'\n", input)
		source, err := os.ReadFile(input)
		if err != nil {
			log.Fatal(err)
		}
		tokens, err := tok.Tokenize(source, input)
		if err != nil {
			log.Fatal(err)
		}
		for _, token := range tokens {
			fmt.Printf("%s %s\n", token.Type, token.Value)
		}
	}

	if output != "" {
		fmt.Printf("Output file: '%s'\n", output)
	}

	// Print any extra (positional) arguments.
	for i, arg := range pflag.Args() {
		fmt.Printf("%d: %s\n", i, arg)
	}
}
