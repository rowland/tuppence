package ast

// Export related node types
const (
	NodeExportAssignment                       NodeType = "ExportAssignment"
	NodeExportDeclaration                      NodeType = "ExportDeclaration"
	NodeExportFunctionDeclaration              NodeType = "ExportFunctionDeclaration"
	NodeExportTypeDeclaration                  NodeType = "ExportTypeDeclaration"
	NodeExportTypeQualifiedDeclaration         NodeType = "ExportTypeQualifiedDeclaration"
	NodeExportTypeQualifiedFunctionDeclaration NodeType = "ExportTypeQualifiedFunctionDeclaration"
	NodeFunctionDeclarationType                NodeType = "FunctionDeclarationType"
	NodeFunctionParameterTypes                 NodeType = "FunctionParameterTypes"
	NodeFunctionTypeDeclaration                NodeType = "FunctionTypeDeclaration"
	NodeTypeQualifiedDeclaration               NodeType = "TypeQualifiedDeclaration"
	NodeTypeQualifiedFunctionDeclaration       NodeType = "TypeQualifiedFunctionDeclaration"
)

// ExportAssignment represents an assignment that is exported
type ExportAssignment struct {
	BaseNode
	Assignment Node // The assignment being exported
}

// NewExportAssignment creates a new ExportAssignment node
func NewExportAssignment(assignment Node) *ExportAssignment {
	return &ExportAssignment{
		BaseNode:   BaseNode{NodeType: NodeExportAssignment},
		Assignment: assignment,
	}
}

// String returns a textual representation of the export assignment
func (e *ExportAssignment) String() string {
	return "export " + e.Assignment.String()
}

// Children returns the child nodes
func (e *ExportAssignment) Children() []Node {
	return []Node{e.Assignment}
}

// ExportDeclaration represents a general declaration being exported
type ExportDeclaration struct {
	BaseNode
	Declaration Node // The declaration being exported
}

// NewExportDeclaration creates a new ExportDeclaration node
func NewExportDeclaration(declaration Node) *ExportDeclaration {
	return &ExportDeclaration{
		BaseNode:    BaseNode{NodeType: NodeExportDeclaration},
		Declaration: declaration,
	}
}

// String returns a textual representation of the export declaration
func (e *ExportDeclaration) String() string {
	return "export " + e.Declaration.String()
}

// Children returns the child nodes
func (e *ExportDeclaration) Children() []Node {
	return []Node{e.Declaration}
}

// ExportFunctionDeclaration represents a function declaration that is exported
type ExportFunctionDeclaration struct {
	BaseNode
	Function *FunctionDeclaration // The function declaration being exported
}

// NewExportFunctionDeclaration creates a new ExportFunctionDeclaration node
func NewExportFunctionDeclaration(function *FunctionDeclaration) *ExportFunctionDeclaration {
	return &ExportFunctionDeclaration{
		BaseNode: BaseNode{NodeType: NodeExportFunctionDeclaration},
		Function: function,
	}
}

// String returns a textual representation of the export function declaration
func (e *ExportFunctionDeclaration) String() string {
	return "export " + e.Function.String()
}

// Children returns the child nodes
func (e *ExportFunctionDeclaration) Children() []Node {
	return []Node{e.Function}
}

// ExportTypeDeclaration represents a type declaration that is exported
type ExportTypeDeclaration struct {
	BaseNode
	Type Node // The type declaration being exported
}

// NewExportTypeDeclaration creates a new ExportTypeDeclaration node
func NewExportTypeDeclaration(typeDecl Node) *ExportTypeDeclaration {
	return &ExportTypeDeclaration{
		BaseNode: BaseNode{NodeType: NodeExportTypeDeclaration},
		Type:     typeDecl,
	}
}

// String returns a textual representation of the export type declaration
func (e *ExportTypeDeclaration) String() string {
	return "export " + e.Type.String()
}

// Children returns the child nodes
func (e *ExportTypeDeclaration) Children() []Node {
	return []Node{e.Type}
}

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
		BaseNode:       BaseNode{NodeType: NodeFunctionDeclarationType},
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

// Children returns the child nodes
func (f *FunctionDeclarationType) Children() []Node {
	var children []Node
	children = append(children, f.Parameters)

	if f.ReturnType != nil {
		children = append(children, f.ReturnType)
	}

	return children
}

// FunctionParameterTypes represents the parameter types in a function declaration
type FunctionParameterTypes struct {
	BaseNode
	Parameters []Node // The function parameter types
}

