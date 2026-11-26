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
	for _, id := range d.Related {
		if !id.Valid() {
			return false
		}
	}
	return true
}
