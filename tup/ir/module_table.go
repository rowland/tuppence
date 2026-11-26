package ir

// ModuleTable is an ordered collection of modules indexed by identity.
type ModuleTable []Module

// Find returns the module with the requested identity.
func (t ModuleTable) Find(id ModuleID) (Module, bool) {
	for _, module := range t {
		if module.ID == id {
			return module, true
		}
	}
	return Module{}, false
}

func (t ModuleTable) ResolutionPath(id ModuleID) string {
	module, ok := t.Find(id)
	if !ok {
		return ""
	}
	return module.ResolutionPath()
}

// Valid reports whether every module is valid and IDs are unique.
func (t ModuleTable) Valid() bool {
	seen := make(map[ModuleID]struct{}, len(t))
	seenNames := make(map[string]struct{}, len(t))
	for _, module := range t {
		if !module.Valid() {
			return false
		}
		if _, exists := seen[module.ID]; exists {
			return false
		}
		if _, exists := seenNames[module.Name]; exists {
			return false
		}
		seen[module.ID] = struct{}{}
		seenNames[module.Name] = struct{}{}
	}
	return true
}
