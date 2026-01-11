package ast

// RangeBound represents a bound in a range (start or end)
type RangeBound struct {
	BaseNode
	Value       Node // The value of the bound
	IsInclusive bool // Whether the bound is inclusive
}

// NewRangeBound creates a new RangeBound node
func NewRangeBound(value Node, isInclusive bool) *RangeBound {
	return &RangeBound{
		BaseNode:    BaseNode{Type: NodeRangeBound},
		Value:       value,
		IsInclusive: isInclusive,
	}
}

// String returns a textual representation of the range bound
func (r *RangeBound) String() string {
	if r.IsInclusive {
		return r.Value.String()
	}
	return r.Value.String() + "!"
}

// Range represents a range expression
type Range struct {
	BaseNode
	Start *RangeBound // The start bound
	End   *RangeBound // The end bound
}

// NewRange creates a new Range node
func NewRange(start, end *RangeBound) *Range {
	return &Range{
		BaseNode: BaseNode{Type: NodeRange},
		Start:    start,
		End:      end,
	}
}

// String returns a textual representation of the range
func (r *Range) String() string {
	return r.Start.String() + ".." + r.End.String()
}

// RestOperator represents the rest/spread operator (...)
type RestOperator struct {
	BaseNode
	Identifier *Identifier // The identifier being spread/rested
}

// NewRestOperator creates a new RestOperator node
func NewRestOperator(identifier *Identifier) *RestOperator {
	return &RestOperator{
		BaseNode:   BaseNode{Type: NodeRestOperator},
		Identifier: identifier,
	}
}

// String returns a textual representation of the rest operator
func (r *RestOperator) String() string {
	return "..." + r.Identifier.String()
}
