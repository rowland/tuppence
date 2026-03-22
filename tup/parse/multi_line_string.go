package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// multi_line_string_literal = "```" [ function_call_context ] eol { indented_line } indented_closing .

func MultiLineStringLiteral(tokens []tok.Token) (*ast.MultiLineStringLiteral, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	t := peek(remainder)
	if t.Type != tok.TokMultiStrLit {
		return nil, tokens, ErrNoMatch
	}

	value := t.Value()
	header, body, err := splitMultiLineStringToken(value)
	if err != nil {
		return nil, remainder, err
	}

	var processor *ast.FunctionCallContext
	if header != "" {
		processor, err = parseFunctionCallContextText(header)
		if err != nil {
			return nil, remainder, err
		}
	}

	contents, err := parseMultiLineStringContents(body)
	if err != nil {
		return nil, remainder, err
	}

	return &ast.MultiLineStringLiteral{
		BaseNode: ast.BaseNode{
			Type:        ast.NodeMultiLineStringLiteral,
			Source:      t.File,
			StartOffset: t.Offset,
			Length:      t.Length,
		},
		Contents:  contents,
		Processor: processor,
	}, remainder[1:], nil
}

func splitMultiLineStringToken(value string) (header string, body string, err error) {
	if !strings.HasPrefix(value, "```") {
		return "", "", ErrNoMatch
	}

	// Split the token into:
	// - the optional processor header on the opening fence line
	// - the dedentable body contents
	// - the closing fence line, which is structural and not part of the value
	rest := value[3:]
	newlineIdx, newlineLen := firstEOL(rest)
	if newlineIdx < 0 {
		return "", "", ErrNoMatch
	}

	header = rest[:newlineIdx]
	bodyAndClosing := rest[newlineIdx+newlineLen:]

	lastNewlineIdx, lastNewlineLen := lastEOL(bodyAndClosing)
	if lastNewlineIdx < 0 {
		return "", "", ErrNoMatch
	}

	// Preserve the line ending immediately before the closing fence so
	// multiline strings behave like mainstream heredocs/triple-quoted strings.
	body = bodyAndClosing[:lastNewlineIdx+lastNewlineLen]
	closing := bodyAndClosing[lastNewlineIdx+lastNewlineLen:]
	closing = strings.TrimRight(closing, "\r\n")
	// Confirm that our first-EOL / last-EOL split really isolated the closing
	// fence line from the body content within this already-tokenized literal.
	if strings.TrimLeft(closing, " \t") != "```" {
		return "", "", ErrNoMatch
	}

	return header, body, nil
}

func parseFunctionCallContextText(text string) (*ast.FunctionCallContext, error) {
	src := source.NewSource([]byte(text), "function_call_context.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		return nil, err
	}

	context, remainder, err := FunctionCallContext(tokens)
	if err != nil {
		return nil, err
	}

	remainder = skipTrivia(remainder)
	if peek(remainder).Type != tok.TokEOF {
		return nil, errorExpectingTokenType(tok.TokEOF, remainder)
	}

	return context, nil
}

func parseMultiLineStringContents(body string) (*ast.InterpolatedStringLiteral, error) {
	lines := splitLinesPreserveEndings(body)
	indent := firstNonEmptyLineIndent(lines)
	var builder strings.Builder

	// Dedent each physical line while preserving its original line ending.
	// The rebuilt body is then parsed with the same interpolation machinery as
	// ordinary interpolated strings, so newlines remain ordinary string content.
	for _, line := range lines {
		builder.WriteString(dedentSegment(line, indent))
	}

	content := builder.String()
	parts, err := parseInterpolatedStringParts(nil, 0, content)
	if err != nil {
		return nil, err
	}

	return &ast.InterpolatedStringLiteral{
		BaseNode: ast.BaseNode{
			Type: ast.NodeInterpolatedStringLiteral,
		},
		Parts: parts,
	}, nil
}

func splitLinesPreserveEndings(body string) []string {
	if body == "" {
		return []string{}
	}

	// SplitAfter keeps '\n' attached to each segment, which lets us preserve
	// LF and CRLF exactly while still dedenting line-by-line. The empty-body
	// special case avoids treating "" as a single fake line segment.
	return strings.SplitAfter(body, "\n")
}

func firstNonEmptyLineIndent(lines []string) string {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		return leadingWhitespace(line)
	}
	return ""
}

func leadingWhitespace(line string) string {
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return line[:i]
}

func dedentSegment(line, indent string) string {
	if strings.TrimSpace(line) == "" {
		return strings.TrimLeft(line, " \t")
	}
	if indent != "" && strings.HasPrefix(line, indent) {
		return line[len(indent):]
	}
	return line
}

func firstEOL(s string) (index int, length int) {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if i > 0 && s[i-1] == '\r' {
				return i - 1, 2
			}
			return i, 1
		}
	}
	return -1, 0
}

func lastEOL(s string) (index int, length int) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			if i > 0 && s[i-1] == '\r' {
				return i - 1, 2
			}
			return i, 1
		}
	}
	return -1, 0
}
