package ast

// function_declaration_type = ( "fn" "(" [ labeled_parameters | parameters ] ")" ( return_type | "_" ) )
//                           | ( "fx" "(" [ labeled_parameters | parameters ] ")" [ return_type | "_" ] ) .

import "strings"

// FunctionDeclarationType represents the type part of a function declaration
type FunctionDeclarationType struct {
	BaseNode
	HasSideEffects bool // Whether the function has side effects (fx vs fn)
	Parameters     Node // The function parameters
	ReturnType     Node // The return type (may be nil)
}

// NewFunctionDeclarationType creates a new FunctionDeclarationType node
func NewFunctionDeclarationType(hasSideEffects bool, parameters Node, returnType Node) *FunctionDeclarationType {
	return &FunctionDeclarationType{
		BaseNode:       BaseNode{Type: NodeFunctionDeclarationType},
		HasSideEffects: hasSideEffects,
		Parameters:     parameters,
		ReturnType:     returnType,
	}
}

// String returns a textual representation of the function declaration type
func (f *FunctionDeclarationType) String() string {
	var result strings.Builder
	if f.HasSideEffects {
		result.WriteString("fx")
	} else {
		result.WriteString("fn")
	}

	result.WriteString(f.Parameters.String())

	if f.ReturnType != nil {
		result.WriteString(f.ReturnType.String())
	}

	return result.String()
}

// function_parameter_types = "[" local_type_reference { "," local_type_reference } "]" .

// FunctionParameterTypes represents the parameter types in a function declaration
type FunctionParameterTypes struct {
	BaseNode
	Parameters []LocalTypeReference // The function parameter types
}

// NewFunctionParameterTypes creates a new FunctionParameterTypes node
func NewFunctionParameterTypes(parameters []LocalTypeReference) *FunctionParameterTypes {
	return &FunctionParameterTypes{
		BaseNode:   BaseNode{Type: NodeFunctionParameterTypes},
		Parameters: parameters,
	}
}

// String returns a textual representation of the function parameter types
func (f *FunctionParameterTypes) String() string {
	if len(f.Parameters) == 0 {
		return ""
	}
	var result strings.Builder
	result.WriteString("[")
	for i, param := range f.Parameters {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.String())
	}
	result.WriteString("]")
	return result.String()
}

// function_type_declaration = function_type_declaration_lhs "=" function_type .

// FunctionTypeDeclaration represents a function type declaration
type FunctionTypeDeclaration struct {
	BaseNode
	Name       *FunctionIdentifier      // The function name
	TypeParams []*GenericTypeParam      // Type parameters if generic
	Type       *FunctionDeclarationType // The function type
}

// NewFunctionTypeDeclaration creates a new FunctionTypeDeclaration node
func NewFunctionTypeDeclaration(name *FunctionIdentifier, typeParams []*GenericTypeParam, functionType *FunctionDeclarationType) *FunctionTypeDeclaration {
	return &FunctionTypeDeclaration{
		BaseNode:   BaseNode{Type: NodeFunctionTypeDeclaration},
		Name:       name,
		TypeParams: typeParams,
		Type:       functionType,
	}
}

// String returns a textual representation of the function type declaration
func (f *FunctionTypeDeclaration) String() string {
	var result strings.Builder
	result.WriteString(f.Name.String())

	if len(f.TypeParams) > 0 {
		result.WriteString("<")
		for i, param := range f.TypeParams {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(param.String())
		}
		result.WriteString(">")
	}

	result.WriteString(": ")
	result.WriteString(f.Type.String())
	return result.String()
}

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

type FunctionCall struct {
	BaseNode
	Function       Node                    // The function being called (can be Identifier, MemberAccess, etc.)
	ParameterTypes *FunctionParameterTypes // The parameter types of the function (may be nil)
	Arguments      *FunctionArguments      // The arguments passed to the function
	FunctionBlock  *FunctionBlock          // Optional function block (for higher-order functions, may be nil)
}

func NewFunctionCall(function Node, parameterTypes *FunctionParameterTypes, arguments *FunctionArguments, functionBlock *FunctionBlock) *FunctionCall {
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
