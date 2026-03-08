package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// argument = ( expression | spread_argument ) .
// spread_argument = "..." expression .

func Argument(tokens []tok.Token) (arg *ast.Argument, remainder []tok.Token, err error) {
	// fmt.Println("Argument", tokens)

	var spread bool
	remainder, spread = SpreadOp(tokens)

	if expression, remainder, err := Expression(remainder); err == nil {
		return ast.NewArgument(expression, spread), remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// arguments = argument { "," argument } .

func Arguments(tokens []tok.Token) (args *ast.Arguments, remainder []tok.Token, err error) {
	// fmt.Println("Arguments", tok.Types(tokens))
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
		remainder2, found := Comma(remainder)
		if !found {
			break
		}
		if _, _, err = ArgumentLabel(remainder2); err == nil {
			break // argument label found, so we're done
		}
		if _, _, err = Argument(remainder2); err != nil {
			break
		}
		remainder = remainder2
	}
	if len(argsList) > 0 {
		return ast.NewArguments(argsList), remainder, nil
	}
	return nil, tokens, ErrNoMatch
}

// argument_label = identifier ":" .

func ArgumentLabel(tokens []tok.Token) (ident *ast.Identifier, remainder []tok.Token, err error) {
	if len(tokens) >= 2 && tokens[0].Type == tok.TokID && tokens[1].Type == tok.TokColon {
		return ast.NewIdentifier(tokens[0].Value(), tokens[0].File, tokens[0].Offset, tokens[0].Length), tokens[2:], nil
	}
	return nil, tokens, ErrNoMatch
}

// labeled_argument = ( identifier ":" argument ) .

func LabeledArgument(tokens []tok.Token) (arg *ast.LabeledArgument, remainder []tok.Token, err error) {
	// fmt.Println("LabeledArgument", tokens)

	identifier, remainder, err := ArgumentLabel(tokens)
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
	// fmt.Println("LabeledArguments", tok.Types(tokens))
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

		remainder2, found := Comma(remainder)
		if !found {
			break
		}
		if _, _, err = LabeledArgument(remainder2); err != nil {
			break
		}
		remainder = remainder2
	}
	if len(argsList) > 0 {
		return ast.NewLabeledArguments(argsList), remainder, nil
	}
	return nil, tokens, ErrNoMatch
}

// arguments_body = labeled_arguments
//                | arguments [ "," labeled_arguments ]

func ArgumentsBody(tokens []tok.Token) (args *ast.Arguments, labeledArgs *ast.LabeledArguments, remainder []tok.Token, err error) {
	// fmt.Println("ArgumentsBody", tok.Types(tokens))

	labeledArgs, remainder, err = LabeledArguments(tokens)
	if err == nil {
		return nil, labeledArgs, remainder, nil
	} else if err != ErrNoMatch {
		return nil, nil, remainder, err
	}

	args, remainder, err = Arguments(remainder)
	if err != nil {
		return nil, nil, remainder, err
	}

	remainder2, found := Comma(remainder)
	if !found {
		return args, nil, remainder, nil
	}

	labeledArgs, remainder2, err = LabeledArguments(remainder2)
	if err == nil {
		return args, labeledArgs, remainder2, nil
	}

	return args, nil, remainder, nil
}
