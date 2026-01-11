package ast

// GenericTypeParam represents a type parameter for generic types and functions
type GenericTypeParam struct {
	BaseNode
	Name       string // The name of the type parameter
	Constraint Node   // Optional type constraint
}

// NewGenericTypeParam creates a new GenericTypeParam node
func NewGenericTypeParam(name string, constraint Node) *GenericTypeParam {
	return &GenericTypeParam{
		BaseNode:   BaseNode{Type: NodeGenericTypeParam},
		Name:       name,
		Constraint: constraint,
	}
}

// String returns a textual representation of the type parameter
func (t *GenericTypeParam) String() string {
	if t.Constraint != nil {
		return t.Name + ": " + t.Constraint.String()
	}
	return t.Name
}

// TypeDeclaration represents a type declaration
type TypeDeclaration struct {
	BaseNode
	Name           *TypeIdentifier     // The name of the type
	TypeParameters []*GenericTypeParam // Type parameters if this is a generic type
	Type           Node                // The type definition
	Annotations    []Annotation        // Annotations applied to the type
}

// NewTypeDeclaration creates a new TypeDeclaration node
func NewTypeDeclaration(name *TypeIdentifier, typeParams []*GenericTypeParam, typeNode Node, annotations []Annotation) *TypeDeclaration {
	return &TypeDeclaration{
		BaseNode:       BaseNode{Type: NodeTypeDeclaration},
		Name:           name,
		TypeParameters: typeParams,
		Type:           typeNode,
		Annotations:    annotations,
	}
}

// String returns a textual representation of the type declaration
func (d *TypeDeclaration) String() string {
	result := ""
	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}
	result += "type " + d.Name.String()

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

	result += " = " + d.Type.String()
	return result
}

// FunctionDeclaration represents a function declaration
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

// NewFunctionDeclaration creates a new FunctionDeclaration node
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

// String returns a textual representation of the function declaration
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

// ErrorDeclaration represents an error type declaration
type ErrorDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The name of the error
	Fields      []Node          // Fields for the error type
	Annotations []Annotation    // Annotations applied to the error
}

// NewErrorDeclaration creates a new ErrorDeclaration node
func NewErrorDeclaration(name *TypeIdentifier, fields []Node, annotations []Annotation) *ErrorDeclaration {
	return &ErrorDeclaration{
		BaseNode:    BaseNode{Type: NodeErrorDeclaration},
		Name:        name,
		Fields:      fields,
		Annotations: annotations,
	}
}

// String returns a textual representation of the error declaration
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
