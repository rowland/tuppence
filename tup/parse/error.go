package parse

import (
	"errors"
	"fmt"

	"github.com/rowland/tuppence/tup/tok"
)

var ErrNoMatch = errors.New("no match")

type Error struct {
	Filename  string
	Expecting string
	Got       string
	Line      int
	Column    int
	// Fatal     bool
	SubErrors []error
}

func (err *Error) Error() string {
	loc := fmt.Sprintf("%d:%d", err.Line+1, err.Column+1)
	if err.Filename != "" {
		loc = fmt.Sprintf("--> %s:%d:%d", err.Filename, err.Line+1, err.Column+1)
	}
	msg := fmt.Sprintf("error: expecting %#v, got %#v", err.Expecting, err.Got)
	if len(err.SubErrors) > 0 {
		return fmt.Sprintf("%s: %v\n%s", msg, err.SubErrors, loc)
	}
	return fmt.Sprintf("%s\n%s", msg, loc)
}

func errorExpecting(expecting string, tokens []tok.Token) *Error {
	filename, got, line, column := errorGot(tokens)
	return &Error{
		Filename:  filename,
		Expecting: expecting,
		Got:       got,
		Line:      line,
		Column:    column,
		// Fatal:     true,
		SubErrors: []error{},
	}
}

func errorExpectingOneOf(expecting string, tokens []tok.Token, errors []error) *Error {
	filename, got, line, column := errorGot(tokens)
	return &Error{
		Filename:  filename,
		Expecting: expecting,
		Got:       got,
		Line:      line,
		Column:    column,
		// Fatal:     true,
		SubErrors: errors,
	}
}

func errorExpectingTokenType(tokenType tok.TokenType, tokens []tok.Token) *Error {
	filename, got, line, column := errorGot(tokens)
	return &Error{
		Filename:  filename,
		Expecting: tok.TokenTypes[tokenType],
		Got:       got,
		Line:      line,
		Column:    column,
		// Fatal:     true,
		SubErrors: []error{},
	}
}

func errorNotExpecting(tokens []tok.Token) *Error {
	filename, got, line, column := errorGot(tokens)
	return &Error{
		Filename:  filename,
		Expecting: "not " + got,
		Got:       got,
		Line:      line,
		Column:    column,
	}
}

func errorGot(tokens []tok.Token) (filename string, got string, line int, column int) {
	tokens = skipTrivia(tokens)
	if len(tokens) > 0 {
		if tokens[0].File != nil {
			filename = tokens[0].File.Filename
		}
		got = tokens[0].Value()
		line = tokens[0].Line()
		column = tokens[0].Column()
	}
	return filename, got, line, column
}
