package ir

// AddPrologAliases installs a batch atomically after declarations are indexed.
func (p *IndexedProgram) AddPrologAliases(aliases []ImportAlias) bool {
	if p == nil {
		return false
	}
	originalLength := len(p.ImportAliases)
	for _, alias := range aliases {
		if !p.AddPrologAlias(alias) {
			p.ImportAliases = p.ImportAliases[:originalLength]
			return false
		}
	}
	return true
}

// Declaration and module ownership are validated later by IndexedProgram.Valid.
func (p *IndexedProgram) AddPrologAlias(alias ImportAlias) bool {
	if p == nil || !p.Prolog.Enabled || !p.Prolog.Includes(alias.Module) {
		return false
	}
	scope, ok := p.Scopes.Find(alias.Scope)
	if !ok || scope.Kind != ScopeProlog {
		return false
	}
	alias.Source = ImportAliasSourceProlog
	p.ImportAliases = append(p.ImportAliases, alias)
	return true
}
