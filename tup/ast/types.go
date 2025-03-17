package ast

import (
	"strings"
)

// TypeReference represents a reference to a defined type
type TypeReference struct {
	BaseNode
	TypeIdentifier *TypeIdentifier // The type identifier
	IsLocal        bool            // Whether it's a local type reference
}

// NewTypeReference creates a new TypeReference node
func NewTypeReference(typeIdentifier *TypeIdentifier, isLocal bool) *TypeReference {
	return &TypeReference{
		BaseNode:       BaseNode{NodeType: NodeTypeReference},
		TypeIdentifier: typeIdentifier,
		IsLocal:        isLocal,
	}
}

// String returns a textual representation of the type reference
func (t *TypeReference) String() string {
	return t.TypeIdentifier.String()
}

// Children returns the child nodes
func (t *TypeReference) Children() []Node {
	return []Node{t.TypeIdentifier}
}

// NilableType represents a nilable type (prefixed with ?)
type NilableType struct {
	BaseNode
	InnerType Node // The type that is made nilable
}

// NewNilableType creates a new NilableType node
func NewNilableType(innerType Node) *NilableType {
	return &NilableType{
		BaseNode:  BaseNode{NodeType: NodeNilableType},
		InnerType: innerType,
	}
}

// String returns a textual representation of the nilable type
func (n *NilableType) String() string {
	return "?" + n.InnerType.String()
}

// Children returns the child nodes
func (n *NilableType) Children() []Node {
	return []Node{n.InnerType}
}

// ArrayType is the base type for array types
type ArrayType struct {
	BaseNode
	ElementType Node // The type of array elements
}

// String returns a textual representation of the array type
func (a *ArrayType) String() string {
	return "[" + "]" + a.ElementType.String()
}

// Children returns the child nodes
func (a *ArrayType) Children() []Node {
	return []Node{a.ElementType}
}

// DynamicArrayType represents a dynamic-size array type
type DynamicArrayType struct {
	ArrayType
}

// NewDynamicArrayType creates a new DynamicArrayType node
func NewDynamicArrayType(elementType Node) *DynamicArrayType {
	return &DynamicArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{NodeType: NodeArrayType},
			ElementType: elementType,
		},
	}
}

// String returns a textual representation of the dynamic array type
func (d *DynamicArrayType) String() string {
	return "[]" + d.ElementType.String()
}

// FixedSizeArrayType represents a fixed-size array type
type FixedSizeArrayType struct {
	ArrayType
	Size Node // Size expression (can be a literal or identifier)
}

// NewFixedSizeArrayType creates a new FixedSizeArrayType node
func NewFixedSizeArrayType(elementType Node, size Node) *FixedSizeArrayType {
	return &FixedSizeArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{NodeType: NodeArrayType},
			ElementType: elementType,
		},
		Size: size,
	}
}

// String returns a textual representation of the fixed-size array type
func (f *FixedSizeArrayType) String() string {
	return "[" + f.Size.String() + "]" + f.ElementType.String()
}

// Children returns the child nodes
func (f *FixedSizeArrayType) Children() []Node {
	return []Node{f.Size, f.ElementType}
}

// Parameter represents a function parameter
type Parameter struct {
	BaseNode
	Annotations []Node // Optional annotations
	Type        Node   // Parameter type
}

// String returns a textual representation of the parameter
func (p *Parameter) String() string {
	var builder strings.Builder

	for _, anno := range p.Annotations {
		builder.WriteString(anno.String())
		builder.WriteString(" ")
	}

	builder.WriteString(p.Type.String())
	return builder.String()
}

// Children returns the child nodes
func (p *Parameter) Children() []Node {
	children := make([]Node, 0, len(p.Annotations)+1)
	children = append(children, p.Annotations...)
	children = append(children, p.Type)
	return children
}

// LabeledParameter represents a labeled function parameter
type LabeledParameter struct {
	BaseNode
	Annotations []Node      // Optional annotations
	Identifier  *Identifier // Parameter name
	Type        Node        // Parameter type
}

// String returns a textual representation of the labeled parameter
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

// Children returns the child nodes
func (l *LabeledParameter) Children() []Node {
	children := make([]Node, 0, len(l.Annotations)+2)
	children = append(children, l.Annotations...)
	children = append(children, l.Identifier, l.Type)
	return children
}

// RestParameter represents a rest parameter (e.g., ...T)
type RestParameter struct {
	BaseNode
	Type Node // Parameter type
}

// NewRestParameter creates a new RestParameter node
func NewRestParameter(paramType Node) *RestParameter {
	return &RestParameter{
		BaseNode: BaseNode{NodeType: NodeRestParameter},
		Type:     paramType,
	}
}

// String returns a textual representation of the rest parameter
func (r *RestParameter) String() string {
	return "..." + r.Type.String()
}

// Children returns the child nodes
func (r *RestParameter) Children() []Node {
	return []Node{r.Type}
}

