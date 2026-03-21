package ast

import "github.com/rowland/tuppence/tup/source"

// boolean_literal = "true" | "false" .

type BooleanLiteral struct {
	BaseNode
	Value        string
	BooleanValue bool
}

func NewBooleanLiteral(value string, booleanValue bool, source *source.Source, startOffset int32, length int32) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode:     BaseNode{Type: NodeBooleanLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		BooleanValue: booleanValue,
	}
}

func (b *BooleanLiteral) String() string {
	return b.Value
}

// rune_literal = "'" ( byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - eol ) "'" .

type RuneLiteral struct {
	BaseNode
	Value     string
	RuneValue rune
}

func NewRuneLiteral(value string, runeValue rune, source *source.Source, startOffset int32, length int32) *RuneLiteral {
	return &RuneLiteral{
		BaseNode:  BaseNode{Type: NodeRuneLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:     value,
		RuneValue: runeValue,
	}
}

func (r *RuneLiteral) String() string {
	return r.Value
}

// symbol_literal = ":" identifier .

// SymbolLiteral represents a symbol literal in the code (e.g., :name)
type SymbolLiteral struct {
	BaseNode
	Value string
}

func NewSymbolLiteral(value string, source *source.Source, startOffset int32, length int32) *SymbolLiteral {
	return &SymbolLiteral{
		BaseNode: BaseNode{Type: NodeSymbolLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:    value,
	}
}

func (s *SymbolLiteral) String() string {
	return s.Value
}
