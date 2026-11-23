package ir

// BindingTable is an ordered collection of binding versions indexed by identity.
type BindingTable []Binding

// Find returns the binding with the requested identity.
func (t BindingTable) Find(id BindingID) (Binding, bool) {
	for _, binding := range t {
		if binding.ID == id {
			return binding, true
		}
	}
	return Binding{}, false
}

// Valid reports whether every binding is valid and IDs are unique.
func (t BindingTable) Valid() bool {
	seen := make(map[BindingID]struct{}, len(t))
	for _, binding := range t {
		if !binding.Valid() {
			return false
		}
		if _, exists := seen[binding.ID]; exists {
			return false
		}
		seen[binding.ID] = struct{}{}
	}
	return true
}
