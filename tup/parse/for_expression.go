package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// for_block = "{" { statement } [ expression ] "}" .

func ForBlock(tokens []tok.Token) (*ast.ForBlock, []tok.Token, error) {
	remainder, found := OpenBrace(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	statements, remainder, err := Statements(remainder)
	if err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var expression ast.Expression
	if len(statements) > 0 {
		if finalExpression, ok := statements[len(statements)-1].(ast.Expression); ok {
			expression = finalExpression
			statements = statements[:len(statements)-1]
		}
	}

	if remainder, found = CloseBrace(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBrace, remainder)
	}

	return ast.NewForBlock(statements, expression), remainder, nil
}

// initializer = assignment .

func Initializer(tokens []tok.Token) (*ast.Initializer, []tok.Token, error) {
	assignment, remainder, err := Assignment(tokens)
	if err != nil {
		return nil, remainder, err
	}
	return ast.NewInitializer(assignment), remainder, nil
}

// iterable = expression .

func Iterable(tokens []tok.Token) (*ast.Iterable, []tok.Token, error) {
	expression, remainder, err := Expression(tokens)
	if err != nil {
		return nil, remainder, err
	}
	return ast.NewIterable(expression), remainder, nil
}

// step_expression = expression .

func StepExpression(tokens []tok.Token) (*ast.StepExpression, []tok.Token, error) {
	expression, remainder, err := Expression(tokens)
	if err != nil {
		return nil, remainder, err
	}
	return ast.NewStepExpression(expression), remainder, nil
}

// for_header = initializer [ ";" condition [ ";" step_expression ] ] .

func ForHeader(tokens []tok.Token) (*ast.ForHeader, []tok.Token, error) {
	initializer, remainder, err := Initializer(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder2, found := skipComments(remainder), false
	if remainder2, found = SemiColon(remainder); !found {
		return ast.NewForHeader(initializer, nil, nil), remainder, nil
	}
	remainder = remainder2

	condition, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("condition", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder2, found = SemiColon(remainder); !found {
		return ast.NewForHeader(initializer, condition, nil), remainder, nil
	}
	remainder = remainder2

	stepExpression, remainder, err := StepExpression(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("step expression", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewForHeader(initializer, condition, stepExpression), remainder, nil
}

// for_expression = "for" [ for_header | for_in_header ] for_block .
//
// This first pass implements the plain for forms using for_header and defers
// for_in_header to a later change.

func ForExpression(tokens []tok.Token) (*ast.ForExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwFor {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	if forBlock, remainder2, err := ForBlock(remainder); err == nil {
		return ast.NewForExpression(nil, forBlock), remainder2, nil
	} else if err != ErrNoMatch {
		return nil, remainder2, err
	}

	header, remainder, err := ForHeader(remainder)
	if err != nil {
		return nil, remainder, err
	}

	forBlock, remainder, err := ForBlock(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewForExpression(header, forBlock), remainder, nil
}
