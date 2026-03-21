package ast

import "github.com/rowland/tuppence/tup/source"

// float_literal = decimal_digit { decimal_digit | "_" } "." decimal_digit { decimal_digit | "_" } [ exponent ]
//               | decimal_digit { decimal_digit | "_" } exponent .

// FloatLiteral represents a floating point literal in the code
type FloatLiteral struct {
	BaseNode
	Value      string
	FloatValue float64
}

// NewFloatLiteral creates a new FloatLiteral node
func NewFloatLiteral(value string, floatValue float64, source *source.Source, startOffset int32, length int32) *FloatLiteral {
	return &FloatLiteral{
		BaseNode:   BaseNode{Type: NodeFloatLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:      value,
		FloatValue: floatValue,
	}
}

func (f *FloatLiteral) String() string {
	return f.Value
}

// integer_literal = binary_literal
//                 | hexadecimal_literal
//                 | octal_literal
//                 | decimal_literal .

type IntegerLiteral struct {
	BaseNode
	Value        string
	IntegerValue int64
	Base         int
}

func (i *IntegerLiteral) String() string {
	return i.Value
}

// binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

func NewBinaryLiteral(value string, integerValue int64, source *source.Source, startOffset int32, length int32) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
		Base:         2,
	}
}

// hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
// hex_digit = decimal_digit | "a"-"f" | "A"-"F" .

func NewHexadecimalLiteral(value string, integerValue int64, source *source.Source, startOffset int32, length int32) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
		Base:         16,
	}
}

// octal_literal = "0o" octal_digit { octal_digit } .
// octal_digit = "0"-"7" .

func NewOctalLiteral(value string, integerValue int64, source *source.Source, startOffset int32, length int32) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
		Base:         8,
	}
}

// decimal_literal = decimal_digit { decimal_digit | "_" } .
// decimal_digit = "0"-"9" .

func NewDecimalLiteral(value string, integerValue int64, source *source.Source, startOffset int32, length int32) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
		Base:         10,
	}
}
