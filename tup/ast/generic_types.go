package ast

import "strings"

// type_argument = type | generic_type .

type TypeArgumentType interface {
	Node
	typeArgumentTypeNode()
}

func (n *TypeReference) typeArgumentTypeNode()      {}
func (n *Identifier) typeArgumentTypeNode()         {}
func (n *DynamicArrayType) typeArgumentTypeNode()   {}
func (n *FixedSizeArrayType) typeArgumentTypeNode() {}
func (n *FunctionType) typeArgumentTypeNode()       {}
func (n *ErrorTuple) typeArgumentTypeNode()         {}
func (n *TupleType) typeArgumentTypeNode()          {}
func (n *GenericType) typeArgumentTypeNode()        {}
func (n *InlineUnion) typeArgumentTypeNode()        {}

type TypeArgument struct {
	BaseNode
	Type TypeArgumentType
}

func NewTypeArgument(item TypeArgumentType) *TypeArgument {
	return &TypeArgument{
		BaseNode: BaseNode{Type: NodeTypeArgument},
		Type:     item,
	}
}

func (t *TypeArgument) String() string {
	return t.Type.String()
}

// type_argument_list = "[" type_argument { "," type_argument } "]" .

type TypeArgumentList struct {
	BaseNode
	Arguments []*TypeArgument
}

func NewTypeArgumentList(arguments []*TypeArgument) *TypeArgumentList {
	return &TypeArgumentList{
		BaseNode:  BaseNode{Type: NodeTypeArgumentList},
		Arguments: arguments,
	}
}

func (t *TypeArgumentList) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, arg := range t.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}
	builder.WriteString("]")
	return builder.String()
}

// generic_type = type_reference type_argument_list .

type GenericType struct {
	BaseNode
	BaseType *TypeReference    // The base type
	TypeArgs *TypeArgumentList // Type arguments
}

func NewGenericType(baseType *TypeReference, typeArgs *TypeArgumentList) *GenericType {
	return &GenericType{
		BaseNode: BaseNode{Type: NodeGenericType},
		BaseType: baseType,
		TypeArgs: typeArgs,
	}
}

func (g *GenericType) String() string {
	return g.BaseType.String() + g.TypeArgs.String()
}

// type_parameter = identifier .

type TypeParameter struct {
	BaseNode
	Identifier *Identifier // The type parameter name
}

func NewTypeParameter(identifier *Identifier) *TypeParameter {
	return &TypeParameter{
		BaseNode:   BaseNode{Type: NodeTypeParameter},
		Identifier: identifier,
	}
}

func (t *TypeParameter) String() string {
	return t.Identifier.String()
}

// type_parameters = "[" type_parameter { "," type_parameter } "]" .

type TypeParameters struct {
	BaseNode
	Parameters []*TypeParameter
}

func NewTypeParameters(parameters []*TypeParameter) *TypeParameters {
	return &TypeParameters{
		BaseNode:   BaseNode{Type: NodeTypeParameters},
		Parameters: parameters,
	}
}

func (t *TypeParameters) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, param := range t.Parameters {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}
	builder.WriteString("]")
	return builder.String()
}
