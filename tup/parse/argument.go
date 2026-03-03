package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// argument = ( expression | spread_argument ) .
// spread_argument = "..." expression .

func Argument(tokens []tok.Token) (arg *ast.Argument, remainder []tok.Token, err error) {
	// fmt.Println("Argument", tokens)
	remainder = skipComments(tokens)

	remainder, err = SpreadOp(remainder)
	spread := (err == nil)

	if expression, remainder, err := Expression(remainder); err == nil {
		return ast.NewArgument(expression, spread), remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// arguments = argument { "," argument } .

func Arguments(tokens []tok.Token) (args *ast.Arguments, remainder []tok.Token, err error) {
	// fmt.Println("Arguments", tokens)
	remainder = skipComments(tokens)

	var argsList []*ast.Argument
	var arg *ast.Argument
	for {
		arg, remainder, err = Argument(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, tokens, err
		}
		argsList = append(argsList, arg)
		remainder, err = Comma(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}
	}
	if len(argsList) > 0 {
		return ast.NewArguments(argsList), remainder, nil
	}
	return nil, tokens, ErrNoMatch
}

// labeled_arguments = labeled_argument { "," ( labeled_argument ) } .

// arguments_body = labeled_arguments
//                | arguments [ "," labeled_arguments ]

// function_arguments = ( arguments_body [ partial_application ]
// 	                    | "*"
// 	                    )
// 	                    [ "," ] .
