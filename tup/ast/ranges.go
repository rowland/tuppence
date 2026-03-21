package ast

// RangeBound represents a bound in a range (start or end)
type RangeBound struct {
	BaseNode
	Value       Node // The value of the bound
	IsInclusive bool // Whether the bound is inclusive
}

func NewRangeBound(value Node, isInclusive bool) *RangeBound {
	return &RangeBound{
		BaseNode:    BaseNode{Type: NodeRangeBound},
		Value:       value,
		IsInclusive: isInclusive,
	}
}

func (r *RangeBound) String() string {
	if r.IsInclusive {
		return r.Value.String()
	}
	return r.Value.String() + "!"
}

type Range struct {
	BaseNode
	Start *RangeBound // The start bound
	End   *RangeBound // The end bound
}

func NewRange(start, end *RangeBound) *Range {
	return &Range{
		BaseNode: BaseNode{Type: NodeRange},
		Start:    start,
		End:      end,
	}
}

func (r *Range) String() string {
	return r.Start.String() + ".." + r.End.String()
}

// RestOperator represents the rest/spread operator (...)
type RestOperator struct {
	BaseNode
	Identifier *Identifier // The identifier being spread/rested
}

func NewRestOperator(identifier *Identifier) *RestOperator {
	return &RestOperator{
		BaseNode:   BaseNode{Type: NodeRestOperator},
		Identifier: identifier,
	}
}

func (r *RestOperator) String() string {
	return "..." + r.Identifier.String()
}
