package ast

import "strings"

// array_literal = type_identifier "[" [ array_members ] "]"
//               | "[" [ array_members ] "]" .

// ArrayLiteral represents an array literal in the code
type ArrayLiteral struct {
	BaseNode
	Elements      []Expression
	TypeSpecifier *TypeIdentifier
}

// NewArrayLiteral creates a new ArrayLiteral node
func NewArrayLiteral(elements []Expression, typeSpecifier *TypeIdentifier) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode:      BaseNode{Type: NodeArrayLiteral},
		Elements:      elements,
		TypeSpecifier: typeSpecifier,
	}
}

// String returns a textual representation of the array literal
func (a *ArrayLiteral) String() string {
	var builder strings.Builder
	if a.TypeSpecifier != nil {
		builder.WriteString(a.TypeSpecifier.String())
	}
	builder.WriteString("[")
	for i, elem := range a.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString("]")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (a *ArrayLiteral) LiteralValue() any {
	return a.String()
}

// fixed_size_array_literal = fixed_size_array "[" array_members "]" .

// FixedSizeArrayLiteral represents a fixed-size array literal in the code
type FixedSizeArrayLiteral struct {
	BaseNode
	ArrayType *ArrayType
	Elements  []Expression
}

// NewFixedSizeArrayLiteral creates a new FixedSizeArrayLiteral node
func NewFixedSizeArrayLiteral(arrayType *ArrayType, elements []Expression) *FixedSizeArrayLiteral {
	return &FixedSizeArrayLiteral{
		BaseNode:  BaseNode{Type: NodeFixedSizeArrayLiteral},
		ArrayType: arrayType,
		Elements:  elements,
	}
}

// String returns a textual representation of the fixed-size array literal
func (f *FixedSizeArrayLiteral) String() string {
	var builder strings.Builder
	builder.WriteString(f.ArrayType.String())
	builder.WriteString("[")
	for i, elem := range f.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString("]")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (f *FixedSizeArrayLiteral) LiteralValue() any {
	return f.String()
}
