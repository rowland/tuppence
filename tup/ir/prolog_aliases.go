package ir

// AddPrologAlias installs one compiler-provided alias into an initial scope.
// Declaration and module ownership are validated later by IndexedProgram.Valid.
func (p *IndexedProgram) AddPrologAlias(alias ImportAlias) bool {
	if p == nil || !p.Prolog.Enabled || !p.Prolog.Includes(alias.Module) {
		return false
	}
	scope, ok := p.Scopes.Find(alias.Scope)
	if !ok || !scope.Initial {
		return false
	}
	alias.Source = ImportAliasSourceProlog
	p.ImportAliases = append(p.ImportAliases, alias)
	return true
}
