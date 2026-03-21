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

// NewTupleLiteral creates a new TupleLiteral node
func NewTupleLiteral(labeled bool, members []*TupleMember) *TupleLiteral {
	return &TupleLiteral{
		BaseNode: BaseNode{Type: NodeTupleLiteral},
		Labeled:  labeled,
		Members:  members,
	}
}

// String returns a textual representation of the tuple literal
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

// LiteralValue returns the Go value of the literal
func (t *TupleLiteral) LiteralValue() any {
	return t.String()
}
