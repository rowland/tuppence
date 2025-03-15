package ast

// String interpolation related node types
const (
	NodeStringInterpolationEscape NodeType = "StringInterpolationEscape"
)

// StringInterpolationEscape represents an escape sequence in a string interpolation
type StringInterpolationEscape struct {
	BaseNode
	Expression Node // The expression being interpolated
}

// NewStringInterpolationEscape creates a new StringInterpolationEscape node
func NewStringInterpolationEscape(expression Node) *StringInterpolationEscape {
	return &StringInterpolationEscape{
		BaseNode:   BaseNode{NodeType: NodeStringInterpolationEscape},
		Expression: expression,
	}
}

// String returns a textual representation of the string interpolation escape
func (s *StringInterpolationEscape) String() string {
	return "${" + s.Expression.String() + "}"
}

// Children returns the child nodes
func (s *StringInterpolationEscape) Children() []Node {
	return []Node{s.Expression}
}
