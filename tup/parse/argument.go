package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// argument = ( expression | spread_argument ) .
// spread_argument = "..." expression .

func Argument(tokens []tok.Token) (arg *ast.Argument, remainder []tok.Token, err error) {
	// fmt.Println("Argument", tokens)

	remainder, err = SpreadOp(tokens)
	spread := (err == nil)

	if expression, remainder, err := Expression(remainder); err == nil {
		return ast.NewArgument(expression, spread), remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// arguments = argument { "," argument } .

func Arguments(tokens []tok.Token) (args *ast.Arguments, remainder []tok.Token, err error) {
	// fmt.Println("Arguments", tokens)
	remainder = tokens

	var argsList []*ast.Argument
	for {
		var arg *ast.Argument
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

// labeled_argument = ( identifier ":" argument ) .

func LabeledArgument(tokens []tok.Token) (arg *ast.LabeledArgument, remainder []tok.Token, err error) {
	// fmt.Println("LabeledArgument", tokens)

	identifier, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder, err = Colon(remainder)
	if err != nil {
		return nil, remainder, err
	}

	argument, remainder, err := Argument(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledArgument(identifier, argument), remainder, nil
}

// labeled_arguments = labeled_argument { "," ( labeled_argument ) } .

func LabeledArguments(tokens []tok.Token) (args *ast.LabeledArguments, remainder []tok.Token, err error) {
	// fmt.Println("LabeledArguments", tokens)
	remainder = tokens

	var argsList []*ast.LabeledArgument
	for {
		var arg *ast.LabeledArgument
		arg, remainder, err = LabeledArgument(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
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
		return ast.NewLabeledArguments(argsList), remainder, nil
	}
	return nil, tokens, ErrNoMatch
}
