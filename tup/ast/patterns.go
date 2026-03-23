package ast

import "strings"

// match_condition = list_match | pattern .

type MatchCondition interface {
	Node
	matchConditionNode()
}

// pattern = wildcard_pattern
//         | typed_pattern
//         | labeled_pattern
//         | tuple_pattern
//         | array_pattern
//         | match_element
//         | type_reference .

type Pattern interface {
	MatchCondition
	patternNode()
}

// match_element = constant | range | type_reference .

type MatchElement interface {
	Node
	matchElementNode()
}

func (n *Constant) matchConditionNode()          {}
func (n *Range) matchConditionNode()             {}
func (n *InferredErrorType) matchConditionNode() {}
func (n *TypeReference) matchConditionNode()     {}
func (n *ListMatch) matchConditionNode()         {}
func (n *WildcardPattern) matchConditionNode()   {}
func (n *TypedPattern) matchConditionNode()      {}
func (n *LabeledPattern) matchConditionNode()    {}
func (n *TuplePattern) matchConditionNode()      {}
func (n *ArrayPattern) matchConditionNode()      {}

func (n *Constant) patternNode()          {}
func (n *Range) patternNode()             {}
func (n *InferredErrorType) patternNode() {}
func (n *TypeReference) patternNode()     {}
func (n *WildcardPattern) patternNode()   {}
func (n *TypedPattern) patternNode()      {}
func (n *LabeledPattern) patternNode()    {}
func (n *TuplePattern) patternNode()      {}
func (n *ArrayPattern) patternNode()      {}

func (n *Constant) matchElementNode()          {}
func (n *Range) matchElementNode()             {}
func (n *InferredErrorType) matchElementNode() {}
func (n *TypeReference) matchElementNode()     {}

// list_match = match_element { "," match_element } .

type ListMatch struct {
	BaseNode
	Elements []MatchElement
}

func NewListMatch(elements []MatchElement) *ListMatch {
	return &ListMatch{
		BaseNode: BaseNode{Type: NodeListMatch},
		Elements: elements,
	}
}

func (l *ListMatch) String() string {
	var builder strings.Builder
	for i, elem := range l.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	return builder.String()
}

// wildcard_pattern = "_" .

type WildcardPattern struct {
	BaseNode
}

func NewWildcardPattern(identifier *Identifier) *WildcardPattern {
	return &WildcardPattern{
		BaseNode: BaseNode{Type: NodeWildcardPattern, Source: identifier.Source, StartOffset: identifier.StartOffset, Length: identifier.Length},
	}
}

func (p *WildcardPattern) String() string {
	return "_"
}

// tuple_pattern = "(" pattern { "," pattern } ")" .

type TuplePattern struct {
	BaseNode
	Elements []Pattern
}

func NewTuplePattern(elements []Pattern) *TuplePattern {
	return &TuplePattern{
		BaseNode: BaseNode{Type: NodeTuplePattern},
		Elements: elements,
	}
}

func (p *TuplePattern) String() string {
	var builder strings.Builder
	builder.WriteString("(")
	for i, elem := range p.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString(")")
	return builder.String()
}

// labeled_pattern = "(" identifier ":" pattern { "," identifier ":" pattern } ")" .

type LabeledPatternMember struct {
	BaseNode
	Label   *Identifier
	Pattern Pattern
}

func NewLabeledPatternMember(label *Identifier, pattern Pattern) *LabeledPatternMember {
	return &LabeledPatternMember{
		BaseNode: BaseNode{Type: NodeLabeledPatternMember},
		Label:    label,
		Pattern:  pattern,
	}
}

func (l *LabeledPatternMember) String() string {
	return l.Label.String() + ": " + l.Pattern.String()
}

type LabeledPattern struct {
	BaseNode
	Members []*LabeledPatternMember
}

func NewLabeledPattern(members []*LabeledPatternMember) *LabeledPattern {
	return &LabeledPattern{
		BaseNode: BaseNode{Type: NodeLabeledPattern},
		Members:  members,
	}
}

func (l *LabeledPattern) String() string {
	var builder strings.Builder
	builder.WriteString("(")
	for i, member := range l.Members {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(member.String())
	}
	builder.WriteString(")")
	return builder.String()
}

// array_pattern = "[" pattern { "," pattern } [ "," rest_operator ] "]" .

type ArrayPattern struct {
	BaseNode
	Elements []Pattern
	HasRest  bool
}

func NewArrayPattern(elements []Pattern, hasRest bool) *ArrayPattern {
	return &ArrayPattern{
		BaseNode: BaseNode{Type: NodeArrayPattern},
		Elements: elements,
		HasRest:  hasRest,
	}
}

func (p *ArrayPattern) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, elem := range p.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	if p.HasRest {
		if len(p.Elements) > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("...")
	}
	builder.WriteString("]")
	return builder.String()
}

// typed_pattern = type_reference pattern .

type TypedPattern struct {
	BaseNode
	Type    *TypeReference
	Pattern Pattern
}

func NewTypedPattern(typeRef *TypeReference, pattern Pattern) *TypedPattern {
	return &TypedPattern{
		BaseNode: BaseNode{Type: NodeTypedPattern},
		Type:     typeRef,
		Pattern:  pattern,
	}
}

func (t *TypedPattern) String() string {
	return t.Type.String() + t.Pattern.String()
}
