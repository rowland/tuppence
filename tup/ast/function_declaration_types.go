package ast

import "strings"

// function_declaration_type = ( "fn" "(" [ labeled_parameters | parameters ] ")" ( return_type | "_" ) )
//                           | ( "fx" "(" [ labeled_parameters | parameters ] ")" [ return_type | "_" ] ) .

type FunctionDeclarationType struct {
	BaseNode
	HasSideEffects bool // Whether the function has side effects (fx vs fn)
	Parameters     []FunctionTypeParameter
	ReturnType     Node // The return type (may be nil)
	InferredReturn bool
}

func NewFunctionDeclarationType(hasSideEffects bool, parameters []FunctionTypeParameter, returnType Node, inferredReturn bool) *FunctionDeclarationType {
	return &FunctionDeclarationType{
		BaseNode:       BaseNode{Type: NodeFunctionDeclarationType},
		HasSideEffects: hasSideEffects,
		Parameters:     parameters,
		ReturnType:     returnType,
		InferredReturn: inferredReturn,
	}
}

func (f *FunctionDeclarationType) String() string {
	var result strings.Builder
	if f.HasSideEffects {
		result.WriteString("fx")
	} else {
		result.WriteString("fn")
	}

	result.WriteString("(")
	for i, param := range f.Parameters {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.String())
	}
	result.WriteString(")")

	if f.InferredReturn {
		result.WriteString(" ")
		result.WriteString("_")
	} else if f.ReturnType != nil {
		result.WriteString(" ")
		result.WriteString(f.ReturnType.String())
	}

	return result.String()
}

// function_parameter_types = "[" local_type_reference { "," local_type_reference } "]" .

type FunctionParameterTypes struct {
	BaseNode
	Parameters []LocalTypeReference // The function parameter types
}

func NewFunctionParameterTypes(parameters []LocalTypeReference) *FunctionParameterTypes {
	return &FunctionParameterTypes{
		BaseNode:   BaseNode{Type: NodeFunctionParameterTypes},
		Parameters: parameters,
	}
}

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
