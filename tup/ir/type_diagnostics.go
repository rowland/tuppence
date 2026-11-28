package ir

import "fmt"

// NewUnresolvedTypeDiagnostic creates a diagnostic for a missing type name.
func NewUnresolvedTypeDiagnostic(origin Origin, name string) Diagnostic {
	return Diagnostic{Severity: SeverityError, Code: "TYPE_NOT_FOUND", Message: fmt.Sprintf("type %q was not found", name), Origin: origin}
}

// NewAmbiguousTypeDiagnostic creates a diagnostic for a type name with several
// candidates in the nearest visible scope.
func NewAmbiguousTypeDiagnostic(origin Origin, name string) Diagnostic {
	return Diagnostic{Severity: SeverityError, Code: "TYPE_AMBIGUOUS", Message: fmt.Sprintf("type %q is ambiguous", name), Origin: origin}
}

// ResolveTypeNameWithDiagnostic resolves a type name and returns a diagnostic
// when no unique lexical candidate exists.
func (p IndexedProgram) ResolveTypeNameWithDiagnostic(scope ScopeID, name string, origin Origin) (TypeID, bool, Diagnostic) {
	matches := p.FindTypeCandidates(scope, name)
	if len(matches) == 1 {
		return matches[0], true, Diagnostic{}
	}
	if len(matches) == 0 {
		return 0, false, NewUnresolvedTypeDiagnostic(origin, name)
	}
	return 0, false, NewAmbiguousTypeDiagnostic(origin, name)
}
