package ast

import (
	"fmt"
	"strings"
)

// type_declaration = type_declaration_lhs "=" type_declaration_rhs .

type TypeDeclaration struct {
	BaseNode
	LHS *TypeDeclarationLHS
	RHS TypeDeclarationRHS
}

func NewTypeDeclaration(
	lhs *TypeDeclarationLHS,
	rhs TypeDeclarationRHS,
) *TypeDeclaration {
	return &TypeDeclaration{
		BaseNode: BaseNode{Type: NodeTypeDeclaration},
		LHS:      lhs,
		RHS:      rhs,
	}
}

func (d *TypeDeclaration) String() string {
	return fmt.Sprintf("%s = %s", d.LHS, d.RHS)
}

// type_declaration_lhs = annotations type_identifier [ type_parameters ] .

type TypeDeclarationLHS struct {
	BaseNode
	Annotations    []Annotation
	Name           *TypeIdentifier
	TypeParameters *TypeParameters
}

func NewTypeDeclarationLHS(annotations []Annotation, name *TypeIdentifier, typeParameters *TypeParameters) *TypeDeclarationLHS {
	return &TypeDeclarationLHS{
		BaseNode:       BaseNode{Type: NodeTypeDeclarationLHS},
		Annotations:    annotations,
		Name:           name,
		TypeParameters: typeParameters,
	}
}

func (d *TypeDeclarationLHS) String() string {
	var result strings.Builder
	for _, a := range d.Annotations {
		fmt.Fprintf(&result, "%s\n", a.String())
	}
	result.WriteString(d.Name.String())
	if d.TypeParameters != nil {
		result.WriteString(d.TypeParameters.String())
	}
	return result.String()
}

// type_declaration_rhs = nilable_type
//                      | type_tuple
//                      | error_tuple
//                      | dynamic_array
//                      | fixed_size_array
//                      | union_type
//                      | union_declaration
//                      | enum_declaration
//                      | contract_declaration
//                      | type_reference .

type TypeDeclarationRHS interface {
	Node
	typeDeclarationRHSNode()
}

func (n NilableType) typeDeclarationRHSNode()      {}
func (n TypeTuple) typeDeclarationRHSNode()        {}
func (n ErrorTuple) typeDeclarationRHSNode()       {}
func (n DynamicArrayType) typeDeclarationRHSNode() {}
func (n FixedSizeArrayType) typeDeclarationRHSNode() {
}

func (n UnionType) typeDeclarationRHSNode()           {}
func (n UnionDeclaration) typeDeclarationRHSNode()    {}
func (n EnumDeclaration) typeDeclarationRHSNode()     {}
func (n ContractDeclaration) typeDeclarationRHSNode() {}
func (n TypeReference) typeDeclarationRHSNode()       {}
