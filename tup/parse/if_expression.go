package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// condition = expression .

func Condition(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
	return Expression(tokens)
}

// else_block = "else" block .

func ElseBlock(tokens []tok.Token) (*ast.ElseBlock, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwElse {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	block, remainder, err := Block(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("block", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewElseBlock(block), remainder, nil
}

// if_expression = "if" condition block { "else" "if" condition block } [ else_block ] .

func IfExpression(tokens []tok.Token) (*ast.IfExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwIf {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	condition, remainder, err := Condition(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("expression", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	block, remainder, err := Block(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("block", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	conditions := []ast.Node{condition}
	blocks := []*ast.Block{block}
	hasElse := false

	for {
		remainder2 := skipTrivia(remainder)
		if peek(remainder2).Type != tok.TokKwElse {
			break
		}
		remainder2 = remainder2[1:]

		if peek(skipTrivia(remainder2)).Type != tok.TokKwIf {
			elseBlock, remainder3, err := ElseBlock(remainder)
			if err != nil {
				return nil, remainder3, err
			}
			blocks = append(blocks, elseBlock.Block)
			remainder = remainder3
			hasElse = true
			break
		}

		remainder2 = skipTrivia(remainder2)
		remainder2 = remainder2[1:]

		elseIfCondition, remainder3, err := Condition(remainder2)
		if err == ErrNoMatch {
			return nil, remainder3, errorExpecting("expression", remainder3)
		} else if err != nil {
			return nil, remainder3, err
		}

		elseIfBlock, remainder3, err := Block(remainder3)
		if err == ErrNoMatch {
			return nil, remainder3, errorExpecting("block", remainder3)
		} else if err != nil {
			return nil, remainder3, err
		}

		conditions = append(conditions, elseIfCondition)
		blocks = append(blocks, elseIfBlock)
		remainder = remainder3
	}

	return ast.NewIfExpression(conditions, blocks, hasElse), remainder, nil
}
