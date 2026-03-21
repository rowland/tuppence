package ast

type PatternIdentifier struct {
	BaseNode
	Name string
}

func NewPatternIdentifier(name string) *PatternIdentifier {
	return &PatternIdentifier{
		BaseNode: BaseNode{Type: NodePatternIdentifier},
		Name:     name,
	}
}

func (p *PatternIdentifier) String() string {
	return p.Name
}

type TuplePattern struct {
	BaseNode
	Elements []Node // Pattern elements
}

func NewTuplePattern(elements []Node) *TuplePattern {
	return &TuplePattern{
		BaseNode: BaseNode{Type: NodeTuplePattern},
		Elements: elements,
	}
}

func (p *TuplePattern) String() string {
	result := "("
	for i, elem := range p.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += ")"
	return result
}

type ArrayPattern struct {
	BaseNode
	Elements []Node // Pattern elements
	Rest     Node   // Optional rest pattern (e.g., [a, b, ...rest])
}

func NewArrayPattern(elements []Node, rest Node) *ArrayPattern {
	return &ArrayPattern{
		BaseNode: BaseNode{Type: NodeArrayPattern},
		Elements: elements,
		Rest:     rest,
	}
}

func (p *ArrayPattern) String() string {
	result := "["
	for i, elem := range p.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	if p.Rest != nil {
		if len(p.Elements) > 0 {
			result += ", "
		}
		result += "..." + p.Rest.String()
	}
	result += "]"
	return result
}

// TypePattern represents a type pattern in pattern matching (e.g., x is Type)
type TypePattern struct {
	BaseNode
	Identifier Node // The identifier to check
	TypeRef    Node // The type to check against
}

func NewTypePattern(identifier, typeRef Node) *TypePattern {
	return &TypePattern{
		BaseNode:   BaseNode{Type: NodeTypePattern},
		Identifier: identifier,
		TypeRef:    typeRef,
	}
}

func (p *TypePattern) String() string {
	return p.Identifier.String() + " is " + p.TypeRef.String()
}

type LiteralPattern struct {
	BaseNode
	Value Node // The literal value to match against
}

func NewLiteralPattern(value Node) *LiteralPattern {
	return &LiteralPattern{
		BaseNode: BaseNode{Type: NodeLiteralPattern},
		Value:    value,
	}
}

func (p *LiteralPattern) String() string {
	return p.Value.String()
}

// WildcardPattern represents a wildcard (_) in pattern matching
type WildcardPattern struct {
	BaseNode
}

func NewWildcardPattern() *WildcardPattern {
	return &WildcardPattern{
		BaseNode: BaseNode{Type: NodeWildcardPattern},
	}
}

func (p *WildcardPattern) String() string {
	return "_"
}

// MatchCase represents a case in a match expression
type MatchCase struct {
	BaseNode
	Pattern Node // The pattern to match against
	Body    Node // The body to execute if pattern matches
	Guard   Node // Optional guard condition
}

func NewMatchCase(pattern, body Node, guard Node) *MatchCase {
	return &MatchCase{
		BaseNode: BaseNode{Type: NodeMatchCase},
		Pattern:  pattern,
		Body:     body,
		Guard:    guard,
	}
}

func (c *MatchCase) String() string {
	result := c.Pattern.String()
	if c.Guard != nil {
		result += " if " + c.Guard.String()
	}
	result += " => " + c.Body.String()
	return result
}

type MatchExpression struct {
	BaseNode
	Subject Node         // The expression being matched
	Cases   []*MatchCase // The match cases
}

func NewMatchExpression(subject Node, cases []*MatchCase) *MatchExpression {
	return &MatchExpression{
		BaseNode: BaseNode{Type: NodeMatchExpression},
		Subject:  subject,
		Cases:    cases,
	}
}

func (m *MatchExpression) String() string {
	result := "match " + m.Subject.String() + " {\n"
	for _, c := range m.Cases {
		result += "  " + c.String() + "\n"
	}
	result += "}"
	return result
}
