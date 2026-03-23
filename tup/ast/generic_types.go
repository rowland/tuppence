package ast

import "strings"

// type_argument = type .

type TypeNode interface {
	ContractFieldType
	ReturnTypeValue
	typeNode()
}

func (n *TypeReference) typeNode()      {}
func (n *Identifier) typeNode()         {}
func (n *DynamicArrayType) typeNode()   {}
func (n *FixedSizeArrayType) typeNode() {}
func (n *FunctionType) typeNode()       {}
func (n *ErrorTuple) typeNode()         {}
func (n *TupleType) typeNode()          {}
func (n *GenericType) typeNode()        {}
func (n *InlineUnion) typeNode()        {}

type TypeArgument struct {
	BaseNode
	Type TypeNode
}

func NewTypeArgument(item TypeNode) *TypeArgument {
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
