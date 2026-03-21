package ast

import "strings"

type GenericTypeParam struct {
	BaseNode
	Name       string // The name of the type parameter
	Constraint Node   // Optional type constraint
}

func NewGenericTypeParam(name string, constraint Node) *GenericTypeParam {
	return &GenericTypeParam{
		BaseNode:   BaseNode{Type: NodeGenericTypeParam},
		Name:       name,
		Constraint: constraint,
	}
}

func (t *GenericTypeParam) String() string {
	if t.Constraint != nil {
		return t.Name + ": " + t.Constraint.String()
	}
	return t.Name
}

// function_declaration_lhs = function_identifier [ function_parameter_types ] .

type FunctionDeclarationLHS struct {
	BaseNode
	Name           *FunctionIdentifier
	ParameterTypes *FunctionParameterTypes
}

func NewFunctionDeclarationLHS(name *FunctionIdentifier, parameterTypes *FunctionParameterTypes) *FunctionDeclarationLHS {
	return &FunctionDeclarationLHS{
		BaseNode:       BaseNode{Type: NodeFunctionDeclarationLHS},
		Name:           name,
		ParameterTypes: parameterTypes,
	}
}

func (f *FunctionDeclarationLHS) String() string {
	var result strings.Builder
	result.WriteString(f.Name.String())
	if f.ParameterTypes != nil {
		result.WriteString(f.ParameterTypes.String())
	}
	return result.String()
}

// function_declaration = annotations function_declaration_lhs "=" function_declaration_type block .

type FunctionDeclaration struct {
	BaseNode
	Annotations []Annotation
	LHS         *FunctionDeclarationLHS
	Type        *FunctionDeclarationType
	Body        *Block
}

func NewFunctionDeclaration(annotations []Annotation, lhs *FunctionDeclarationLHS, functionType *FunctionDeclarationType, body *Block) *FunctionDeclaration {
	return &FunctionDeclaration{
		BaseNode:    BaseNode{Type: NodeFunctionDeclaration},
		Annotations: annotations,
		LHS:         lhs,
		Type:        functionType,
		Body:        body,
	}
}

func (d *FunctionDeclaration) String() string {
	var result strings.Builder
	for _, a := range d.Annotations {
		result.WriteString(a.String())
		result.WriteString("\n")
	}
	result.WriteString(d.LHS.String())
	result.WriteString(" = ")
	result.WriteString(d.Type.String())
	result.WriteString(" ")
	result.WriteString(d.Body.String())
	return result.String()
}

type ErrorDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The name of the error
	Fields      []Node          // Fields for the error type
	Annotations []Annotation    // Annotations applied to the error
}

func NewErrorDeclaration(name *TypeIdentifier, fields []Node, annotations []Annotation) *ErrorDeclaration {
	return &ErrorDeclaration{
		BaseNode:    BaseNode{Type: NodeErrorDeclaration},
		Name:        name,
		Fields:      fields,
		Annotations: annotations,
	}
}

func (d *ErrorDeclaration) String() string {
	result := ""

	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	result += "error " + d.Name.String()

	if len(d.Fields) > 0 {
		result += " {\n"
		for _, field := range d.Fields {
			result += "  " + field.String() + ",\n"
		}
		result += "}"
	}

	return result
}
