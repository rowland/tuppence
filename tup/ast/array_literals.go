package ast

import "strings"

// array_literal = type_identifier "[" [ array_members ] "]"
//               | "[" [ array_members ] "]" .

type ArrayLiteral struct {
	BaseNode
	Elements      []Expression
	TypeSpecifier *TypeIdentifier
}

func NewArrayLiteral(elements []Expression, typeSpecifier *TypeIdentifier) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode:      BaseNode{Type: NodeArrayLiteral},
		Elements:      elements,
		TypeSpecifier: typeSpecifier,
	}
}

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

// fixed_size_array_literal = fixed_size_array "[" array_members "]" .

type FixedSizeArrayLiteralType interface {
	Node
	fixedSizeArrayLiteralTypeNode()
}

func (n *FixedSizeArrayType) fixedSizeArrayLiteralTypeNode() {}
func (n *TypeReference) fixedSizeArrayLiteralTypeNode()      {}

type FixedSizeArrayLiteral struct {
	BaseNode
	ArrayType   FixedSizeArrayLiteralType
	Elements    []Expression
	Initializer *FunctionBlock
}

func NewFixedSizeArrayLiteral(arrayType FixedSizeArrayLiteralType, elements []Expression, initializer *FunctionBlock) *FixedSizeArrayLiteral {
	return &FixedSizeArrayLiteral{
		BaseNode:    BaseNode{Type: NodeFixedSizeArrayLiteral},
		ArrayType:   arrayType,
		Elements:    elements,
		Initializer: initializer,
	}
}

func (f *FixedSizeArrayLiteral) String() string {
	var builder strings.Builder
	builder.WriteString(f.ArrayType.String())
	if f.Initializer != nil {
		builder.WriteString(" ")
		builder.WriteString(f.Initializer.String())
		return builder.String()
	}
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
