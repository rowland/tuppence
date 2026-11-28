package ir

import "github.com/rowland/tuppence/tup/ast"

// NewTypeReferenceFromAST converts a parsed local type reference into the IR
// form while preserving its source span. Qualified names remain intact for a
// later qualified-type resolution pass.
func NewTypeReferenceFromAST(id TypeReferenceID, scope ScopeID, file FileID, reference *ast.TypeReference) TypeReference {
	if reference == nil {
		return TypeReference{}
	}
	return TypeReference{
		ID:     id,
		Scope:  scope,
		Name:   reference.String(),
		Origin: Origin{Primary: Span{File: file, Start: uint32(reference.StartOffset), End: uint32(reference.StartOffset + reference.Length)}},
	}
}