// LabeledRestParameter represents a labeled rest parameter (e.g., rest: ...T)
type LabeledRestParameter struct {
	BaseNode
	Annotations []Node      // Optional annotations
	Identifier  *Identifier // Parameter name
	RestType    Node        // Rest parameter type (changed from *RestParameter to Node)
}

// String returns a textual representation of the labeled rest parameter
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

// Children returns the child nodes
func (l *LabeledRestParameter) Children() []Node {
	var children []Node

	children = append(children, l.Annotations...)
	children = append(children, l.Identifier)
	children = append(children, l.RestType)
	return children
}

// ReturnType represents a return type
type ReturnType struct {
	BaseNode
	Type Node // Return type
}

// NewReturnType creates a new ReturnType node
func NewReturnType(returnType Node) *ReturnType {
	return &ReturnType{
		BaseNode: BaseNode{NodeType: NodeReturnType},
		Type:     returnType,
	}
}

// String returns a textual representation of the return type
func (r *ReturnType) String() string {
	return r.Type.String()
}

// Children returns the child nodes
func (r *ReturnType) Children() []Node {
	return []Node{r.Type}
}

// FunctionType represents a function type (fn(params) -> return_type or fx(params))
type FunctionType struct {
	BaseNode
	HasSideEffects bool   // True for 'fx', false for 'fn'
	Parameters     []Node // Can be Parameter, LabeledParameter, RestParameter, or LabeledRestParameter
	ReturnType     Node   // Return type (may be omitted for fx, resulting in nil) (changed from *ReturnType to Node)
}

// NewFunctionType creates a new FunctionType node
func NewFunctionType(hasSideEffects bool, parameters []Node, returnType Node) *FunctionType {
	return &FunctionType{
		BaseNode:       BaseNode{NodeType: NodeFunctionType},
		HasSideEffects: hasSideEffects,
		Parameters:     parameters,
		ReturnType:     returnType,
	}
}

// String returns a textual representation of the function type
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

// Children returns the child nodes
func (f *FunctionType) Children() []Node {
	children := make([]Node, 0, len(f.Parameters)+1)

	children = append(children, f.Parameters...)

	if f.ReturnType != nil {
		children = append(children, f.ReturnType)
	}
	return children
}

// TupleTypeMember represents a member of a tuple type
type TupleTypeMember struct {
	BaseNode
	Annotations []Node // Optional annotations
	Type        Node   // Member type
}

// String returns a textual representation of the tuple type member
func (t *TupleTypeMember) String() string {
	var builder strings.Builder

	for _, anno := range t.Annotations {
		builder.WriteString(anno.String())
		builder.WriteString(" ")
	}

	builder.WriteString(t.Type.String())
	return builder.String()
}

// Children returns the child nodes
func (t *TupleTypeMember) Children() []Node {
	children := make([]Node, 0, len(t.Annotations)+1)
	children = append(children, t.Annotations...)
	children = append(children, t.Type)
	return children
}

// LabeledTupleTypeMember represents a labeled member of a tuple type
type LabeledTupleTypeMember struct {
	BaseNode
	Annotations []Node      // Optional annotations
	Identifier  *Identifier // Field name
	Type        Node        // Field type
}

