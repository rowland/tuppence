package ast

// Assignment represents a variable assignment
type Assignment struct {
	BaseNode
	Left  Node // Target of the assignment (typically an identifier)
	Right Node // Value being assigned
}

// NewAssignment creates a new Assignment node
func NewAssignment(left, right Node) *Assignment {
	return &Assignment{
		BaseNode: BaseNode{NodeType: NodeAssignment},
		Left:     left,
		Right:    right,
	}
}

// String returns a textual representation of the assignment
func (a *Assignment) String() string {
	return a.Left.String() + " = " + a.Right.String()
}

// Children returns the child nodes
func (a *Assignment) Children() []Node {
	return []Node{a.Left, a.Right}
}

// DestructuringPattern represents a destructuring pattern for assignments
type DestructuringPattern struct {
	BaseNode
	Pattern Node // The pattern used for destructuring
}

// NewDestructuringPattern creates a new DestructuringPattern node
func NewDestructuringPattern(pattern Node) *DestructuringPattern {
	return &DestructuringPattern{
		BaseNode: BaseNode{NodeType: NodeDestructuringPattern},
		Pattern:  pattern,
	}
}

// String returns a textual representation of the destructuring pattern
func (d *DestructuringPattern) String() string {
	return d.Pattern.String()
}

// Children returns the child nodes
func (d *DestructuringPattern) Children() []Node {
	return []Node{d.Pattern}
}

// DestructuringAssignment represents a destructuring assignment
type DestructuringAssignment struct {
	BaseNode
	Pattern Node // The destructuring pattern
	Value   Node // The value being destructured
}

// NewDestructuringAssignment creates a new DestructuringAssignment node
func NewDestructuringAssignment(pattern, value Node) *DestructuringAssignment {
	return &DestructuringAssignment{
		BaseNode: BaseNode{NodeType: NodeDestructuringAssignment},
		Pattern:  pattern,
		Value:    value,
	}
}

// String returns a textual representation of the destructuring assignment
func (d *DestructuringAssignment) String() string {
	return d.Pattern.String() + " = " + d.Value.String()
}

// Children returns the child nodes
func (d *DestructuringAssignment) Children() []Node {
	return []Node{d.Pattern, d.Value}
}
