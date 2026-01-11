package ast

import "strings"

// export_assignment = assignment_lhs ":" expression .

// ExportAssignment represents an assignment that is exported
type ExportAssignment struct {
	BaseNode
	Assignment Assignment // The assignment being exported
}

// NewExportAssignment creates a new ExportAssignment node
func NewExportAssignment(assignment Assignment) *ExportAssignment {
	return &ExportAssignment{
		BaseNode:   BaseNode{Type: NodeExportAssignment},
		Assignment: assignment,
	}
}

// String returns a textual representation of the export assignment
func (e *ExportAssignment) String() string {
	return strings.Replace(e.Assignment.String(), " = ", ": ", 1)
}

// export_function_declaration = annotations function_declaration_lhs ":" function_declaration_type block .

// ExportFunctionDeclaration represents a function declaration that is exported
type ExportFunctionDeclaration struct {
	BaseNode
	Function *FunctionDeclaration // The function declaration being exported
}

// NewExportFunctionDeclaration creates a new ExportFunctionDeclaration node
func NewExportFunctionDeclaration(function *FunctionDeclaration) *ExportFunctionDeclaration {
	return &ExportFunctionDeclaration{
		BaseNode: BaseNode{Type: NodeExportFunctionDeclaration},
		Function: function,
	}
}

// String returns a textual representation of the export function declaration
func (e *ExportFunctionDeclaration) String() string {
	return strings.Replace(e.Function.String(), " = ", ": ", 1)
}

// export_type_declaration = type_declaration_lhs ":" type_declaration_rhs .

// ExportTypeDeclaration represents a type declaration that is exported
type ExportTypeDeclaration struct {
	BaseNode
	Type TypeDeclaration // The type declaration being exported
}

// NewExportTypeDeclaration creates a new ExportTypeDeclaration node
func NewExportTypeDeclaration(typeDecl TypeDeclaration) *ExportTypeDeclaration {
	return &ExportTypeDeclaration{
		BaseNode: BaseNode{Type: NodeExportTypeDeclaration},
		Type:     typeDecl,
	}
}

// String returns a textual representation of the export type declaration
func (e *ExportTypeDeclaration) String() string {
	return strings.Replace(e.Type.String(), " = ", ": ", 1)
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
		BaseNode:    BaseNode{Type: NodeTypeQualifiedDeclaration},
		TypeName:    typeName,
		Declaration: declaration,
	}
}

// String returns a textual representation of the type-qualified declaration
func (t *TypeQualifiedDeclaration) String() string {
	return t.TypeName.String() + "." + t.Declaration.String()
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
		BaseNode: BaseNode{Type: NodeTypeQualifiedFunctionDeclaration},
		TypeName: typeName,
		Function: function,
	}
}

// String returns a textual representation of the type-qualified function declaration
func (t *TypeQualifiedFunctionDeclaration) String() string {
	return t.TypeName.String() + "." + t.Function.String()
}

// export_type_qualified_declaration = type_identifier "." identifier ":" expression .

// ExportTypeQualifiedDeclaration represents an exported type-qualified declaration
type ExportTypeQualifiedDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedDeclaration // The type-qualified declaration being exported
}

// NewExportTypeQualifiedDeclaration creates a new ExportTypeQualifiedDeclaration node
func NewExportTypeQualifiedDeclaration(declaration *TypeQualifiedDeclaration) *ExportTypeQualifiedDeclaration {
	return &ExportTypeQualifiedDeclaration{
		BaseNode:    BaseNode{Type: NodeExportTypeQualifiedDeclaration},
		Declaration: declaration,
	}
}

// String returns a textual representation of the export type-qualified declaration
func (e *ExportTypeQualifiedDeclaration) String() string {
	return strings.Replace(e.Declaration.String(), " = ", ": ", 1)
}

// export_type_qualified_function_declaration = annotations type_identifier "." function_declaration_lhs ":" function_declaration_type block .

// ExportTypeQualifiedFunctionDeclaration represents an exported type-qualified function declaration
type ExportTypeQualifiedFunctionDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedFunctionDeclaration // The type-qualified function declaration being exported
}

// NewExportTypeQualifiedFunctionDeclaration creates a new ExportTypeQualifiedFunctionDeclaration node
func NewExportTypeQualifiedFunctionDeclaration(declaration *TypeQualifiedFunctionDeclaration) *ExportTypeQualifiedFunctionDeclaration {
	return &ExportTypeQualifiedFunctionDeclaration{
		BaseNode:    BaseNode{Type: NodeExportTypeQualifiedFunctionDeclaration},
		Declaration: declaration,
	}
}

// String returns a textual representation of the export type-qualified function declaration
func (e *ExportTypeQualifiedFunctionDeclaration) String() string {
	return strings.Replace(e.Declaration.String(), " = ", ": ", 1)
}
