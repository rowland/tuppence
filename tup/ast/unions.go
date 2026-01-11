package ast

import "strings"

// UnionMemberDeclaration represents a member of a union declaration
type UnionMemberDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The union member name
	Parameters  []Node          // Optional parameters (for tuple-like variant)
	Annotations []Annotation    // Member annotations
}

// NewUnionMemberDeclaration creates a new UnionMemberDeclaration node
func NewUnionMemberDeclaration(name *TypeIdentifier, parameters []Node, annotations []Annotation) *UnionMemberDeclaration {
	return &UnionMemberDeclaration{
		BaseNode:    BaseNode{Type: NodeUnionMemberDeclaration},
		Name:        name,
		Parameters:  parameters,
		Annotations: annotations,
	}
}

// String returns a textual representation of the union member declaration
func (u *UnionMemberDeclaration) String() string {
	var builder strings.Builder
	for _, annotation := range u.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString(" ")
	}

	builder.WriteString(u.Name.String())

	if len(u.Parameters) > 0 {
		builder.WriteString("(")
		for i, param := range u.Parameters {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(param.String())
		}
		builder.WriteString(")")
	}

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
func NewUnionDeclaration() *UnionDeclaration {
	return &UnionDeclaration{
		BaseNode: BaseNode{Type: NodeUnionDeclaration},
		Members:  UnionMembers{},
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

// UnionMember represents a member of a union type
type UnionMember struct {
	BaseNode
	Type Node // Member type
}

// String returns a textual representation of the union member
func (u *UnionMember) String() string {
	return u.Type.String()
}

// UnionType represents a union type
type UnionType struct {
	BaseNode
	Members []Node // List of UnionMember nodes
}

// NewUnionType creates a new UnionType node
func NewUnionType(members []Node) *UnionType {
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
