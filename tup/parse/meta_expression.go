package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// meta_expression = "$" labeled_tuple .

func MetaExpression(tokens []tok.Token) (expr *ast.MetaExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = Dollar(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var labeledArgs *ast.LabeledArguments
	if labeledArgs, remainder, err = LabeledArguments(remainder); err == ErrNoMatch {
		return nil, remainder, errorExpecting("labeled argument", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	remainder, _ = Comma(remainder)

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	keyValues := make(map[string]ast.Node, len(labeledArgs.Args))
	for _, arg := range labeledArgs.Args {
		keyValues[arg.Identifier.Name] = arg.Argument
	}

	return ast.NewMetaExpression(keyValues), remainder, nil
}
