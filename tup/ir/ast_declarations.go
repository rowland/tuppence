package ir

import "github.com/rowland/tuppence/tup/ast"

// CollectedTypeDeclarations contains IR declarations and scope entries derived
// from parsed type declarations. Type bodies are resolved separately.
type CollectedTypeDeclarations struct {
	Declarations DeclarationTable
	Entries      []ScopeEntry
	Diagnostics  []Diagnostic
}

// CollectTypeDeclarations converts parsed type declaration names into semantic
// declarations without interpreting their right-hand-side type bodies.
func CollectTypeDeclarations(scope ScopeID, file FileID, first DeclID, declarations []*ast.TypeDeclaration) CollectedTypeDeclarations {
	result := CollectedTypeDeclarations{}
	next := uint32(first)
	for _, declaration := range declarations {
		if declaration == nil || declaration.LHS == nil || declaration.LHS.Name == nil {
			continue
		}
		next++
		name := declaration.LHS.Name
		origin := Origin{Primary: Span{File: file, Start: uint32(name.StartOffset), End: uint32(name.StartOffset + name.Length)}}
		irDeclaration := Declaration{ID: DeclID(next), Scope: scope, Name: name.String(), Kind: DeclType, Origin: origin}
		result.Declarations = append(result.Declarations, irDeclaration)
		result.Entries = append(result.Entries, ScopeEntry{Name: irDeclaration.Name, Declaration: irDeclaration.ID})
	}
	return result
}
