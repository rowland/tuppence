package ir

import "strings"

// TypeReference is an unresolved source-level type name.
type TypeReference struct {
	ID     TypeReferenceID
	Scope  ScopeID
	Name   string
	Origin Origin
}

// Valid reports whether the reference has the required identity and source data.
func (r TypeReference) Valid() bool {
	return r.ID.Valid() && r.Scope.Valid() && strings.TrimSpace(r.Name) != "" && r.Origin.Valid()
}

// ResolvedTypeReference preserves the source reference and its resolved type.
type ResolvedTypeReference struct {
	Reference TypeReference
	Type      TypeID
}

// Valid reports whether the reference and resolved type are valid.
func (r ResolvedTypeReference) Valid() bool {
	return r.Reference.Valid() && r.Type.Valid()
}

// ResolveTypeReferences resolves references in input order and returns all
// diagnostics rather than stopping at the first failure.
func (p IndexedProgram) ResolveTypeReferences(references []TypeReference) ([]ResolvedTypeReference, []Diagnostic) {
	resolved := make([]ResolvedTypeReference, 0, len(references))
	diagnostics := make([]Diagnostic, 0)
	for _, reference := range references {
		if !reference.Valid() {
			diagnostics = append(diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_REFERENCE_INVALID", Message: "invalid type reference", Origin: reference.Origin})
			continue
		}
		typeID, ok, diagnostic := p.ResolveTypeNameWithDiagnostic(reference.Scope, reference.Name, reference.Origin)
		if !ok {
			diagnostics = append(diagnostics, diagnostic)
			continue
		}
		resolved = append(resolved, ResolvedTypeReference{Reference: reference, Type: typeID})
	}
	return resolved, diagnostics
}
