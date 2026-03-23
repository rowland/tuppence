package ast

import (
	"strings"

	"github.com/rowland/tuppence/tup/source"
)

// string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - '"' - eol } '"' .

type StringLiteral struct {
	BaseNode
	Value       string
	StringValue string
}

func NewStringLiteral(value string, stringValue string, source *source.Source, startOffset int32, length int32) *StringLiteral {
	return &StringLiteral{
		BaseNode:    BaseNode{Type: NodeStringLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:       value,
		StringValue: stringValue,
	}
}

func (s *StringLiteral) String() string {
	return s.Value
}

// raw_string_literal = "`" { "``" | character - "`" } "`" .

// RawStringLiteral represents a raw string literal enclosed in backticks
type RawStringLiteral struct {
	BaseNode
	Value       string
	StringValue string
}

func NewRawStringLiteral(value string, stringValue string, source *source.Source, startOffset int32, length int32) *RawStringLiteral {
	return &RawStringLiteral{
		BaseNode:    BaseNode{Type: NodeRawStringLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:       value,
		StringValue: stringValue,
	}
}

func (r *RawStringLiteral) String() string {
	return r.Value
}

// interpolation = "\\(" expression ")" .

type InterpolatedStringPart interface {
	Node
	interpolatedStringPartNode()
}

func (s *StringLiteral) interpolatedStringPartNode() {}
func (i *Interpolation) interpolatedStringPartNode() {}

// Interpolation represents an interpolated expression within a string
type Interpolation struct {
	BaseNode
	Expression Expression
}

// NewInterpolation creates a new Interpolation node
func NewInterpolation(expr Expression, source *source.Source) *Interpolation {
	return &Interpolation{
		BaseNode:   BaseNode{Type: NodeInterpolation, Source: source, StartOffset: int32(expr.Pos().Offset), Length: int32(expr.End().Offset - expr.Pos().Offset)},
		Expression: expr,
	}
}

func (i *Interpolation) String() string {
	return `\(` + i.Expression.String() + `)`
}

// interpolated_string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | interpolation | character - '"' - eol } '"' .

// InterpolatedStringLiteral represents a string literal with interpolated expressions
type InterpolatedStringLiteral struct {
	BaseNode
	Parts []InterpolatedStringPart
}

// NewInterpolatedStringLiteral creates a new InterpolatedStringLiteral node
func NewInterpolatedStringLiteral(parts []InterpolatedStringPart, source *source.Source) *InterpolatedStringLiteral {
	return &InterpolatedStringLiteral{
		BaseNode: BaseNode{Type: NodeInterpolatedStringLiteral, Source: source, StartOffset: int32(parts[0].Pos().Offset), Length: int32(parts[len(parts)-1].End().Offset - parts[0].Pos().Offset)},
		Parts:    parts,
	}
}

func (i *InterpolatedStringLiteral) String() string {
	var builder strings.Builder
	builder.WriteString("\"")
	for _, part := range i.Parts {
		builder.WriteString(part.String())
	}
	builder.WriteString("\"")
	return builder.String()
}

// multi_line_string_literal = "```" [ function_call_context ] eol { indented_line } indented_closing .

// MultiLineStringLiteral represents a multi-line string literal with optional processor
type MultiLineStringLiteral struct {
	BaseNode
	Contents  *InterpolatedStringLiteral
	Processor *FunctionCallContext
}

// NewMultiLineStringLiteral creates a new MultiLineStringLiteral node
func NewMultiLineStringLiteral(contents *InterpolatedStringLiteral, processor *FunctionCallContext) *MultiLineStringLiteral {
	return &MultiLineStringLiteral{
		BaseNode:  BaseNode{Type: NodeMultiLineStringLiteral},
		Contents:  contents,
		Processor: processor,
	}
}

// String returns a textual representation of the multi-line string literal
func (m *MultiLineStringLiteral) String() string {
	var builder strings.Builder
	builder.WriteString("```")
	if m.Processor != nil {
		builder.WriteString(m.Processor.String())
	}
	builder.WriteString("\n")
	if m.Contents != nil {
		for _, part := range m.Contents.Parts {
			builder.WriteString(part.String())
		}
	}
	builder.WriteString("```")
	return builder.String()
}
