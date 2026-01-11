package ast

// primary_expression = "(" expression ")"
//	                  | block
//	                  | if_expression
//	                  | for_expression
//	                  | inline_for_expression
//	                  | array_function_call
//	                  | import_expression
//	                  | typeof_expression
//	                  | function_call
//	                  | type_constructor_call
//	                  | return_expression
//	                  | break_expression
//	                  | continue_expression
//	                  | member_access
//	                  | tuple_update_expression
//	                  | safe_indexed_access
//	                  | indexed_access
//	                  | range
//	                  | identifier
//	                  | literal .

type PrimaryExpression interface {
	primaryExpressionNode()
}

func (n *Block) primaryExpressionNode()               {}
func (n *IfExpression) primaryExpressionNode()        {}
func (n *ForExpression) primaryExpressionNode()       {}
func (n *InlineForExpression) primaryExpressionNode() {}
func (n *ArrayFunctionCall) primaryExpressionNode()   {}

// func (n *ImportExpression) primaryExpressionNode()      {}
func (n *TypeofExpression) primaryExpressionNode()      {}
func (n *FunctionCall) primaryExpressionNode()          {}
func (n *TypeConstructorCall) primaryExpressionNode()   {}
func (n *ReturnExpression) primaryExpressionNode()      {}
func (n *BreakExpression) primaryExpressionNode()       {}
func (n *ContinueExpression) primaryExpressionNode()    {}
func (n *MemberAccess) primaryExpressionNode()          {}
func (n *TupleUpdateExpression) primaryExpressionNode() {}
func (n *SafeIndexedAccess) primaryExpressionNode()     {}
func (n *IndexedAccess) primaryExpressionNode()         {}
func (n *Range) primaryExpressionNode()                 {}
func (n *Identifier) primaryExpressionNode()            {}

// func (n *Literal) primaryExpressionNode() {}

// function_block = "{" [ block_parameters ] block_body "}" .

// FunctionBlock represents a block used as a function body
type FunctionBlock struct {
	BaseNode
	Parameters *BlockParameters // Optional block parameters
	Body       *BlockBody       // The function body
}

// NewFunctionBlock creates a new FunctionBlock node
func NewFunctionBlock(parameters *BlockParameters, body *BlockBody) *FunctionBlock {
	return &FunctionBlock{
		BaseNode:   BaseNode{Type: NodeFunctionBlock},
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

// function_call_context = function_identifier [ "(" ( labeled_arguments | arguments [ "," labeled_arguments ] ) [ partial_application ] ")" ] .

// FunctionCallContext represents a context in which a function is called
type FunctionCallContext struct {
	BaseNode
	Function  Node               // The function being called
	Arguments *FunctionArguments // The arguments passed to the function
}

// NewFunctionCallContext creates a new FunctionCallContext node
func NewFunctionCallContext(function Node, arguments *FunctionArguments) *FunctionCallContext {
	return &FunctionCallContext{
		BaseNode:  BaseNode{Type: NodeFunctionCallContext},
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the function call context
func (f *FunctionCallContext) String() string {
	return f.Function.String() + f.Arguments.String()
}

// function_arguments = ( labeled_arguments
// 	| arguments [ "," labeled_arguments ]
// 	) [ partial_application ] .

// FunctionArguments represents the arguments passed to a function
type FunctionArguments struct {
	BaseNode
	Arguments []Node // List of arguments
}

// NewFunctionArguments creates a new FunctionArguments node
func NewFunctionArguments(arguments []Node) *FunctionArguments {
	return &FunctionArguments{
		BaseNode:  BaseNode{Type: NodeFunctionArguments},
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

// labeled_argument = ( identifier ":" expression | spread_argument ) .

// LabeledArgument represents a labeled argument in a function call
type LabeledArgument struct {
	BaseNode
	Label *Identifier // The argument label
	Value Node        // The argument value
}

// NewLabeledArgument creates a new LabeledArgument node
func NewLabeledArgument(label *Identifier, value Node) *LabeledArgument {
	return &LabeledArgument{
		BaseNode: BaseNode{Type: NodeLabeledArgument},
		Label:    label,
		Value:    value,
	}
}

// String returns a textual representation of the labeled argument
func (l *LabeledArgument) String() string {
	return l.Label.String() + ": " + l.Value.String()
}

// spread_argument = "..." expression .

// SpreadArgument represents a spread argument in a function call
type SpreadArgument struct {
	BaseNode
	Expression Node // The expression being spread
}

// NewSpreadArgument creates a new SpreadArgument node
func NewSpreadArgument(expression Node) *SpreadArgument {
	return &SpreadArgument{
		BaseNode:   BaseNode{Type: NodeSpreadArgument},
		Expression: expression,
	}
}

// String returns a textual representation of the spread argument
func (s *SpreadArgument) String() string {
	return "..." + s.Expression.String()
}

// partial_application = [ "," ] "*" .

// PartialApplication represents a partial function application using the _ placeholder
type PartialApplication struct {
	BaseNode
	Function  Node               // The function being partially applied
	Arguments *FunctionArguments // The arguments with placeholders
}

// NewPartialApplication creates a new PartialApplication node
func NewPartialApplication(function Node, arguments *FunctionArguments) *PartialApplication {
	return &PartialApplication{
		BaseNode:  BaseNode{Type: NodePartialApplication},
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the partial application
func (p *PartialApplication) String() string {
	return p.Function.String() + p.Arguments.String()
}

// initializer = assignment .

// Initializer represents an initializer expression
type Initializer struct {
	BaseNode
	Expression Node // The initializer expression
}

// NewInitializer creates a new Initializer node
func NewInitializer(expression Node) *Initializer {
	return &Initializer{
		BaseNode:   BaseNode{Type: NodeInitializer},
		Expression: expression,
	}
}

// String returns a textual representation of the initializer
func (i *Initializer) String() string {
	return i.Expression.String()
}

// step_expression = expression .

// StepExpression represents a step expression in a for loop
type StepExpression struct {
	BaseNode
	Expression Node // The step expression
}

// NewStepExpression creates a new StepExpression node
func NewStepExpression(expression Node) *StepExpression {
	return &StepExpression{
		BaseNode:   BaseNode{Type: NodeStepExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the step expression
func (s *StepExpression) String() string {
	return s.Expression.String()
}

// statement = ( type_qualified_function_declaration
// 	| type_qualified_declaration
// 	| type_declaration
// 	| function_declaration
// 	| compound_assignment
// 	| assignment
// 	| expression
// 	) .

// Statement represents a statement in a block
type Statement struct {
	BaseNode
	Expression Node // The statement expression
}

// NewStatement creates a new Statement node
func NewStatement(expression Node) *Statement {
	return &Statement{
		BaseNode:   BaseNode{Type: NodeStatement},
		Expression: expression,
	}
}

// String returns a textual representation of the statement
func (s *Statement) String() string {
	return s.Expression.String() + ";"
}

// iterable = expression .

// Iterable represents an iterable expression used in for-in loops
type Iterable struct {
	BaseNode
	Expression Node // The expression that produces an iterable value
}

// NewIterable creates a new Iterable node
func NewIterable(expression Node) *Iterable {
	return &Iterable{
		BaseNode:   BaseNode{Type: NodeIterable},
		Expression: expression,
	}
}

// String returns a textual representation of the iterable
func (i *Iterable) String() string {
	return i.Expression.String()
}
