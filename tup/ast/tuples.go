package ast

// LabeledTupleMembers represents the members of a labeled tuple
type LabeledTupleMembers struct {
	BaseNode
	Members []Node // The labeled tuple members
}

// NewLabeledTupleMembers creates a new LabeledTupleMembers node
func NewLabeledTupleMembers(members []Node) *LabeledTupleMembers {
	return &LabeledTupleMembers{
		BaseNode: BaseNode{NodeType: NodeLabeledTupleMembers},
		Members:  members,
	}
}

// String returns a textual representation of the labeled tuple members
func (l *LabeledTupleMembers) String() string {
	result := ""
	for i, member := range l.Members {
		if i > 0 {
			result += ", "
		}
		result += member.String()
	}
	return result
}

// Children returns the child nodes
func (l *LabeledTupleMembers) Children() []Node {
	return l.Members
}

// LabeledTuple represents a labeled tuple
type LabeledTuple struct {
	BaseNode
	Members *LabeledTupleMembers // The labeled tuple members
}

// NewLabeledTuple creates a new LabeledTuple node
func NewLabeledTuple(members *LabeledTupleMembers) *LabeledTuple {
	return &LabeledTuple{
		BaseNode: BaseNode{NodeType: NodeLabeledTuple},
		Members:  members,
	}
}

// String returns a textual representation of the labeled tuple
func (l *LabeledTuple) String() string {
	return "(" + l.Members.String() + ")"
}

// Children returns the child nodes
func (l *LabeledTuple) Children() []Node {
	return []Node{l.Members}
}
