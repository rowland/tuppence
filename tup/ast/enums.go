package ast

import "strings"

// EnumMember represents a member of an enum
type EnumMember struct {
	BaseNode
	Name  string // The name of the enum member
	Value Node   // Optional explicit value
	Docs  string // Documentation comments
}

// NewEnumMember creates a new EnumMember node
func NewEnumMember(name string, value Node, docs string) *EnumMember {
	return &EnumMember{
		BaseNode: BaseNode{Type: NodeEnumMember},
		Name:     name,
		Value:    value,
		Docs:     docs,
	}
}

// String returns a textual representation of the enum member
func (m *EnumMember) String() string {
	if m.Value != nil {
		return m.Name + " = " + m.Value.String()
	}
	return m.Name
}

// EnumDeclaration represents an enum declaration
type EnumDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The name of the enum
	Members     []*EnumMember   // The enum members
	Annotations []Annotation    // Annotations applied to the enum
}

// NewEnumDeclaration creates a new EnumDeclaration node
func NewEnumDeclaration(name *TypeIdentifier, members []*EnumMember, annotations []Annotation) *EnumDeclaration {
	return &EnumDeclaration{
		BaseNode:    BaseNode{Type: NodeEnumDeclaration},
		Name:        name,
		Members:     members,
		Annotations: annotations,
	}
}

// String returns a textual representation of the enum declaration
func (d *EnumDeclaration) String() string {
	result := ""

	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	result += "enum " + d.Name.String() + " {\n"

	for _, member := range d.Members {
		result += "  " + member.String() + ",\n"
	}

	result += "}"
	return result
}

// EnumMembers represents a collection of enum members
type EnumMembers struct {
	BaseNode
	Members []*EnumMember // The enum members
}

// NewEnumMembers creates a new EnumMembers node
func NewEnumMembers(members []*EnumMember) *EnumMembers {
	return &EnumMembers{
		BaseNode: BaseNode{Type: NodeEnumMembers},
		Members:  members,
	}
}

// String returns a textual representation of the enum members
func (e *EnumMembers) String() string {
	var builder strings.Builder
	for i, member := range e.Members {
		if i > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(member.String())
	}
	return builder.String()
}
