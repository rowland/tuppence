package ast

// PrimaryExpression represents the most basic form of expression
// This includes literals, identifiers, and parenthesized expressions
type PrimaryExpression struct {
	BaseNode
	Expression Node // The actual expression (literal, identifier, etc.)
}

// NewPrimaryExpression creates a new PrimaryExpression node
func NewPrimaryExpression(expression Node) *PrimaryExpression {
	return &PrimaryExpression{
		BaseNode:   BaseNode{NodeType: NodePrimaryExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the primary expression
func (p *PrimaryExpression) String() string {
	return p.Expression.String()
}

// Children returns the child nodes
func (p *PrimaryExpression) Children() []Node {
	return []Node{p.Expression}
}

// FunctionBlock represents a block used as a function body
type FunctionBlock struct {
	BaseNode
	Parameters *BlockParameters // Optional block parameters
	Body       *BlockBody       // The function body
}

// NewFunctionBlock creates a new FunctionBlock node
func NewFunctionBlock(parameters *BlockParameters, body *BlockBody) *FunctionBlock {
	return &FunctionBlock{
		BaseNode:   BaseNode{NodeType: NodeFunctionBlock},
		Parameters: parameters,
		Body:       body,
	}
}

// String returns a textual representation of the function block
func (f *FunctionBlock) String() string {
	result := "{"
	if f.Parameters != nil {
		result += " " + f.Parameters.String() + " "
	}

	if f.Body != nil {
		result += f.Body.String()
	}

	result += "}"
	return result
}

// Children returns the child nodes
func (f *FunctionBlock) Children() []Node {
	var children []Node
	if f.Parameters != nil {
		children = append(children, f.Parameters)
	}
	if f.Body != nil {
		children = append(children, f.Body)
	}
	return children
}

// FunctionCallContext represents a context in which a function is called
type FunctionCallContext struct {
	BaseNode
	Function  Node               // The function being called
	Arguments *FunctionArguments // The arguments passed to the function
}

// NewFunctionCallContext creates a new FunctionCallContext node
func NewFunctionCallContext(function Node, arguments *FunctionArguments) *FunctionCallContext {
	return &FunctionCallContext{
		BaseNode:  BaseNode{NodeType: NodeFunctionCallContext},
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the function call context
func (f *FunctionCallContext) String() string {
	return f.Function.String() + f.Arguments.String()
}

// Children returns the child nodes
func (f *FunctionCallContext) Children() []Node {
	return []Node{f.Function, f.Arguments}
}

// FunctionArguments represents the arguments passed to a function
type FunctionArguments struct {
	BaseNode
	Arguments []Node // List of arguments
}

// NewFunctionArguments creates a new FunctionArguments node
func NewFunctionArguments(arguments []Node) *FunctionArguments {
	return &FunctionArguments{
		BaseNode:  BaseNode{NodeType: NodeFunctionArguments},
		Arguments: arguments,
	}
}

// String returns a textual representation of the function arguments
func (f *FunctionArguments) String() string {
	result := "("
	for i, arg := range f.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

// Children returns the child nodes
func (f *FunctionArguments) Children() []Node {
	return f.Arguments
}

// LabeledArgument represents a labeled argument in a function call
type LabeledArgument struct {
	BaseNode
	Label *Identifier // The argument label
	Value Node        // The argument value
}

// NewLabeledArgument creates a new LabeledArgument node
func NewLabeledArgument(label *Identifier, value Node) *LabeledArgument {
	return &LabeledArgument{
		BaseNode: BaseNode{NodeType: NodeLabeledArgument},
		Label:    label,
		Value:    value,
	}
}

// String returns a textual representation of the labeled argument
func (l *LabeledArgument) String() string {
	return l.Label.String() + ": " + l.Value.String()
}

// Children returns the child nodes
func (l *LabeledArgument) Children() []Node {
	return []Node{l.Label, l.Value}
}

// SpreadArgument represents a spread argument in a function call
type SpreadArgument struct {
	BaseNode
	Expression Node // The expression being spread
}

// NewSpreadArgument creates a new SpreadArgument node
func NewSpreadArgument(expression Node) *SpreadArgument {
	return &SpreadArgument{
		BaseNode:   BaseNode{NodeType: NodeSpreadArgument},
		Expression: expression,
	}
}

// String returns a textual representation of the spread argument
func (s *SpreadArgument) String() string {
	return "..." + s.Expression.String()
}

// Children returns the child nodes
func (s *SpreadArgument) Children() []Node {
	return []Node{s.Expression}
}

// PartialApplication represents a partial function application using the _ placeholder
type PartialApplication struct {
	BaseNode
	Function  Node               // The function being partially applied
	Arguments *FunctionArguments // The arguments with placeholders
}

// NewPartialApplication creates a new PartialApplication node
func NewPartialApplication(function Node, arguments *FunctionArguments) *PartialApplication {
	return &PartialApplication{
		BaseNode:  BaseNode{NodeType: NodePartialApplication},
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the partial application
func (p *PartialApplication) String() string {
	return p.Function.String() + p.Arguments.String()
}

// Children returns the child nodes
func (p *PartialApplication) Children() []Node {
	return []Node{p.Function, p.Arguments}
}

// Initializer represents an initializer expression
type Initializer struct {
	BaseNode
	Expression Node // The initializer expression
}

// NewInitializer creates a new Initializer node
func NewInitializer(expression Node) *Initializer {
	return &Initializer{
		BaseNode:   BaseNode{NodeType: NodeInitializer},
		Expression: expression,
	}
}

// String returns a textual representation of the initializer
func (i *Initializer) String() string {
	return i.Expression.String()
}

// Children returns the child nodes
func (i *Initializer) Children() []Node {
	return []Node{i.Expression}
}

// StepExpression represents a step expression in a for loop
type StepExpression struct {
	BaseNode
	Expression Node // The step expression
}

// NewStepExpression creates a new StepExpression node
func NewStepExpression(expression Node) *StepExpression {
	return &StepExpression{
		BaseNode:   BaseNode{NodeType: NodeStepExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the step expression
func (s *StepExpression) String() string {
	return s.Expression.String()
}

// Children returns the child nodes
func (s *StepExpression) Children() []Node {
	return []Node{s.Expression}
}

// Statement represents a statement in a block
type Statement struct {
	BaseNode
	Expression Node // The statement expression
}

// NewStatement creates a new Statement node
func NewStatement(expression Node) *Statement {
	return &Statement{
		BaseNode:   BaseNode{NodeType: NodeStatement},
		Expression: expression,
	}
}

// String returns a textual representation of the statement
func (s *Statement) String() string {
	return s.Expression.String() + ";"
}

// Children returns the child nodes
func (s *Statement) Children() []Node {
	return []Node{s.Expression}
}

// TopLevelItem represents a top-level item in a module
type TopLevelItem struct {
	BaseNode
	Item Node // The top-level item
}

// NewTopLevelItem creates a new TopLevelItem node
func NewTopLevelItem(item Node) *TopLevelItem {
	return &TopLevelItem{
		BaseNode: BaseNode{NodeType: NodeTopLevelItem},
		Item:     item,
	}
}

// String returns a textual representation of the top-level item
func (t *TopLevelItem) String() string {
	return t.Item.String()
}

// Children returns the child nodes
func (t *TopLevelItem) Children() []Node {
	return []Node{t.Item}
}

// Iterable represents an iterable expression used in for-in loops
type Iterable struct {
	BaseNode
	Expression Node // The expression that produces an iterable value
}

// NewIterable creates a new Iterable node
func NewIterable(expression Node) *Iterable {
	return &Iterable{
		BaseNode:   BaseNode{NodeType: NodeIterable},
		Expression: expression,
	}
}

// String returns a textual representation of the iterable
func (i *Iterable) String() string {
	return i.Expression.String()
}

// Children returns the child nodes
func (i *Iterable) Children() []Node {
	return []Node{i.Expression}
}
