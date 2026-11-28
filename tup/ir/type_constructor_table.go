package ir

// TypeConstructorTable is an ordered collection of literal storage constructors.
type TypeConstructorTable []TypeConstructor

// Find returns the constructor for a type and storage kind.
func (t TypeConstructorTable) Find(typ TypeID, storage StorageKind) (TypeConstructor, bool) {
	for _, constructor := range t {
		if constructor.Type == typ && constructor.Storage == storage {
			return constructor, true
		}
	}
	return TypeConstructor{}, false
}

// Valid reports whether every constructor is valid and each pair is unique.
func (t TypeConstructorTable) Valid() bool {
	type key struct {
		typ     TypeID
		storage StorageKind
	}
	seen := make(map[key]struct{}, len(t))
	for _, constructor := range t {
		if !constructor.Valid() {
			return false
		}
		pair := key{typ: constructor.Type, storage: constructor.Storage}
		if _, exists := seen[pair]; exists {
			return false
		}
		seen[pair] = struct{}{}
	}
	return true
}
