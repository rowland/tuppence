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
		BaseNode: BaseNode{Type: NodeLabeledPattern},
		Label:    label,
		Pattern:  pattern,
	}
}

// String returns a textual representation of the labeled pattern
func (l *LabeledPattern) String() string {
	return l.Label.String() + ": " + l.Pattern.String()
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
		BaseNode: BaseNode{Type: NodeListMatch},
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

// MatchCondition represents a condition in a match case
type MatchCondition struct {
	BaseNode
	Expression Node // The condition expression
}

// NewMatchCondition creates a new MatchCondition node
func NewMatchCondition(expression Node) *MatchCondition {
	return &MatchCondition{
		BaseNode:   BaseNode{Type: NodeMatchCondition},
		Expression: expression,
	}
}

// String returns a textual representation of the match condition
func (m *MatchCondition) String() string {
	return "if " + m.Expression.String()
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
		BaseNode:   BaseNode{Type: NodePatternMatch},
		Expression: expression,
		Pattern:    pattern,
	}
}

// String returns a textual representation of the pattern match
func (p *PatternMatch) String() string {
	return p.Expression.String() + " is " + p.Pattern.String()
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
		BaseNode:   BaseNode{Type: NodeStructuredMatch},
		Pattern:    pattern,
		Expression: expression,
	}
}

// String returns a textual representation of the structured match
func (s *StructuredMatch) String() string {
	return s.Pattern.String() + " = " + s.Expression.String()
}
