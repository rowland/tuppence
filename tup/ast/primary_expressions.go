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
	Expression
	primaryExpressionNode()
}

func (n *Block) primaryExpressionNode()                 {}
func (n *IfExpression) primaryExpressionNode()          {}
func (n *ForExpression) primaryExpressionNode()         {}
func (n *InlineForExpression) primaryExpressionNode()   {}
func (n *ArrayFunctionCall) primaryExpressionNode()     {}
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

// func (n *ImportExpression) primaryExpressionNode()      {}

// function_block = "{" [ block_parameters ] block_body "}" .

type FunctionBlock struct {
	BaseNode
	Parameters *BlockParameters // Optional block parameters
	Body       *BlockBody       // The function body
}

func NewFunctionBlock(parameters *BlockParameters, body *BlockBody) *FunctionBlock {
	return &FunctionBlock{
		BaseNode:   BaseNode{Type: NodeFunctionBlock},
		Parameters: parameters,
		Body:       body,
	}
}

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

type FunctionCallContext struct {
	BaseNode
	Function  Node
	Arguments *FunctionArguments
}

func NewFunctionCallContext(function Node, arguments *FunctionArguments) *FunctionCallContext {
	return &FunctionCallContext{
		BaseNode:  BaseNode{Type: NodeFunctionCallContext},
		Function:  function,
		Arguments: arguments,
	}
}

func (f *FunctionCallContext) String() string {
	return f.Function.String() + f.Arguments.String()
}

// function_arguments = ( labeled_arguments
// 	                    | arguments [ "," labeled_arguments ]
// 	                    ) [ partial_application ] .

type FunctionArguments struct {
	BaseNode
	Arguments []Node // List of arguments
}

func NewFunctionArguments(arguments []Node) *FunctionArguments {
	return &FunctionArguments{
		BaseNode:  BaseNode{Type: NodeFunctionArguments},
		Arguments: arguments,
	}
}

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

type LabeledArgument struct {
	BaseNode
	Label *Identifier // The argument label
	Value Node        // The argument value
}

func NewLabeledArgument(label *Identifier, value Node) *LabeledArgument {
	return &LabeledArgument{
		BaseNode: BaseNode{Type: NodeLabeledArgument},
		Label:    label,
		Value:    value,
	}
}

func (l *LabeledArgument) String() string {
	return l.Label.String() + ": " + l.Value.String()
}

// spread_argument = "..." expression .

type SpreadArgument struct {
	BaseNode
	Expression Node // The expression being spread
}

func NewSpreadArgument(expression Node) *SpreadArgument {
	return &SpreadArgument{
		BaseNode:   BaseNode{Type: NodeSpreadArgument},
		Expression: expression,
	}
}

func (s *SpreadArgument) String() string {
	return "..." + s.Expression.String()
}

// partial_application = [ "," ] "*" .

type PartialApplication struct {
	BaseNode
	Function  Node               // The function being partially applied
	Arguments *FunctionArguments // The arguments with placeholders
}

func NewPartialApplication(function Node, arguments *FunctionArguments) *PartialApplication {
	return &PartialApplication{
		BaseNode:  BaseNode{Type: NodePartialApplication},
		Function:  function,
		Arguments: arguments,
	}
}

func (p *PartialApplication) String() string {
	return p.Function.String() + p.Arguments.String()
}

// initializer = assignment .

type Initializer struct {
	BaseNode
	Assignment Assignment
}

func NewInitializer(assignment Assignment) *Initializer {
	return &Initializer{
		BaseNode:   BaseNode{Type: NodeInitializer},
		Assignment: assignment,
	}
}

func (i *Initializer) String() string {
	return i.Assignment.String()
}

// step_expression = expression .

type StepExpression struct {
	BaseNode
	Expression Node // The step expression
}

func NewStepExpression(expression Node) *StepExpression {
	return &StepExpression{
		BaseNode:   BaseNode{Type: NodeStepExpression},
		Expression: expression,
	}
}

func (s *StepExpression) String() string {
	return s.Expression.String()
}

// statement = ( type_qualified_function_declaration
// 	           | type_qualified_declaration
// 	           | type_declaration
// 	           | function_declaration
// 	           | compound_assignment
// 	           | assignment
// 	           | expression
// 	           ) .

type Statement struct {
	BaseNode
	Expression Node // The statement expression
}

func NewStatement(expression Node) *Statement {
	return &Statement{
		BaseNode:   BaseNode{Type: NodeStatement},
		Expression: expression,
	}
}

func (s *Statement) String() string {
	return s.Expression.String() + ";"
}

// iterable = expression .

// Iterable represents an iterable expression used in for-in loops
type Iterable struct {
	BaseNode
	Expression Node // The expression that produces an iterable value
}

func NewIterable(expression Node) *Iterable {
	return &Iterable{
		BaseNode:   BaseNode{Type: NodeIterable},
		Expression: expression,
	}
}

func (i *Iterable) String() string {
	return i.Expression.String()
}
