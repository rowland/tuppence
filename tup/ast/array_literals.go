package ast

import "strings"

// array_initializer = "[" [ array_members ] "]"
//                   | function_block .

type ArrayLiteralType interface {
	Node
	arrayLiteralTypeNode()
}

func (n *TypeReference) arrayLiteralTypeNode()      {}
func (n *FixedSizeArrayType) arrayLiteralTypeNode() {}

// array_literal = fixed_size_array array_initializer
//               | type_reference array_initializer
//               | "[" [ array_members ] "]" .

type ArrayLiteral struct {
	BaseNode
	ArrayType   ArrayLiteralType
	Elements    []Expression
	Initializer *FunctionBlock
}

func NewArrayLiteral(arrayType ArrayLiteralType, elements []Expression, initializer *FunctionBlock) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode:    BaseNode{Type: NodeArrayLiteral},
		ArrayType:   arrayType,
		Elements:    elements,
		Initializer: initializer,
	}
}

func (a *ArrayLiteral) String() string {
	var builder strings.Builder
	if a.ArrayType != nil {
		builder.WriteString(a.ArrayType.String())
	}
	if a.Initializer != nil {
		builder.WriteString(" ")
		builder.WriteString(a.Initializer.String())
		return builder.String()
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
