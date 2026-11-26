package ir

// TypeConstructor describes one semantic type's admissible literal storage
// representation. It is separate from the semantic type and storage descriptor.
type TypeConstructor struct {
	Type    TypeID
	Storage StorageKind
	Origin  Origin
}

// Valid reports whether the constructor has valid identity and provenance.
func (c TypeConstructor) Valid() bool {
	return c.Type.Valid() && c.Storage.Valid() && c.Origin.Valid()
}