// NewFunctionParameterTypes creates a new FunctionParameterTypes node
func NewFunctionParameterTypes(parameters []Node) *FunctionParameterTypes {
	return &FunctionParameterTypes{
		BaseNode:   BaseNode{NodeType: NodeFunctionParameterTypes},
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

// Children returns the child nodes
func (f *FunctionParameterTypes) Children() []Node {
	return f.Parameters
}

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
		BaseNode:   BaseNode{NodeType: NodeFunctionTypeDeclaration},
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

// Children returns the child nodes
func (f *FunctionTypeDeclaration) Children() []Node {
	children := make([]Node, 0, len(f.TypeParams)+2)
	children = append(children, f.Name)

	for _, param := range f.TypeParams {
		children = append(children, param)
	}

	children = append(children, f.Type)
	return children
}

// TypeQualifiedDeclaration represents a type-qualified declaration
type TypeQualifiedDeclaration struct {
	BaseNode
	TypeName    *TypeIdentifier // The type being qualified
	Declaration Node            // The declaration being qualified
}

// NewTypeQualifiedDeclaration creates a new TypeQualifiedDeclaration node
func NewTypeQualifiedDeclaration(typeName *TypeIdentifier, declaration Node) *TypeQualifiedDeclaration {
	return &TypeQualifiedDeclaration{
		BaseNode:    BaseNode{NodeType: NodeTypeQualifiedDeclaration},
		TypeName:    typeName,
		Declaration: declaration,
	}
}

// String returns a textual representation of the type-qualified declaration
func (t *TypeQualifiedDeclaration) String() string {
	return t.TypeName.String() + "." + t.Declaration.String()
}

// Children returns the child nodes
func (t *TypeQualifiedDeclaration) Children() []Node {
	return []Node{t.TypeName, t.Declaration}
}

// TypeQualifiedFunctionDeclaration represents a function declaration for a specific type
type TypeQualifiedFunctionDeclaration struct {
	BaseNode
	TypeName *TypeIdentifier      // The type being qualified
	Function *FunctionDeclaration // The function declaration
}

// NewTypeQualifiedFunctionDeclaration creates a new TypeQualifiedFunctionDeclaration node
func NewTypeQualifiedFunctionDeclaration(typeName *TypeIdentifier, function *FunctionDeclaration) *TypeQualifiedFunctionDeclaration {
	return &TypeQualifiedFunctionDeclaration{
		BaseNode: BaseNode{NodeType: NodeTypeQualifiedFunctionDeclaration},
		TypeName: typeName,
		Function: function,
	}
}

// String returns a textual representation of the type-qualified function declaration
func (t *TypeQualifiedFunctionDeclaration) String() string {
	return t.TypeName.String() + "." + t.Function.String()
}

// Children returns the child nodes
func (t *TypeQualifiedFunctionDeclaration) Children() []Node {
	return []Node{t.TypeName, t.Function}
}

// ExportTypeQualifiedDeclaration represents an exported type-qualified declaration
type ExportTypeQualifiedDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedDeclaration // The type-qualified declaration being exported
}

// NewExportTypeQualifiedDeclaration creates a new ExportTypeQualifiedDeclaration node
func NewExportTypeQualifiedDeclaration(declaration *TypeQualifiedDeclaration) *ExportTypeQualifiedDeclaration {
	return &ExportTypeQualifiedDeclaration{
		BaseNode:    BaseNode{NodeType: NodeExportTypeQualifiedDeclaration},
		Declaration: declaration,
	}
}

// String returns a textual representation of the export type-qualified declaration
func (e *ExportTypeQualifiedDeclaration) String() string {
	return "export " + e.Declaration.String()
}

// Children returns the child nodes
func (e *ExportTypeQualifiedDeclaration) Children() []Node {
	return []Node{e.Declaration}
}

// ExportTypeQualifiedFunctionDeclaration represents an exported type-qualified function declaration
type ExportTypeQualifiedFunctionDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedFunctionDeclaration // The type-qualified function declaration being exported
}

// NewExportTypeQualifiedFunctionDeclaration creates a new ExportTypeQualifiedFunctionDeclaration node
func NewExportTypeQualifiedFunctionDeclaration(declaration *TypeQualifiedFunctionDeclaration) *ExportTypeQualifiedFunctionDeclaration {
	return &ExportTypeQualifiedFunctionDeclaration{
		BaseNode:    BaseNode{NodeType: NodeExportTypeQualifiedFunctionDeclaration},
		Declaration: declaration,
	}
}

// String returns a textual representation of the export type-qualified function declaration
func (e *ExportTypeQualifiedFunctionDeclaration) String() string {
	return "export " + e.Declaration.String()
}

// Children returns the child nodes
func (e *ExportTypeQualifiedFunctionDeclaration) Children() []Node {
	return []Node{e.Declaration}
}
