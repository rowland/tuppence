package ir

import (
	"fmt"
	"strings"
)

// NewUnresolvedImportDiagnostic creates a source-backed diagnostic containing
// the requested path and the exact candidates considered by resolution.
func NewUnresolvedImportDiagnostic(origin Origin, requested string, candidates []string) Diagnostic {
	copied := append([]string(nil), candidates...)
	return Diagnostic{
		Severity:         SeverityError,
		Code:             "IMPORT_NOT_FOUND",
		Message:          fmt.Sprintf("module import not found: %s", requested),
		Origin:           origin,
		ImportPath:       requested,
		ImportCandidates: copied,
	}
}

// FormatImportResolution returns a compact explanation of the attempted paths.
func FormatImportResolution(diagnostic Diagnostic) string {
	if diagnostic.ImportPath == "" {
		return ""
	}
	if len(diagnostic.ImportCandidates) == 0 {
		return fmt.Sprintf("import %s (no candidates)", diagnostic.ImportPath)
	}
	return fmt.Sprintf("import %s (tried: %s)", diagnostic.ImportPath, strings.Join(diagnostic.ImportCandidates, ", "))
}
