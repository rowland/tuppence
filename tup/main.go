package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var input string
	var output string
	// use standard flag package; short flags duplicate the same variable.
	flag.StringVar(&input, "input", "", "Input file.")
	flag.StringVar(&input, "i", "", "Input file (shorthand).")
	flag.StringVar(&output, "output", "", "Output file.")
	flag.StringVar(&output, "o", "", "Output file (shorthand).")
	flag.Parse()

	if input != "" {
		fmt.Printf("Input file: '%s'\n", input)
		source, err := os.ReadFile(input)
		if err != nil {
			log.Fatal(err)
		}
		tokens, err := Tokenize(source, input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("tokens: %+v\n", tokens)
	}

	if output != "" {
		fmt.Printf("Output file: '%s'\n", output)
	}

	// Print any extra (positional) arguments.
	for i, arg := range flag.Args() {
		fmt.Printf("%d: %s\n", i, arg)
	}
}
