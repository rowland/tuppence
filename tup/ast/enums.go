package ast

import "strings"

// EnumMembers represents a collection of enum members
type EnumMembers struct {
	BaseNode
	Members []*EnumMember // The enum members
}

// NewEnumMembers creates a new EnumMembers node
func NewEnumMembers(members []*EnumMember) *EnumMembers {
	return &EnumMembers{
		BaseNode: BaseNode{NodeType: NodeEnumMembers},
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

// Children returns the child nodes
func (e *EnumMembers) Children() []Node {
	children := make([]Node, len(e.Members))
	for i, member := range e.Members {
		children[i] = member
	}
	return children
}
