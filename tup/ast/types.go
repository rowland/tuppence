package ast

import (
	"strings"

	"github.com/rowland/tuppence/tup/source"
)

// type_reference = [ identifier { "." identifier } "." ] type_identifier .

type TypeReference struct {
	BaseNode
	Identifiers    []*Identifier   // The identifiers in the type reference
	TypeIdentifier *TypeIdentifier // The type identifier
}

func NewTypeReference(identifiers []*Identifier, typeIdentifier *TypeIdentifier, source *source.Source, startOffset int32, length int32) *TypeReference {
	return &TypeReference{
		BaseNode:       BaseNode{Type: NodeTypeReference, Source: source, StartOffset: startOffset, Length: length},
		TypeIdentifier: typeIdentifier,
		Identifiers:    identifiers,
	}
}

func (t *TypeReference) String() string {
	return t.TypeIdentifier.String()
}

// local_type_reference = type_reference | identifier .

type LocalTypeReference interface {
	Node
	localTypeReferenceNode()
}

func (n *TypeReference) localTypeReferenceNode() {}
func (n *Identifier) localTypeReferenceNode()    {}

// nilable_type = "?" local_type_reference .

type NilableType struct {
	BaseNode
	InnerType Node // The type that is made nilable
}

func NewNilableType(innerType Node) *NilableType {
	return &NilableType{
		BaseNode:  BaseNode{Type: NodeNilableType},
		InnerType: innerType,
	}
}

func (n *NilableType) String() string {
	return "?" + n.InnerType.String()
}

// error_tuple = "error" tuple_type .

type ErrorTuple struct {
	BaseNode
	TupleType *TupleType
}

func NewErrorTuple(tupleType *TupleType) *ErrorTuple {
	return &ErrorTuple{
		BaseNode:  BaseNode{Type: NodeErrorTuple},
		TupleType: tupleType,
	}
}

func (e *ErrorTuple) String() string {
	return "error" + e.TupleType.String()
}

type ArrayElementType interface {
	Node
	arrayElementTypeNode()
}

func (n *TypeReference) arrayElementTypeNode()      {}
func (n *DynamicArrayType) arrayElementTypeNode()   {}
func (n *FixedSizeArrayType) arrayElementTypeNode() {}

// array_type = fixed_size_array | dynamic_array .

type ArrayType struct {
	BaseNode
	ElementType ArrayElementType
}

func (a *ArrayType) String() string {
	return "[" + "]" + a.ElementType.String()
}

// dynamic_array = "[" "]" (type_reference | array_type) .

type DynamicArrayType struct {
	ArrayType
}

func NewDynamicArrayType(elementType ArrayElementType) *DynamicArrayType {
	return &DynamicArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{Type: NodeArrayType},
			ElementType: elementType,
		},
	}
}

func (d *DynamicArrayType) String() string {
	return "[]" + d.ElementType.String()
}

// fixed_size_array = "[" size "]" (type_reference | array_type) .

type FixedSizeArrayType struct {
	ArrayType
	Size Size
}

func NewFixedSizeArrayType(elementType ArrayElementType, size Size) *FixedSizeArrayType {
	return &FixedSizeArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{Type: NodeArrayType},
			ElementType: elementType,
		},
		Size: size,
	}
}

func (f *FixedSizeArrayType) String() string {
	return "[" + f.Size.String() + "]" + f.ElementType.String()
}

// parameter = annotations ( nilable_type
//	                       | type
//	                       | literal
//	                       | union_type
//	                       | union_declaration ) .

type Parameter struct {
	BaseNode
	Annotations []Node // Optional annotations
	Type        Node   // Parameter type
}

func (p *Parameter) String() string {
	var builder strings.Builder

	for _, anno := range p.Annotations {
		builder.WriteString(anno.String())
		builder.WriteString(" ")
	}

	builder.WriteString(p.Type.String())
	return builder.String()
}

// labeled_parameter = annotations identifier ":" ( nilable_type
//	                                              | type
//	                                              | literal
//	                                              | union_type
//	                                              | union_declaration ) .

type LabeledParameter struct {
	BaseNode
	Annotations []Node      // Optional annotations
	Identifier  *Identifier // Parameter name
	Type        Node        // Parameter type
}

func (l *LabeledParameter) String() string {
	var builder strings.Builder

	for _, anno := range l.Annotations {
		builder.WriteString(anno.String())
		builder.WriteString(" ")
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.Type.String())
	return builder.String()
}

// rest_parameter = "..." type .

type RestParameter struct {
	BaseNode
	Type Node // Parameter type
}

func NewRestParameter(paramType Node) *RestParameter {
	return &RestParameter{
		BaseNode: BaseNode{Type: NodeRestParameter},
		Type:     paramType,
	}
}

func (r *RestParameter) String() string {
	return "..." + r.Type.String()
}

// labeled_rest_parameter = annotations identifier ":" rest_parameter .

type LabeledRestParameter struct {
	BaseNode
	Annotations []Node      // Optional annotations
	Identifier  *Identifier // Parameter name
	RestType    Node        // Rest parameter type (changed from *RestParameter to Node)
}

