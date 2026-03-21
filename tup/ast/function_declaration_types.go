package ast

import "strings"

// function_declaration_type = ( "fn" "(" [ labeled_parameters | parameters ] ")" ( return_type | "_" ) )
//                           | ( "fx" "(" [ labeled_parameters | parameters ] ")" [ return_type | "_" ] ) .

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
	Name           *TypeIdentifier
	ParameterTypes *FunctionParameterTypes
	Type           *FunctionType
}

func NewFunctionTypeDeclaration(
	name *TypeIdentifier,
	parameterTypes *FunctionParameterTypes,
	functionType *FunctionType,
) *FunctionTypeDeclaration {
	return &FunctionTypeDeclaration{
		BaseNode:       BaseNode{Type: NodeFunctionTypeDeclaration},
		Name:           name,
		ParameterTypes: parameterTypes,
		Type:           functionType,
	}
}

func (f *FunctionTypeDeclaration) String() string {
	var result strings.Builder
	result.WriteString(f.Name.String())

	if f.ParameterTypes != nil {
		result.WriteString(f.ParameterTypes.String())
	}

	result.WriteString(" = ")
	result.WriteString(f.Type.String())
	return result.String()
}
