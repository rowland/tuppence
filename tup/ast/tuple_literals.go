package ast

import "strings"

// tuple_member = expression .
// labeled_tuple_member = identifier ":" tuple_member .

type TupleMember struct {
	BaseNode
	Label *Identifier
	Value Expression
}

func NewTupleMember(label *Identifier, value Expression) *TupleMember {
	return &TupleMember{
		BaseNode: BaseNode{Type: NodeTupleMember},
		Label:    label,
		Value:    value,
	}
}

func (t *TupleMember) String() string {
	var builder strings.Builder
	if t.Label != nil {
		builder.WriteString(t.Label.String())
		builder.WriteString(": ")
	}
	builder.WriteString(t.Value.String())
	return builder.String()
}

// tuple_literal = empty_tuple | labeled_tuple_members | tuple_members .

type TupleLiteral struct {
	BaseNode
	Labeled bool
	Members []*TupleMember
}

func NewTupleLiteral(labeled bool, members []*TupleMember) *TupleLiteral {
	return &TupleLiteral{
		BaseNode: BaseNode{Type: NodeTupleLiteral},
		Labeled:  labeled,
		Members:  members,
	}
}

func (t *TupleLiteral) String() string {
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
