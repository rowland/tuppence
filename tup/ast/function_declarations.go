package ast

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

type FunctionDeclaration struct {
	BaseNode
	Name           *FunctionIdentifier // The name of the function
	TypeParameters []*GenericTypeParam // Type parameters for generic functions
	Parameters     []Node              // Function parameters
	ReturnType     Node                // The return type
	Body           Node                // The function body
	Annotations    []Annotation        // Annotations applied to the function
	IsLocal        bool                // Whether this is a local function
}

func NewFunctionDeclaration(name *FunctionIdentifier, typeParams []*GenericTypeParam, params []Node, returnType Node, body Node, annotations []Annotation, isLocal bool) *FunctionDeclaration {
	return &FunctionDeclaration{
		BaseNode:       BaseNode{Type: NodeFunctionDeclaration},
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		Body:           body,
		Annotations:    annotations,
		IsLocal:        isLocal,
	}
}

func (d *FunctionDeclaration) String() string {
	result := ""
	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	if d.IsLocal {
		result += "local "
	}

	result += "fn " + d.Name.String()

	if len(d.TypeParameters) > 0 {
		result += "<"
		for i, param := range d.TypeParameters {
			if i > 0 {
				result += ", "
			}
			result += param.String()
		}
		result += ">"
	}

	result += "("
	for i, param := range d.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param.String()
	}
	result += ")"

	if d.ReturnType != nil {
		result += " -> " + d.ReturnType.String()
	}

	result += " " + d.Body.String()
	return result
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
