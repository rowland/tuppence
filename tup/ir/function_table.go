package ir

// FunctionTable is an ordered collection of functions indexed by identity.
type FunctionTable []Function

// Find returns the function with the requested identity.
func (t FunctionTable) Find(id FunctionID) (Function, bool) {
	for _, function := range t {
		if function.ID == id {
			return function, true
		}
	}
	return Function{}, false
}

// Valid reports whether every function is valid and IDs are unique.
func (t FunctionTable) Valid() bool {
	seen := make(map[FunctionID]struct{}, len(t))
	for _, function := range t {
		if !function.Valid() {
			return false
		}
		if _, exists := seen[function.ID]; exists {
			return false
		}
		seen[function.ID] = struct{}{}
	}
	return true
}