func (l *LabeledRestParameter) String() string {
	var builder strings.Builder

	if len(l.Annotations) > 0 {
		for _, ann := range l.Annotations {
			builder.WriteString(ann.String())
			builder.WriteString(" ")
		}
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.RestType.String())

	return builder.String()
}

// return_type = union_with_error
//             | union_declaration_with_error
//             | nilable_type
//             | type
//             | "error" .

type ReturnType struct {
	BaseNode
	Type Node // Return type
}

func NewReturnType(returnType Node) *ReturnType {
	return &ReturnType{
		BaseNode: BaseNode{Type: NodeReturnType},
		Type:     returnType,
	}
}

func (r *ReturnType) String() string {
	return r.Type.String()
}

// function_type = ( "fn" | "fx" ) "(" [ labeled_parameters | parameters ] ")" return_type .

type FunctionType struct {
	BaseNode
	HasSideEffects bool   // True for 'fx', false for 'fn'
	Parameters     []Node // Can be Parameter, LabeledParameter, RestParameter, or LabeledRestParameter
	ReturnType     Node   // Return type (may be omitted for fx, resulting in nil) (changed from *ReturnType to Node)
}

func NewFunctionType(hasSideEffects bool, parameters []Node, returnType Node) *FunctionType {
	return &FunctionType{
		BaseNode:       BaseNode{Type: NodeFunctionType},
		HasSideEffects: hasSideEffects,
		Parameters:     parameters,
		ReturnType:     returnType,
	}
}

func (f *FunctionType) String() string {
	var builder strings.Builder

	if f.HasSideEffects {
		builder.WriteString("fx")
	} else {
		builder.WriteString("fn")
	}

	builder.WriteString("(")

	for i, param := range f.Parameters {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}

	builder.WriteString(")")

	if !f.HasSideEffects && f.ReturnType != nil {
		builder.WriteString(" -> ")
		builder.WriteString(f.ReturnType.String())
	}

	return builder.String()
}

type TupleTypeMemberNode interface {
	Node
	tupleTypeMemberNode()
}

func (n *TupleTypeMember) tupleTypeMemberNode()        {}
func (n *LabeledTupleTypeMember) tupleTypeMemberNode() {}

// tuple_type_member = annotations ( nilable_type
// 	                               | type
// 	                               | union_type
// 	                               | union_declaration
// 	                               | literal ) .

type TupleTypeMember struct {
	BaseNode
	Annotations *Annotations // Optional annotations
	Type        Node         // Member type
}

func NewTupleTypeMember(annotations *Annotations, memberType Node) *TupleTypeMember {
	return &TupleTypeMember{
		BaseNode:    BaseNode{Type: NodeTupleTypeMember},
		Annotations: annotations,
		Type:        memberType,
	}
}

func (t *TupleTypeMember) String() string {
	var builder strings.Builder

	if t.Annotations != nil {
		builder.WriteString(t.Annotations.String())
	}

	builder.WriteString(t.Type.String())
	return builder.String()
}

// labeled_tuple_type_member = annotations identifier ":" tuple_type_member .

type LabeledTupleTypeMember struct {
	BaseNode
	Annotations *Annotations // Optional annotations
	Identifier  *Identifier  // Field name
	Type        Node         // Field type
}

func NewLabeledTupleTypeMember(annotations *Annotations, identifier *Identifier, memberType Node) *LabeledTupleTypeMember {
	return &LabeledTupleTypeMember{
		BaseNode:    BaseNode{Type: NodeLabeledTupleTypeMember},
		Annotations: annotations,
		Identifier:  identifier,
		Type:        memberType,
	}
}

func (l *LabeledTupleTypeMember) String() string {
	var builder strings.Builder

	if l.Annotations != nil {
		builder.WriteString(l.Annotations.String())
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.Type.String())
	return builder.String()
}

// tuple_type = "(" [ labeled_tuple_type_members | tuple_type_members ] ")" .

type TupleType struct {
	BaseNode
	Members []TupleTypeMemberNode
}

func NewTupleType(members []TupleTypeMemberNode) *TupleType {
	return &TupleType{
		BaseNode: BaseNode{Type: NodeTupleType},
		Members:  members,
	}
}

func (t *TupleType) String() string {
	var builder strings.Builder
	builder.WriteString("(")
	for i, member := range t.Members {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(member.String())
	}
	builder.WriteString(")")
	return builder.String()
}

// type_argument = type | generic_type .

type TypeArgument struct {
	BaseNode
	Type Node // The type
}

func (t *TypeArgument) String() string {
	return t.Type.String()
}

// type_argument_list = "[" type_argument { "," type_argument } "]" .

type TypeArgumentList struct {
	BaseNode
	Arguments []Node // List of TypeArgument nodes
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

// named_tuple = type_identifier tuple_type .

type NamedTuple struct {
	BaseNode
	TypeIdentifier *TypeIdentifier
	TupleType      *TupleType
}

func NewNamedTuple(typeIdentifier *TypeIdentifier, tupleType *TupleType) *NamedTuple {
	return &NamedTuple{
		BaseNode:       BaseNode{Type: NodeNamedTuple},
		TypeIdentifier: typeIdentifier,
		TupleType:      tupleType,
	}
}

func (n *NamedTuple) String() string {
	return n.TypeIdentifier.String() + " " + n.TupleType.String()
}
