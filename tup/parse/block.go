package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// block = "{" block_body "}" .

func Block(tokens []tok.Token) (expr *ast.Block, remainder []tok.Token, err error) {
	// fmt.Println("Block", tok.Types(tokens))
	var found bool
	if remainder, found = OpenBrace(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var body *ast.BlockBody
	if body, remainder, err = BlockBody(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseBrace(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBrace, remainder)
	}

	return ast.NewBlock(body), remainder, nil
}

// block_body = { statement } expression .

func BlockBody(tokens []tok.Token) (expr *ast.BlockBody, remainder []tok.Token, err error) {
	// fmt.Println("BlockBody", tok.Types(tokens))
	var statements []ast.Statement
	if statements, remainder, err = Statements(tokens); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	if len(statements) == 0 {
		return nil, remainder, errorExpecting("expression", remainder)
	}

	expression, ok := statements[len(statements)-1].(ast.Expression)
	if !ok {
		return nil, remainder, errorExpecting("expression", remainder)
	}

	return ast.NewBlockBody(statements[:len(statements)-1], expression), remainder, nil
}
