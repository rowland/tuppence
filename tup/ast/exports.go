package ast

import "strings"

// export_assignment = assignment_lhs ":" expression .

type ExportAssignment struct {
	BaseNode
	Assignment Assignment // The assignment being exported
}

func NewExportAssignment(assignment Assignment) *ExportAssignment {
	return &ExportAssignment{
		BaseNode:   BaseNode{Type: NodeExportAssignment},
		Assignment: assignment,
	}
}

func (e *ExportAssignment) String() string {
	return strings.Replace(e.Assignment.String(), " = ", ": ", 1)
}

// export_function_declaration = annotations function_declaration_lhs ":" function_declaration_type block .

type ExportFunctionDeclaration struct {
	BaseNode
	Function *FunctionDeclaration // The function declaration being exported
}

func NewExportFunctionDeclaration(function *FunctionDeclaration) *ExportFunctionDeclaration {
	return &ExportFunctionDeclaration{
		BaseNode: BaseNode{Type: NodeExportFunctionDeclaration},
		Function: function,
	}
}

func (e *ExportFunctionDeclaration) String() string {
	return strings.Replace(e.Function.String(), " = ", ": ", 1)
}

// export_function_type_declaration = function_type_declaration_lhs ":" function_type .

type ExportFunctionTypeDeclaration struct {
	BaseNode
	FunctionType *FunctionTypeDeclaration
}

func NewExportFunctionTypeDeclaration(functionType *FunctionTypeDeclaration) *ExportFunctionTypeDeclaration {
	return &ExportFunctionTypeDeclaration{
		BaseNode:     BaseNode{Type: NodeExportFunctionTypeDeclaration},
		FunctionType: functionType,
	}
}

func (e *ExportFunctionTypeDeclaration) String() string {
	return strings.Replace(e.FunctionType.String(), " = ", ": ", 1)
}

// export_type_declaration = type_declaration_lhs ":" type_declaration_rhs .

type ExportTypeDeclaration struct {
	BaseNode
	Type TypeDeclaration // The type declaration being exported
}

func NewExportTypeDeclaration(typeDecl TypeDeclaration) *ExportTypeDeclaration {
	return &ExportTypeDeclaration{
		BaseNode: BaseNode{Type: NodeExportTypeDeclaration},
		Type:     typeDecl,
	}
}

func (e *ExportTypeDeclaration) String() string {
	return strings.Replace(e.Type.String(), " = ", ": ", 1)
}

type TypeQualifiedDeclaration struct {
	BaseNode
	TypeName    *TypeIdentifier // The type being qualified
	Declaration Node            // The declaration being qualified
}

func NewTypeQualifiedDeclaration(typeName *TypeIdentifier, declaration Node) *TypeQualifiedDeclaration {
	return &TypeQualifiedDeclaration{
		BaseNode:    BaseNode{Type: NodeTypeQualifiedDeclaration},
		TypeName:    typeName,
		Declaration: declaration,
	}
}

func (t *TypeQualifiedDeclaration) String() string {
	return t.TypeName.String() + "." + t.Declaration.String()
}

type TypeQualifiedFunctionDeclaration struct {
	BaseNode
	TypeName *TypeIdentifier      // The type being qualified
	Function *FunctionDeclaration // The function declaration
}

func NewTypeQualifiedFunctionDeclaration(typeName *TypeIdentifier, function *FunctionDeclaration) *TypeQualifiedFunctionDeclaration {
	return &TypeQualifiedFunctionDeclaration{
		BaseNode: BaseNode{Type: NodeTypeQualifiedFunctionDeclaration},
		TypeName: typeName,
		Function: function,
	}
}

func (t *TypeQualifiedFunctionDeclaration) String() string {
	return t.TypeName.String() + "." + t.Function.String()
}

// export_type_qualified_declaration = type_identifier "." identifier ":" expression .

type ExportTypeQualifiedDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedDeclaration // The type-qualified declaration being exported
}

func NewExportTypeQualifiedDeclaration(declaration *TypeQualifiedDeclaration) *ExportTypeQualifiedDeclaration {
	return &ExportTypeQualifiedDeclaration{
		BaseNode:    BaseNode{Type: NodeExportTypeQualifiedDeclaration},
		Declaration: declaration,
	}
}

func (e *ExportTypeQualifiedDeclaration) String() string {
	return strings.Replace(e.Declaration.String(), " = ", ": ", 1)
}

// export_type_qualified_function_declaration = annotations type_identifier "." function_declaration_lhs ":" function_declaration_type block .

type ExportTypeQualifiedFunctionDeclaration struct {
	BaseNode
	Declaration *TypeQualifiedFunctionDeclaration // The type-qualified function declaration being exported
}

func NewExportTypeQualifiedFunctionDeclaration(declaration *TypeQualifiedFunctionDeclaration) *ExportTypeQualifiedFunctionDeclaration {
	return &ExportTypeQualifiedFunctionDeclaration{
		BaseNode:    BaseNode{Type: NodeExportTypeQualifiedFunctionDeclaration},
		Declaration: declaration,
	}
}

func (e *ExportTypeQualifiedFunctionDeclaration) String() string {
	return strings.Replace(e.Declaration.String(), " = ", ": ", 1)
}
