package ast

// range_bound = postfix_expression .

type RangeBound struct {
	BaseNode
	Value Expression
}

func NewRangeBound(value Expression) *RangeBound {
	return &RangeBound{
		BaseNode: BaseNode{Type: NodeRangeBound},
		Value:    value,
	}
}

func (r *RangeBound) String() string {
	return r.Value.String()
}

// range = range_bound ".." range_bound .

type Range struct {
	BaseNode
	StartBound *RangeBound
	EndBound   *RangeBound
}

func NewRange(start, end *RangeBound) *Range {
	return &Range{
		BaseNode:   BaseNode{Type: NodeRange},
		StartBound: start,
		EndBound:   end,
	}
}

func (r *Range) String() string {
	return r.StartBound.String() + ".." + r.EndBound.String()
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
	if r.Identifier == nil {
		return "..."
	}
	return "..." + r.Identifier.String()
}
