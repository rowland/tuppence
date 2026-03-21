package ast

import "strings"

// enum_member_declaration = annotations identifier [ "=" integer_literal ] .

type EnumMember struct {
	BaseNode
	Annotations *Annotations
	Name        *Identifier
	Value       *IntegerLiteral
}

func NewEnumMember(annotations *Annotations, name *Identifier, value *IntegerLiteral) *EnumMember {
	return &EnumMember{
		BaseNode:    BaseNode{Type: NodeEnumMember},
		Annotations: annotations,
		Name:        name,
		Value:       value,
	}
}

func (m *EnumMember) String() string {
	var builder strings.Builder
	if m.Annotations != nil {
		builder.WriteString(m.Annotations.String())
	}
	builder.WriteString(m.Name.String())
	if m.Value != nil {
		builder.WriteString(" = ")
		builder.WriteString(m.Value.String())
	}
	return builder.String()
}

// enum_members = enum_member_declaration { eol enum_member_declaration } eol .

type EnumMembers struct {
	BaseNode
	Members []*EnumMember
}

func NewEnumMembers(members []*EnumMember) *EnumMembers {
	return &EnumMembers{
		BaseNode: BaseNode{Type: NodeEnumMembers},
		Members:  members,
	}
}

func (e *EnumMembers) String() string {
	var builder strings.Builder
	for i, member := range e.Members {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(member.String())
	}
	return builder.String()
}

// enum_declaration = "enum" "(" eol enum_members ")" .

type EnumDeclaration struct {
	BaseNode
	Members *EnumMembers
}

func NewEnumDeclaration(members *EnumMembers) *EnumDeclaration {
	return &EnumDeclaration{
		BaseNode: BaseNode{Type: NodeEnumDeclaration},
		Members:  members,
	}
}

func (d *EnumDeclaration) String() string {
	var builder strings.Builder
	builder.WriteString("enum(\n")
	if d.Members != nil {
		builder.WriteString(indentString(d.Members.String()))
		builder.WriteString("\n")
	}
	builder.WriteString(")")
	return builder.String()
}
