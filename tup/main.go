package main

import (
	"fmt"
	"os"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/parse"
	"github.com/rowland/tuppence/tup/source"
)

func main() {
	paths := os.Args[1:]
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "usage: tup <file> [file...]")
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range paths {
		if err := checkFile(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func checkFile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	src := source.NewSource(contents, path)
	mod := ast.NewModule(path)
	_, err = parse.Module(src, mod)
	return err
}
