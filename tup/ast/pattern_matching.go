package ast

// LabeledPattern represents a labeled pattern in pattern matching
type LabeledPattern struct {
	BaseNode
	Label   *Identifier // The pattern label
	Pattern Node        // The underlying pattern
}

// NewLabeledPattern creates a new LabeledPattern node
func NewLabeledPattern(label *Identifier, pattern Node) *LabeledPattern {
	return &LabeledPattern{
		BaseNode: BaseNode{NodeType: NodeLabeledPattern},
		Label:    label,
		Pattern:  pattern,
	}
}

// String returns a textual representation of the labeled pattern
func (l *LabeledPattern) String() string {
	return l.Label.String() + ": " + l.Pattern.String()
}

// Children returns the child nodes
func (l *LabeledPattern) Children() []Node {
	return []Node{l.Label, l.Pattern}
}

// ListMatch represents a list pattern match (similar to array pattern but for linked lists)
type ListMatch struct {
	BaseNode
	Elements []Node // Pattern elements
	Rest     Node   // Optional rest pattern
}

// NewListMatch creates a new ListMatch node
func NewListMatch(elements []Node, rest Node) *ListMatch {
	return &ListMatch{
		BaseNode: BaseNode{NodeType: NodeListMatch},
		Elements: elements,
		Rest:     rest,
	}
}

// String returns a textual representation of the list match
func (l *ListMatch) String() string {
	result := "list["
	for i, elem := range l.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	if l.Rest != nil {
		if len(l.Elements) > 0 {
			result += ", "
		}
		result += "..." + l.Rest.String()
	}
	result += "]"
	return result
}

// Children returns the child nodes
func (l *ListMatch) Children() []Node {
	children := make([]Node, len(l.Elements))
	copy(children, l.Elements)
	if l.Rest != nil {
		children = append(children, l.Rest)
	}
	return children
}

// MatchCondition represents a condition in a match case
type MatchCondition struct {
	BaseNode
	Expression Node // The condition expression
}

// NewMatchCondition creates a new MatchCondition node
func NewMatchCondition(expression Node) *MatchCondition {
	return &MatchCondition{
		BaseNode:   BaseNode{NodeType: NodeMatchCondition},
		Expression: expression,
	}
}

// String returns a textual representation of the match condition
func (m *MatchCondition) String() string {
	return "if " + m.Expression.String()
}

// Children returns the child nodes
func (m *MatchCondition) Children() []Node {
	return []Node{m.Expression}
}

// PatternMatch represents a pattern match operation
type PatternMatch struct {
	BaseNode
	Expression Node // Expression being matched
	Pattern    Node // Pattern to match against
}

// NewPatternMatch creates a new PatternMatch node
func NewPatternMatch(expression Node, pattern Node) *PatternMatch {
	return &PatternMatch{
		BaseNode:   BaseNode{NodeType: NodePatternMatch},
		Expression: expression,
		Pattern:    pattern,
	}
}

// String returns a textual representation of the pattern match
func (p *PatternMatch) String() string {
	return p.Expression.String() + " is " + p.Pattern.String()
}

// Children returns the child nodes
func (p *PatternMatch) Children() []Node {
	return []Node{p.Expression, p.Pattern}
}

// StructuredMatch represents a structured match operation
type StructuredMatch struct {
	BaseNode
	Pattern    Node // Pattern to match against
	Expression Node // Expression being matched
}

// NewStructuredMatch creates a new StructuredMatch node
func NewStructuredMatch(pattern Node, expression Node) *StructuredMatch {
	return &StructuredMatch{
		BaseNode:   BaseNode{NodeType: NodeStructuredMatch},
		Pattern:    pattern,
		Expression: expression,
	}
}

// String returns a textual representation of the structured match
func (s *StructuredMatch) String() string {
	return s.Pattern.String() + " = " + s.Expression.String()
}

// Children returns the child nodes
func (s *StructuredMatch) Children() []Node {
	return []Node{s.Pattern, s.Expression}
}

// AssignmentLhs represents the left-hand side of an assignment
type AssignmentLhs struct {
	BaseNode
	Target Node // The assignment target (identifier, member access, etc.)
}

// NewAssignmentLhs creates a new AssignmentLhs node
func NewAssignmentLhs(target Node) *AssignmentLhs {
	return &AssignmentLhs{
		BaseNode: BaseNode{NodeType: NodeAssignmentLhs},
		Target:   target,
	}
}

// String returns a textual representation of the assignment left-hand side
func (a *AssignmentLhs) String() string {
	return a.Target.String()
}

// Children returns the child nodes
func (a *AssignmentLhs) Children() []Node {
	return []Node{a.Target}
}

// LabeledAssignmentLhs represents a labeled left-hand side of an assignment
type LabeledAssignmentLhs struct {
	BaseNode
	Label  *Identifier    // The label
	Target *AssignmentLhs // The assignment target
}

// NewLabeledAssignmentLhs creates a new LabeledAssignmentLhs node
func NewLabeledAssignmentLhs(label *Identifier, target *AssignmentLhs) *LabeledAssignmentLhs {
	return &LabeledAssignmentLhs{
		BaseNode: BaseNode{NodeType: NodeLabeledAssignmentLhs},
		Label:    label,
		Target:   target,
	}
}

// String returns a textual representation of the labeled assignment left-hand side
func (l *LabeledAssignmentLhs) String() string {
	return l.Label.String() + ": " + l.Target.String()
}

// Children returns the child nodes
func (l *LabeledAssignmentLhs) Children() []Node {
	return []Node{l.Label, l.Target}
}
