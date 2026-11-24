package ir

// TupleMember is one logical tuple field. An empty Name denotes an ordinal
// field; LogicalSlot preserves source/logical order independently of layout.
type TupleMember struct {
	Name        string
	LogicalSlot uint32
	Type        TypeID
}

// Valid reports whether the member references a semantic type.
func (m TupleMember) Valid() bool {
	return m.Type.Valid()
}

// TupleSchema describes the logical members of a tuple-shaped type or
// parameter list. Defaults, rest modes, and constructors are added later.
type TupleSchema struct {
	ID      TupleSchemaID
	Members []TupleMember
	Origin  Origin
}

// Valid reports whether the schema and its logical slots are valid.
func (s TupleSchema) Valid() bool {
	if !s.ID.Valid() || !s.Origin.Valid() {
		return false
	}
	seenSlots := make(map[uint32]struct{}, len(s.Members))
	for _, member := range s.Members {
		if !member.Valid() {
			return false
		}
		if _, exists := seenSlots[member.LogicalSlot]; exists {
			return false
		}
		seenSlots[member.LogicalSlot] = struct{}{}
	}
	return true
}
