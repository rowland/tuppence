package ast

type LabeledPattern struct {
	BaseNode
	Label   *Identifier // The pattern label
	Pattern Node        // The underlying pattern
}

func NewLabeledPattern(label *Identifier, pattern Node) *LabeledPattern {
	return &LabeledPattern{
		BaseNode: BaseNode{Type: NodeLabeledPattern},
		Label:    label,
		Pattern:  pattern,
	}
}

func (l *LabeledPattern) String() string {
	return l.Label.String() + ": " + l.Pattern.String()
}

// ListMatch represents a list pattern match (similar to array pattern but for linked lists)
type ListMatch struct {
	BaseNode
	Elements []Node // Pattern elements
	Rest     Node   // Optional rest pattern
}

func NewListMatch(elements []Node, rest Node) *ListMatch {
	return &ListMatch{
		BaseNode: BaseNode{Type: NodeListMatch},
		Elements: elements,
		Rest:     rest,
	}
}

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

type MatchCondition struct {
	BaseNode
	Expression Node // The condition expression
}

func NewMatchCondition(expression Node) *MatchCondition {
	return &MatchCondition{
		BaseNode:   BaseNode{Type: NodeMatchCondition},
		Expression: expression,
	}
}

func (m *MatchCondition) String() string {
	return "if " + m.Expression.String()
}

type PatternMatch struct {
	BaseNode
	Expression Node // Expression being matched
	Pattern    Node // Pattern to match against
}

func NewPatternMatch(expression Node, pattern Node) *PatternMatch {
	return &PatternMatch{
		BaseNode:   BaseNode{Type: NodePatternMatch},
		Expression: expression,
		Pattern:    pattern,
	}
}

func (p *PatternMatch) String() string {
	return p.Expression.String() + " is " + p.Pattern.String()
}

type StructuredMatch struct {
	BaseNode
	Pattern    Node // Pattern to match against
	Expression Node // Expression being matched
}

func NewStructuredMatch(pattern Node, expression Node) *StructuredMatch {
	return &StructuredMatch{
		BaseNode:   BaseNode{Type: NodeStructuredMatch},
		Pattern:    pattern,
		Expression: expression,
	}
}

func (s *StructuredMatch) String() string {
	return s.Pattern.String() + " = " + s.Expression.String()
}
