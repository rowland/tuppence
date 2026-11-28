package ir

import "github.com/rowland/tuppence/tup/ast"

// CollectedTypeDeclarations contains IR declarations and scope entries derived
// from parsed type declarations. Type bodies are resolved separately.
type CollectedTypeDeclarations struct {
	Declarations DeclarationTable
	Entries      []ScopeEntry
	Diagnostics  []Diagnostic
}

// InstallCollectedTypeDeclarations adds collected declarations to a module's
// initial scope atomically.
func (p *IndexedProgram) InstallCollectedTypeDeclarations(module ModuleID, collected CollectedTypeDeclarations) bool {
	if p == nil || !module.Valid() || len(collected.Diagnostics) != 0 {
		return false
	}
	var scope *Scope
	for i := range p.Scopes {
		candidate := &p.Scopes[i]
		if candidate.Module == module && candidate.Kind == ScopeModule && candidate.Initial {
			scope = candidate
			break
		}
	}
	if scope == nil || len(collected.Declarations) != len(collected.Entries) {
		return false
	}
	seenNames := make(map[string]struct{}, len(scope.Entries)+len(collected.Entries))
	for _, entry := range scope.Entries {
		seenNames[entry.Name] = struct{}{}
	}
	for i, declaration := range collected.Declarations {
		entry := collected.Entries[i]
		if declaration.Scope != scope.ID || entry.Declaration != declaration.ID || entry.Name != declaration.Name || declaration.ID == 0 {
			return false
		}
		if _, exists := seenNames[entry.Name]; exists {
			return false
		}
		if _, exists := p.Declarations.Find(declaration.ID); exists {
			return false
		}
		seenNames[entry.Name] = struct{}{}
	}
	p.Declarations = append(p.Declarations, collected.Declarations...)
	scope.Entries = append(scope.Entries, collected.Entries...)
	return true
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
