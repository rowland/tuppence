package ast

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
	result := ""
	if f.HasSideEffects {
		result += "fx"
	} else {
		result += "fn"
	}

	result += f.Parameters.String()

	if f.ReturnType != nil {
		result += " -> " + f.ReturnType.String()
	}

	return result
}

// function_parameter_types = "[" [ local_type_reference { "," local_type_reference } ] "]" .

// FunctionParameterTypes represents the parameter types in a function declaration
type FunctionParameterTypes struct {
	BaseNode
	Parameters []Node // The function parameter types
}

// NewFunctionParameterTypes creates a new FunctionParameterTypes node
func NewFunctionParameterTypes(parameters []Node) *FunctionParameterTypes {
	return &FunctionParameterTypes{
		BaseNode:   BaseNode{Type: NodeFunctionParameterTypes},
		Parameters: parameters,
	}
}

// String returns a textual representation of the function parameter types
func (f *FunctionParameterTypes) String() string {
	result := "("
	for i, param := range f.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param.String()
	}
	result += ")"
	return result
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
	result := f.Name.String()

	if len(f.TypeParams) > 0 {
		result += "<"
		for i, param := range f.TypeParams {
			if i > 0 {
				result += ", "
			}
			result += param.String()
		}
		result += ">"
	}

	result += ": " + f.Type.String()
	return result
}
