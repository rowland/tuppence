package ast

// PatternIdentifier represents an identifier in a pattern
type PatternIdentifier struct {
	BaseNode
	Name string
}

// NewPatternIdentifier creates a new PatternIdentifier node
func NewPatternIdentifier(name string) *PatternIdentifier {
	return &PatternIdentifier{
		BaseNode: BaseNode{NodeType: NodePatternIdentifier},
		Name:     name,
	}
}

// String returns a textual representation of the pattern identifier
func (p *PatternIdentifier) String() string {
	return p.Name
}

// TuplePattern represents a tuple pattern in pattern matching
type TuplePattern struct {
	BaseNode
	Elements []Node // Pattern elements
}

// NewTuplePattern creates a new TuplePattern node
func NewTuplePattern(elements []Node) *TuplePattern {
	return &TuplePattern{
		BaseNode: BaseNode{NodeType: NodeTuplePattern},
		Elements: elements,
	}
}

// String returns a textual representation of the tuple pattern
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

// Children returns the child nodes
func (p *TuplePattern) Children() []Node {
	return p.Elements
}

// ArrayPattern represents an array pattern in pattern matching
type ArrayPattern struct {
	BaseNode
	Elements []Node // Pattern elements
	Rest     Node   // Optional rest pattern (e.g., [a, b, ...rest])
}

// NewArrayPattern creates a new ArrayPattern node
func NewArrayPattern(elements []Node, rest Node) *ArrayPattern {
	return &ArrayPattern{
		BaseNode: BaseNode{NodeType: NodeArrayPattern},
		Elements: elements,
		Rest:     rest,
	}
}

// String returns a textual representation of the array pattern
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

// Children returns the child nodes
func (p *ArrayPattern) Children() []Node {
	children := make([]Node, len(p.Elements))
	copy(children, p.Elements)
	if p.Rest != nil {
		children = append(children, p.Rest)
	}
	return children
}

// TypePattern represents a type pattern in pattern matching (e.g., x is Type)
type TypePattern struct {
	BaseNode
	Identifier Node // The identifier to check
	TypeRef    Node // The type to check against
}

// NewTypePattern creates a new TypePattern node
func NewTypePattern(identifier, typeRef Node) *TypePattern {
	return &TypePattern{
		BaseNode:   BaseNode{NodeType: NodeTypePattern},
		Identifier: identifier,
		TypeRef:    typeRef,
	}
}

// String returns a textual representation of the type pattern
func (p *TypePattern) String() string {
	return p.Identifier.String() + " is " + p.TypeRef.String()
}

// Children returns the child nodes
func (p *TypePattern) Children() []Node {
	return []Node{p.Identifier, p.TypeRef}
}

// LiteralPattern represents a literal value in pattern matching
type LiteralPattern struct {
	BaseNode
	Value Node // The literal value to match against
}

// NewLiteralPattern creates a new LiteralPattern node
func NewLiteralPattern(value Node) *LiteralPattern {
	return &LiteralPattern{
		BaseNode: BaseNode{NodeType: NodeLiteralPattern},
		Value:    value,
	}
}

// String returns a textual representation of the literal pattern
func (p *LiteralPattern) String() string {
	return p.Value.String()
}

// Children returns the child nodes
func (p *LiteralPattern) Children() []Node {
	return []Node{p.Value}
}

// WildcardPattern represents a wildcard (_) in pattern matching
type WildcardPattern struct {
	BaseNode
}

// NewWildcardPattern creates a new WildcardPattern node
func NewWildcardPattern() *WildcardPattern {
	return &WildcardPattern{
		BaseNode: BaseNode{NodeType: NodeWildcardPattern},
	}
}

// String returns a textual representation of the wildcard pattern
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

// NewMatchCase creates a new MatchCase node
func NewMatchCase(pattern, body Node, guard Node) *MatchCase {
	return &MatchCase{
		BaseNode: BaseNode{NodeType: NodeMatchCase},
		Pattern:  pattern,
		Body:     body,
		Guard:    guard,
	}
}

// String returns a textual representation of the match case
func (c *MatchCase) String() string {
	result := c.Pattern.String()
	if c.Guard != nil {
		result += " if " + c.Guard.String()
	}
	result += " => " + c.Body.String()
	return result
}

// Children returns the child nodes
func (c *MatchCase) Children() []Node {
	if c.Guard != nil {
		return []Node{c.Pattern, c.Guard, c.Body}
	}
	return []Node{c.Pattern, c.Body}
}

// MatchExpression represents a match expression
type MatchExpression struct {
	BaseNode
	Subject Node         // The expression being matched
	Cases   []*MatchCase // The match cases
}

// NewMatchExpression creates a new MatchExpression node
func NewMatchExpression(subject Node, cases []*MatchCase) *MatchExpression {
	return &MatchExpression{
		BaseNode: BaseNode{NodeType: NodeMatchExpression},
		Subject:  subject,
		Cases:    cases,
	}
}

// String returns a textual representation of the match expression
func (m *MatchExpression) String() string {
	result := "match " + m.Subject.String() + " {\n"
	for _, c := range m.Cases {
		result += "  " + c.String() + "\n"
	}
	result += "}"
	return result
}

// Children returns the child nodes
func (m *MatchExpression) Children() []Node {
	children := []Node{m.Subject}
	for _, c := range m.Cases {
		children = append(children, c)
	}
	return children
}
