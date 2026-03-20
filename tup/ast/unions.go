package ast

import "strings"

// union_member_declaration = annotations named_tuple
//                          | union_member_no_annotations .

type UnionDeclarationMemberType interface {
	Node
	unionDeclarationMemberTypeNode()
}

func (n *NamedTuple) unionDeclarationMemberTypeNode()      {}
func (n *GenericType) unionDeclarationMemberTypeNode()     {}
func (n *DynamicArrayType) unionDeclarationMemberTypeNode() {}
func (n *FixedSizeArrayType) unionDeclarationMemberTypeNode() {}
func (n *TypeReference) unionDeclarationMemberTypeNode()   {}

// UnionMemberDeclaration represents a member of a union declaration.
// The member itself mirrors the grammar directly: either an introduced
// named_tuple or an existing union_member_no_annotations form.
type UnionMemberDeclaration struct {
	BaseNode
	Annotations []Annotation
	Member      UnionDeclarationMemberType
}

// NewUnionMemberDeclaration creates a new UnionMemberDeclaration node.
func NewUnionMemberDeclaration(annotations []Annotation, member UnionDeclarationMemberType) *UnionMemberDeclaration {
	return &UnionMemberDeclaration{
		BaseNode:    BaseNode{Type: NodeUnionMemberDeclaration},
		Annotations: annotations,
		Member:      member,
	}
}

// String returns a textual representation of the union member declaration.
func (u *UnionMemberDeclaration) String() string {
	var builder strings.Builder
	for _, annotation := range u.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString(" ")
	}

	builder.WriteString(u.Member.String())
	return builder.String()
}

// UnionMembers represents a collection of union members
type UnionMembers []*UnionMemberDeclaration // The union members

// UnionDeclaration represents a union declaration
type UnionDeclaration struct {
	BaseNode
	Members UnionMembers
}

// NewUnionDeclaration creates a new UnionDeclaration node
func NewUnionDeclaration(members UnionMembers) *UnionDeclaration {
	return &UnionDeclaration{
		BaseNode: BaseNode{Type: NodeUnionDeclaration},
		Members:  members,
	}
}

// String returns a textual representation of the union declaration
func (u *UnionDeclaration) String() string {
	var builder strings.Builder

	builder.WriteString("union (\n")

	for _, member := range u.Members {
		builder.WriteString(member.String())
		builder.WriteString("\n")
	}
	builder.WriteString(")\n")

	return builder.String()
}

// union_member = named_tuple
//              | generic_type
//              | dynamic_array
//              | fixed_size_array
//              | local_type_reference
//              | contract_declaration .

type UnionMemberType interface {
	Node
	unionMemberTypeNode()
}

func (n *NamedTuple) unionMemberTypeNode()          {}
func (n *GenericType) unionMemberTypeNode()         {}
func (n *DynamicArrayType) unionMemberTypeNode()    {}
func (n *FixedSizeArrayType) unionMemberTypeNode()  {}
func (n *TypeReference) unionMemberTypeNode()       {}
func (n *Identifier) unionMemberTypeNode()          {}
func (n *ContractDeclaration) unionMemberTypeNode() {}

// UnionType represents a union type
type UnionType struct {
	BaseNode
	Members []UnionMemberType
}

// NewUnionType creates a new UnionType node
func NewUnionType(members []UnionMemberType) *UnionType {
	return &UnionType{
		BaseNode: BaseNode{Type: NodeUnionType},
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

// UnionWithError represents a union type that includes an error
type UnionWithError struct {
	BaseNode
	Members       []UnionMemberType // Union members excluding error
	IsExclamation bool   // True if using the ! prefix syntax
}

func NewUnionWithError(members []UnionMemberType, isExclamation bool) *UnionWithError {
	return &UnionWithError{
		BaseNode:       BaseNode{Type: NodeUnionWithError},
		Members:        members,
		IsExclamation:  isExclamation,
	}
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

// InlineUnion represents an inline union type
type InlineUnion struct {
	BaseNode
	UnionType *UnionType
}

// NewInlineUnion creates a new InlineUnion node
func NewInlineUnion(unionType *UnionType) *InlineUnion {
	return &InlineUnion{
		BaseNode:  BaseNode{Type: NodeInlineUnion},
		UnionType: unionType,
	}
}

// String returns a textual representation of the inline union
func (i *InlineUnion) String() string {
	return "(" + i.UnionType.String() + ")"
}
