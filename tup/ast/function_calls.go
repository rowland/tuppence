package ast

import "strings"

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

type FunctionCall struct {
	BaseNode
	Function       Expression              // The function being called
	ParameterTypes *FunctionParameterTypes // The parameter types of the function (may be nil)
	Arguments      *FunctionArguments      // The arguments passed to the function
	FunctionBlock  *FunctionBlock          // Optional function block (for higher-order functions, may be nil)
}

func NewFunctionCall(
	function Expression,
	parameterTypes *FunctionParameterTypes,
	arguments *FunctionArguments,
	functionBlock *FunctionBlock,
) *FunctionCall {
	return &FunctionCall{
		BaseNode:       BaseNode{Type: NodeFunctionCall},
		Function:       function,
		ParameterTypes: parameterTypes,
		Arguments:      arguments,
		FunctionBlock:  functionBlock,
	}
}

// String returns a textual representation of the function call
func (f *FunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString(f.Function.String())

	if f.ParameterTypes != nil {
		builder.WriteString(f.ParameterTypes.String())
	}

	if f.Arguments != nil {
		builder.WriteString(f.Arguments.String())
	}

	if f.FunctionBlock != nil {
		builder.WriteString(" ")
		builder.WriteString(f.FunctionBlock.String())
	}

	return builder.String()
}

// UFCSFunctionCall represents a Uniform Function Call Syntax function call
type UFCSFunctionCall struct {
	BaseNode
	Receiver  Node   // The receiver object
	Function  Node   // The function being called (typically an Identifier)
	Arguments []Node // The arguments passed to the function (excluding the receiver)
}

// NewUFCSFunctionCall creates a new UFCSFunctionCall node
func NewUFCSFunctionCall(receiver Node, function Node, arguments []Node) *UFCSFunctionCall {
	return &UFCSFunctionCall{
		BaseNode:  BaseNode{Type: NodeUFCSFunctionCall},
		Receiver:  receiver,
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the UFCS function call
func (u *UFCSFunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString(u.Receiver.String())
	builder.WriteString(".")
	builder.WriteString(u.Function.String())
	builder.WriteString("(")

	for i, arg := range u.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}

	builder.WriteString(")")
	return builder.String()
}

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
	var result strings.Builder
	result.WriteString("{ ")

	if f.Parameters != nil {
		result.WriteString(f.Parameters.String())
		result.WriteString(" ")
	}
	result.WriteString(f.Body.String())

	result.WriteString(" }")
	return result.String()
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
	Args               *Arguments
	LabeledArgs        *LabeledArguments
	PartialApplication bool // True if the function arguments are a partial application
}

func NewFunctionArguments(args *Arguments, labeledArgs *LabeledArguments, partialApplication bool) *FunctionArguments {
	return &FunctionArguments{
		BaseNode:           BaseNode{Type: NodeFunctionArguments},
		Args:               args,
		LabeledArgs:        labeledArgs,
		PartialApplication: partialApplication,
	}
}

func (f *FunctionArguments) String() string {
	var result strings.Builder
	result.WriteString("(")
	if f.Args != nil {
		result.WriteString(f.Args.String())
	}
	if f.LabeledArgs != nil {
		if f.Args != nil {
			result.WriteString(", ")
		}
		result.WriteString(f.LabeledArgs.String())
	}
	if f.PartialApplication {
		if f.Args != nil || f.LabeledArgs != nil {
			result.WriteString(", ")
		}
		result.WriteString("*")
	}
	result.WriteString(")")
	return result.String()
}
