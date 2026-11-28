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
	seenIDs := make(map[ModuleID]struct{}, len(t))
	seenPaths := make(map[string]struct{}, len(t))
	for _, module := range t {
		if !module.Valid() {
			return false
		}
		if _, exists := seenIDs[module.ID]; exists {
			return false
		}
		modulePath := module.ResolutionPath()
		if _, exists := seenPaths[modulePath]; exists {
			return false
		}
		seenIDs[module.ID] = struct{}{}
		seenPaths[modulePath] = struct{}{}
	}
	return true
}
