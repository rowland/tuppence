package ir

import "strings"

// DiagnosticSeverity classifies a compiler diagnostic.
type DiagnosticSeverity uint8

const (
	SeverityInfo DiagnosticSeverity = iota + 1
	SeverityWarning
	SeverityError
)

// Diagnostic is a source-backed compiler diagnostic. Code is intended to be
// stable for tooling; Message is human-readable and may evolve.
type Diagnostic struct {
	Severity             DiagnosticSeverity
	Code                 string
	Message              string
	Origin               Origin
	Related              []OriginID
	CanonicalModule      ModuleID
	CanonicalDeclaration DeclID
	ImportPath           string
	ImportCandidates     []string
	ImportSearchRoots    []string
	ImportResolvedPath   string
}

// Valid reports whether the diagnostic has the required stable identity,
// message, source origin, and related-origin references.
func (d Diagnostic) Valid() bool {
	if d.Severity < SeverityInfo || d.Severity > SeverityError {
		return false
	}
	if strings.TrimSpace(d.Code) == "" || strings.TrimSpace(d.Message) == "" {
		return false
	}
	if !d.Origin.Valid() {
		return false
	}
	if d.CanonicalModule.Valid() != d.CanonicalDeclaration.Valid() {
		return false
	}
	if strings.TrimSpace(d.ImportPath) == "" && (len(d.ImportCandidates) != 0 || len(d.ImportSearchRoots) != 0 || d.ImportResolvedPath != "") {
		return false
	}
	seenRoots := make(map[string]struct{}, len(d.ImportSearchRoots))
	for _, root := range d.ImportSearchRoots {
		if strings.TrimSpace(root) == "" {
			return false
		}
		if _, exists := seenRoots[root]; exists {
			return false
		}
		seenRoots[root] = struct{}{}
	}
	seenCandidates := make(map[string]struct{}, len(d.ImportCandidates))
	for _, candidate := range d.ImportCandidates {
		if strings.TrimSpace(candidate) == "" {
			return false
		}
		if _, exists := seenCandidates[candidate]; exists {
			return false
		}
		seenCandidates[candidate] = struct{}{}
	}
	for _, id := range d.Related {
		if !id.Valid() {
			return false
		}
	}
	return true
}
