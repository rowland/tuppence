package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// interpolation = "\\(" expression ")" .

func Interpolation(tokens []tok.Token) (*ast.Interpolation, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	t := peek(remainder)
	if t.Type != tok.TokInterpStrLit {
		return nil, tokens, ErrNoMatch
	}

	value := t.Value()
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}

	expression, err := parseInterpolationExpression(value[2 : len(value)-1])
	if err != nil {
		return nil, remainder, err
	}

	return &ast.Interpolation{
		BaseNode: ast.BaseNode{
			Type:        ast.NodeInterpolation,
			Source:      t.File,
			StartOffset: t.Offset,
			Length:      t.Length,
		},
		Expression: expression,
	}, remainder[1:], nil
}

// interpolated_string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | interpolation | character - '"' - eol } '"' .

func InterpolatedStringLiteral(tokens []tok.Token) (*ast.InterpolatedStringLiteral, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	t := peek(remainder)
	if t.Type != tok.TokInterpStrLit {
		return nil, tokens, ErrNoMatch
	}

	value := t.Value()
	if len(value) < 2 {
		return nil, remainder, errorExpecting("interpolated string literal", remainder)
	}

	content := value[1 : len(value)-1]
	parts, err := parseInterpolatedStringParts(t.File, t.Offset+1, content)
	if err != nil {
		return nil, remainder, err
	}

	return &ast.InterpolatedStringLiteral{
		BaseNode: ast.BaseNode{
			Type:        ast.NodeInterpolatedStringLiteral,
			Source:      t.File,
			StartOffset: t.Offset,
			Length:      t.Length,
		},
		Parts: parts,
	}, remainder[1:], nil
}

func parseInterpolatedStringParts(src *source.Source, startOffset int32, content string) ([]ast.InterpolatedStringPart, error) {
	parts := []ast.InterpolatedStringPart{}
	segmentStart := 0

	for i := 0; i < len(content); i++ {
		if content[i] != '\\' || i+1 >= len(content) || content[i+1] != '(' {
			continue
		}

		// Flush the preceding literal segment, if any, before parsing the
		// interpolation that starts at this byte offset.
		if i > segmentStart {
			parts = append(parts, ast.NewStringLiteral(
				content[segmentStart:i],
				content[segmentStart:i],
				src,
				startOffset+int32(segmentStart),
				int32(i-segmentStart),
			))
		}

		// Find the matching closing ')' for this interpolation by tokenizing the
		// interpolation body and reusing the tokenizer's parenthesis handling.
		end, expression, err := parseInterpolationAt(content, i+2)
		if err != nil {
			return nil, err
		}

		parts = append(parts, &ast.Interpolation{
			BaseNode: ast.BaseNode{
				Type:        ast.NodeInterpolation,
				Source:      src,
				StartOffset: startOffset + int32(i),
				Length:      int32(end - i + 1),
			},
			Expression: expression,
		})

		i = end
		segmentStart = end + 1
	}

	if segmentStart < len(content) {
		parts = append(parts, ast.NewStringLiteral(
			content[segmentStart:],
			content[segmentStart:],
			src,
			startOffset+int32(segmentStart),
			int32(len(content)-segmentStart),
		))
	}

	return parts, nil
}

func parseInterpolationAt(content string, exprStart int) (end int, expression ast.Expression, err error) {
	// Tokenize from just after the opening "\(" and scan until the tokenizer
	// reports the ')' that closes this interpolation, accounting for nested
	// parentheses in the embedded expression.
	src := source.NewSource([]byte(content[exprStart:]), "interpolation.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		return 0, nil, err
	}

	parens := 0
	closeOffset := -1
	for _, token := range tokens {
		switch token.Type {
		case tok.TokOpenParen:
			parens++
		case tok.TokCloseParen:
			parens--
			if parens < 0 {
				closeOffset = int(token.Offset)
			}
		}
		if closeOffset >= 0 {
			break
		}
		if token.Type == tok.TokEOF {
			break
		}
	}

	if closeOffset < 0 {
		return 0, nil, ErrNoMatch
	}

	// Re-tokenize just the expression body so Expression(...) sees the same
	// input it would see in ordinary source code, without the closing ')'.
	expression, err = parseInterpolationExpression(content[exprStart : exprStart+closeOffset])
	if err != nil {
		return 0, nil, err
	}

	return exprStart + closeOffset, expression, nil
}

func parseInterpolationExpression(exprText string) (ast.Expression, error) {
	src := source.NewSource([]byte(exprText), "interpolation.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		return nil, err
	}

	expression, remainder, err := Expression(tokens)
	if err != nil {
		return nil, err
	}

	remainder = skipTrivia(remainder)
	if peek(remainder).Type != tok.TokEOF {
		return nil, errorExpectingTokenType(tok.TokEOF, remainder)
	}

	return expression, nil
}
