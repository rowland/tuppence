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

// TODO: Snapshot construction should append these diagnostics to IndexedProgram.Diagnostics
// once import loading owns diagnostic materialization.
// ImportResolutionDiagnostics reports unresolved or mismatched recorded imports.
func (p IndexedProgram) ImportResolutionDiagnostics() []Diagnostic {
	var diagnostics []Diagnostic
	for _, module := range p.Modules {
		for _, binding := range module.Imports {
			if binding.Path == "" {
				continue
			}
			candidates := p.Modules.ResolveImportCandidates(module.ID, binding.Path, p.ModuleResolution)
			resolved, ok := p.Modules.ResolveImport(module.ID, binding.Path, p.ModuleResolution)
			if !ok {
				diagnostics = append(diagnostics, NewUnresolvedImportDiagnostic(binding.Origin, binding.Path, candidates))
				continue
			}
			if resolved != binding.Target {
				diagnostics = append(diagnostics, Diagnostic{
					Severity:         SeverityError,
					Code:             "IMPORT_TARGET_MISMATCH",
					Message:          fmt.Sprintf("import %s resolves to module %d, not module %d", binding.Path, resolved, binding.Target),
					Origin:           binding.Origin,
					ImportPath:       binding.Path,
					ImportCandidates: candidates,
				})
			}
		}
	}
	return diagnostics
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
