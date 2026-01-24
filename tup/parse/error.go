package parse

import (
	"fmt"

	"github.com/rowland/tuppence/tup/tok"
)

type Error struct {
	Expecting string
	Got       string
	Line      int
	Column    int
	// Fatal     bool
	SubErrors []error
}

func (err *Error) Error() string {
	return fmt.Sprintf(`expecting %#v, got %#v at row %d, col %d`, err.Expecting, err.Got, err.Line+1, err.Column+1)
}

func errorExpecting(expecting string, tokens []tok.Token) *Error {
	got, line, column := errorGot(tokens)
	return &Error{
		Expecting: expecting,
		Got:       got,
		Line:      line,
		Column:    column,
		// Fatal:     true,
		SubErrors: []error{},
	}
}

func errorExpectingOneOf(expecting string, tokens []tok.Token, errors []error) *Error {
	got, line, column := errorGot(tokens)
	return &Error{
		Expecting: expecting,
		Got:       got,
		Line:      line,
		Column:    column,
		// Fatal:     true,
		SubErrors: errors,
	}
}

func errorGot(tokens []tok.Token) (got string, line int, column int) {
	if len(tokens) > 0 {
		got = tokens[0].Value()
		line = tokens[0].Line()
		column = tokens[0].Column()
	} else {
		got = "EOF"
	}
	return got, line, column
}
