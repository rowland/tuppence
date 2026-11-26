package ir

// ImportAliasTable is an ordered collection of local imported-member aliases.
type ImportAliasTable []ImportAlias

// Valid reports whether aliases are valid and unique within each scope.
func (t ImportAliasTable) Valid() bool {
	type key struct {
		scope ScopeID
		name  string
	}
	seen := make(map[key]struct{}, len(t))
	for _, alias := range t {
		if !alias.Valid() {
			return false
		}
		pair := key{scope: alias.Scope, name: alias.LocalName}
		if _, exists := seen[pair]; exists {
			return false
		}
		seen[pair] = struct{}{}
	}
	return true
}
