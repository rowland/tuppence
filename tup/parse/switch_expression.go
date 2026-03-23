package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// switch_expression = "switch" expression "{" switch_case { switch_case } [ switch_else_block ] "}" .

func SwitchExpression(tokens []tok.Token) (*ast.SwitchExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwSwitch {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	subject, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("expression", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = OpenBrace(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenBrace, remainder)
	}

	remainder = skipTrivia(remainder)

	firstCase, remainder, err := SwitchCase(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("switch case", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	cases := []*ast.SwitchCase{firstCase}
	remainder = skipTrivia(remainder)

	for {
		if peek(skipTrivia(remainder)).Type == tok.TokKwElse || peek(skipTrivia(remainder)).Type == tok.TokCloseBrace {
			break
		}

		nextCase, remainder2, err := SwitchCase(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder2, err
		}

		cases = append(cases, nextCase)
		remainder = skipTrivia(remainder2)
	}

	var elseBlock *ast.FunctionBlock
	if elseBlock, remainder, err = SwitchElseBlock(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseBrace(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBrace, remainder)
	}

	return ast.NewSwitchExpression(subject, cases, elseBlock), remainder, nil
}

// switch_case = match_condition function_block .

func SwitchCase(tokens []tok.Token) (*ast.SwitchCase, []tok.Token, error) {
	condition, remainder, err := MatchCondition(tokens)
	if err != nil {
		return nil, remainder, err
	}

	body, remainder, err := FunctionBlock(skipTrivia(remainder))
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("function block", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewSwitchCase(condition, body), remainder, nil
}

// switch_else_block = "else" function_block .

func SwitchElseBlock(tokens []tok.Token) (*ast.FunctionBlock, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwElse {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	body, remainder, err := FunctionBlock(skipTrivia(remainder))
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("function block", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return body, remainder, nil
}

// match_condition = list_match | pattern .

func MatchCondition(tokens []tok.Token) (ast.MatchCondition, []tok.Token, error) {
	pattern, remainder, err := Pattern(tokens)
	if err != nil {
		return nil, remainder, err
	}

	matchElement, ok := pattern.(ast.MatchElement)
	if !ok {
		return pattern, remainder, nil
	}

	remainder2, found := Comma(remainder)
	if !found {
		return pattern, remainder, nil
	}

	return listMatchTail(matchElement, remainder2)
}

// list_match = match_element "," match_element { "," match_element } .

func ListMatch(tokens []tok.Token) (*ast.ListMatch, []tok.Token, error) {
	first, remainder, err := MatchElement(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder2, found := Comma(remainder)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	return listMatchTail(first, remainder2)
}

func listMatchTail(first ast.MatchElement, tokens []tok.Token) (*ast.ListMatch, []tok.Token, error) {
	remainder := tokens

	second, remainder, err := MatchElement(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("match element", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	elements := []ast.MatchElement{first, second}

	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		next, remainder3, err := MatchElement(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("match element", remainder2)
		} else if err != nil {
			return nil, remainder3, err
		}

		elements = append(elements, next)
		remainder = remainder3
	}

	return ast.NewListMatch(elements), remainder, nil
}

// match_element = constant | range | inferred_error_type | type_reference .

func MatchElement(tokens []tok.Token) (ast.MatchElement, []tok.Token, error) {
	if rangeExpr, remainder, err := Range(tokens); err == nil {
		return rangeExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if inferredErrorType, remainder, err := InferredErrorType(tokens); err == nil {
		return inferredErrorType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if constant, remainder, err := Constant(tokens); err == nil {
		return constant, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// pattern = wildcard_pattern | pattern_match | match_element .

func Pattern(tokens []tok.Token) (ast.Pattern, []tok.Token, error) {
	if wildcard, remainder, err := WildcardPattern(tokens); err == nil {
		return wildcard, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if patternMatch, remainder, err := PatternMatch(tokens); err == nil {
		return patternMatch, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if matchElement, remainder, err := MatchElement(tokens); err == nil {
		switch value := matchElement.(type) {
		case *ast.Constant:
			return value, remainder, nil
		case *ast.Range:
			return value, remainder, nil
		case *ast.InferredErrorType:
			return value, remainder, nil
		case *ast.TypeReference:
			return value, remainder, nil
		default:
			return nil, remainder, errorExpecting("pattern", remainder)
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// pattern_match = typed_pattern | structured_match .

func PatternMatch(tokens []tok.Token) (ast.Pattern, []tok.Token, error) {
	if typedPattern, remainder, err := TypedPattern(tokens); err == nil {
		return typedPattern, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if structuredPattern, remainder, err := StructuredMatch(tokens); err == nil {
		return structuredPattern, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// wildcard_pattern = "_" .

func WildcardPattern(tokens []tok.Token) (*ast.WildcardPattern, []tok.Token, error) {
	identifier, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	if identifier.Name != "_" {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewWildcardPattern(identifier), remainder, nil
}

// typed_pattern = type_reference structured_match .

func TypedPattern(tokens []tok.Token) (*ast.TypedPattern, []tok.Token, error) {
	typeReference, remainder, err := TypeReference(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder = skipTrivia(remainder)
	switch peek(remainder).Type {
	case tok.TokOpenParen, tok.TokOpenBracket:
	default:
		return nil, tokens, ErrNoMatch
	}

	pattern, remainder, err := StructuredMatch(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypedPattern(typeReference, pattern), remainder, nil
}

// structured_match = labeled_pattern | tuple_pattern | array_pattern .

func StructuredMatch(tokens []tok.Token) (ast.Pattern, []tok.Token, error) {
	if labeledPattern, remainder, err := LabeledPattern(tokens); err == nil {
		return labeledPattern, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if tuplePattern, remainder, err := TuplePattern(tokens); err == nil {
		return tuplePattern, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if arrayPattern, remainder, err := ArrayPattern(tokens); err == nil {
		return arrayPattern, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// labeled_pattern = "(" identifier ":" pattern { "," identifier ":" pattern } ")" .

func LabeledPattern(tokens []tok.Token) (*ast.LabeledPattern, []tok.Token, error) {
	remainder, found := OpenParen(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	member, remainder, err := labeledPatternMember(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	members := []*ast.LabeledPatternMember{member}
	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		member, remainder, err = labeledPatternMember(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("labeled pattern member", remainder2)
		} else if err != nil {
			return nil, remainder, err
		}
		members = append(members, member)
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewLabeledPattern(members), remainder, nil
}

func labeledPatternMember(tokens []tok.Token) (*ast.LabeledPatternMember, []tok.Token, error) {
	label, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder2, found := Colon(remainder)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	pattern, remainder, err := Pattern(remainder2)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledPatternMember(label, pattern), remainder, nil
}

// tuple_pattern = "(" pattern { "," pattern } ")" .

func TuplePattern(tokens []tok.Token) (*ast.TuplePattern, []tok.Token, error) {
	remainder, found := OpenParen(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	first, remainder, err := Pattern(remainder)
	if err != nil {
		return nil, remainder, err
	}

	elements := []ast.Pattern{first}
	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		next, remainder3, err := Pattern(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("pattern", remainder2)
		} else if err != nil {
			return nil, remainder3, err
		}

		elements = append(elements, next)
		remainder = remainder3
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewTuplePattern(elements), remainder, nil
}

// array_pattern = "[" [ pattern { "," pattern } [ "," "..." ] | "..." ] "]" .

func ArrayPattern(tokens []tok.Token) (*ast.ArrayPattern, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = CloseBracket(remainder); found {
		return ast.NewArrayPattern(nil, false), remainder, nil
	}

	hasRest := false
	if remainder, found = expectFunc(tok.TokOpRest)(remainder); found {
		hasRest = true
		if remainder, found = CloseBracket(remainder); !found {
			return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
		}
		return ast.NewArrayPattern(nil, hasRest), remainder, nil
	}

	first, remainder, err := Pattern(remainder)
	if err != nil {
		return nil, remainder, err
	}

	elements := []ast.Pattern{first}
	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		if remainder3, found := expectFunc(tok.TokOpRest)(remainder2); found {
			remainder = remainder3
			hasRest = true
			break
		}

		next, remainder3, err := Pattern(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("pattern", remainder2)
		} else if err != nil {
			return nil, remainder3, err
		}

		elements = append(elements, next)
		remainder = remainder3
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewArrayPattern(elements, hasRest), remainder, nil
}
