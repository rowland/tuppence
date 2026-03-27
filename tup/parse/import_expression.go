package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// import_expression = "import" "(" string_literal ")" .

func ImportExpression(tokens []tok.Token) (expr *ast.ImportExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	if peek(remainder).Type != tok.TokKwImport {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var path *ast.StringLiteral
	if path, remainder, err = StringLiteral(remainder); err == ErrNoMatch {
		return nil, remainder, errorExpecting("string literal", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewImportExpression(path), remainder, nil
}
