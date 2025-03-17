package ast

import "strings"

// UnionMemberDeclaration represents a member of a union declaration
type UnionMemberDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The union member name
	Parameters  []Node          // Optional parameters (for tuple-like variant)
	Annotations []*Annotation   // Member annotations
	Docs        string          // Documentation comments
}

// NewUnionMemberDeclaration creates a new UnionMemberDeclaration node
func NewUnionMemberDeclaration(name *TypeIdentifier, parameters []Node, annotations []*Annotation, docs string) *UnionMemberDeclaration {
	return &UnionMemberDeclaration{
		BaseNode:    BaseNode{NodeType: NodeUnionMemberDeclaration},
		Name:        name,
		Parameters:  parameters,
		Annotations: annotations,
		Docs:        docs,
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

// Children returns the child nodes
func (u *UnionMemberDeclaration) Children() []Node {
	children := make([]Node, 0, len(u.Annotations)+len(u.Parameters)+1)
	for _, annotation := range u.Annotations {
		children = append(children, annotation)
	}

	children = append(children, u.Name)

	children = append(children, u.Parameters...)

	return children
}

// UnionMembers represents a collection of union members
type UnionMembers struct {
	BaseNode
	Members []*UnionMemberDeclaration // The union members
}

// NewUnionMembers creates a new UnionMembers node
func NewUnionMembers(members []*UnionMemberDeclaration) *UnionMembers {
	return &UnionMembers{
		BaseNode: BaseNode{NodeType: NodeUnionMembers},
		Members:  members,
	}
}

// String returns a textual representation of the union members
func (u *UnionMembers) String() string {
	var builder strings.Builder
	for i, member := range u.Members {
		if i > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(member.String())
	}
	return builder.String()
}

// Children returns the child nodes
func (u *UnionMembers) Children() []Node {
	children := make([]Node, len(u.Members))
	for i, member := range u.Members {
		children[i] = member
	}
	return children
}

// UnionDeclaration represents a union declaration
type UnionDeclaration struct {
	BaseNode
	Name        *TypeIdentifier     // The union name
	TypeParams  []*GenericTypeParam // Type parameters if generic
	Members     *UnionMembers       // The union members
	Annotations []*Annotation       // Union annotations
	Docs        string              // Documentation comments
}

// NewUnionDeclaration creates a new UnionDeclaration node
func NewUnionDeclaration(name *TypeIdentifier, typeParams []*GenericTypeParam, members *UnionMembers, annotations []*Annotation, docs string) *UnionDeclaration {
	return &UnionDeclaration{
		BaseNode:    BaseNode{NodeType: NodeUnionDeclaration},
		Name:        name,
		TypeParams:  typeParams,
		Members:     members,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the union declaration
func (u *UnionDeclaration) String() string {
	var builder strings.Builder
	for _, annotation := range u.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString("\n")
	}

	builder.WriteString("union ")
	builder.WriteString(u.Name.String())

	if len(u.TypeParams) > 0 {
		builder.WriteString("<")
		for i, param := range u.TypeParams {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(param.String())
		}
		builder.WriteString(">")
	}

	builder.WriteString(" {\n")
	if u.Members != nil {
		builder.WriteString(u.Members.String())
	}
	builder.WriteString("\n}")

	return builder.String()
}

// Children returns the child nodes
func (u *UnionDeclaration) Children() []Node {
	children := make([]Node, 0, len(u.Annotations)+len(u.TypeParams)+2)
	for _, annotation := range u.Annotations {
		children = append(children, annotation)
	}

	children = append(children, u.Name)

	for _, param := range u.TypeParams {
		children = append(children, param)
	}

	if u.Members != nil {
		children = append(children, u.Members)
	}

	return children
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