// String returns a textual representation of the labeled tuple type member
func (l *LabeledTupleTypeMember) String() string {
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

// Children returns the child nodes
func (l *LabeledTupleTypeMember) Children() []Node {
	children := make([]Node, 0, len(l.Annotations)+2)
	children = append(children, l.Annotations...)
	children = append(children, l.Identifier, l.Type)
	return children
}

// TupleType represents a tuple type
type TupleType struct {
	BaseNode
	Members []Node // Mix of TupleTypeMember and LabeledTupleTypeMember nodes
}

// NewTupleType creates a new TupleType node
func NewTupleType(members []Node) *TupleType {
	return &TupleType{
		BaseNode: BaseNode{NodeType: NodeTupleType},
		Members:  members,
	}
}

// String returns a textual representation of the tuple type
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

// Children returns the child nodes
func (t *TupleType) Children() []Node {
	return t.Members
}

// ErrorTuple represents an error tuple type
type ErrorTuple struct {
	BaseNode
	TupleType *TupleType // The tuple type
}

// NewErrorTuple creates a new ErrorTuple node
func NewErrorTuple(tupleType *TupleType) *ErrorTuple {
	return &ErrorTuple{
		BaseNode:  BaseNode{NodeType: NodeErrorTuple},
		TupleType: tupleType,
	}
}

// String returns a textual representation of the error tuple
func (e *ErrorTuple) String() string {
	return "error " + e.TupleType.String()
}

// Children returns the child nodes
func (e *ErrorTuple) Children() []Node {
	return []Node{e.TupleType}
}

// TypeArgument represents a type argument in a generic type
type TypeArgument struct {
	BaseNode
	Type Node // The type
}

// String returns a textual representation of the type argument
func (t *TypeArgument) String() string {
	return t.Type.String()
}

// Children returns the child nodes
func (t *TypeArgument) Children() []Node {
	return []Node{t.Type}
}

// TypeArgumentList represents a list of type arguments for a generic type
type TypeArgumentList struct {
	BaseNode
	Arguments []Node // List of TypeArgument nodes
}

// String returns a textual representation of the type argument list
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

// Children returns the child nodes
func (t *TypeArgumentList) Children() []Node {
	return t.Arguments
}

// GenericType represents a generic type with type arguments
type GenericType struct {
	BaseNode
	BaseType *TypeReference    // The base type
	TypeArgs *TypeArgumentList // Type arguments
}

// NewGenericType creates a new GenericType node
func NewGenericType(baseType *TypeReference, typeArgs *TypeArgumentList) *GenericType {
	return &GenericType{
		BaseNode: BaseNode{NodeType: NodeGenericType},
		BaseType: baseType,
		TypeArgs: typeArgs,
	}
}

// String returns a textual representation of the generic type
func (g *GenericType) String() string {
	return g.BaseType.String() + g.TypeArgs.String()
}

// Children returns the child nodes
func (g *GenericType) Children() []Node {
	return []Node{g.BaseType, g.TypeArgs}
}

// TypeParameter represents a type parameter in a generic type definition
type TypeParameter struct {
	BaseNode
	Identifier *Identifier // The type parameter name
}

// String returns a textual representation of the type parameter
func (t *TypeParameter) String() string {
	return t.Identifier.String()
}

// Children returns the child nodes
func (t *TypeParameter) Children() []Node {
	return []Node{t.Identifier}
}

// TypeParameters represents a list of type parameters for a generic type definition
type TypeParameters struct {
	BaseNode
	Parameters []Node // List of TypeParameter nodes
}

// String returns a textual representation of the type parameters
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

// Children returns the child nodes
func (t *TypeParameters) Children() []Node {
	return t.Parameters
}

// UnionMember represents a member of a union type
type UnionMember struct {
	BaseNode
	Type Node // Member type
}

// String returns a textual representation of the union member
func (u *UnionMember) String() string {
	return u.Type.String()
}

// Children returns the child nodes
func (u *UnionMember) Children() []Node {
	return []Node{u.Type}
}

// UnionType represents a union type
type UnionType struct {
	BaseNode
	Members []Node // List of UnionMember nodes
}

// NewUnionType creates a new UnionType node
func NewUnionType(members []Node) *UnionType {
	return &UnionType{
		BaseNode: BaseNode{NodeType: NodeUnionType},
		Members:  members,
	}
}

// String returns a textual representation of the union type
func (u *UnionType) String() string {
	if len(u.Members) == 0 {
		return "any"
	}

	var builder strings.Builder
	for i, member := range u.Members {
		if i > 0 {
			builder.WriteString(" | ")
		}
		builder.WriteString(member.String())
	}
	return builder.String()
}

// Children returns the child nodes
func (u *UnionType) Children() []Node {
	return u.Members
}

// UnionWithError represents a union type that includes an error
type UnionWithError struct {
	BaseNode
	Members       []Node // List of UnionMember nodes (excluding error)
	IsExclamation bool   // True if using the ! prefix syntax
}

// String returns a textual representation of the union with error
func (u *UnionWithError) String() string {
	if u.IsExclamation {
		if len(u.Members) == 1 {
			return "!" + u.Members[0].String()
		}
	}

	var builder strings.Builder
	for i, member := range u.Members {
		if i > 0 {
			builder.WriteString(" | ")
		}
		builder.WriteString(member.String())
	}
	builder.WriteString(" | error")
	return builder.String()
}

// Children returns the child nodes
func (u *UnionWithError) Children() []Node {
	return u.Members
}

// InlineUnion represents an inline union type
type InlineUnion struct {
	BaseNode
	UnionType *UnionType
}

// NewInlineUnion creates a new InlineUnion node
func NewInlineUnion(unionType *UnionType) *InlineUnion {
	return &InlineUnion{
		BaseNode:  BaseNode{NodeType: NodeInlineUnion},
		UnionType: unionType,
	}
}

// String returns a textual representation of the inline union
func (i *InlineUnion) String() string {
	return "(" + i.UnionType.String() + ")"
}

// Children returns the child nodes
func (i *InlineUnion) Children() []Node {
	return []Node{i.UnionType}
}

// NamedTuple represents a named tuple type
type NamedTuple struct {
	BaseNode
	TypeIdentifier *TypeIdentifier
	TupleType      *TupleType
}

// NewNamedTuple creates a new NamedTuple node
func NewNamedTuple(typeIdentifier *TypeIdentifier, tupleType *TupleType) *NamedTuple {
	return &NamedTuple{
		BaseNode:       BaseNode{NodeType: NodeNamedTuple},
		TypeIdentifier: typeIdentifier,
		TupleType:      tupleType,
	}
}

// String returns a textual representation of the named tuple
func (n *NamedTuple) String() string {
	return n.TypeIdentifier.String() + " " + n.TupleType.String()
}

// Children returns the child nodes
func (n *NamedTuple) Children() []Node {
	return []Node{n.TypeIdentifier, n.TupleType}
}
