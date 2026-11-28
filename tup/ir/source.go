package ir

// Span identifies a half-open byte range in an immutable source snapshot.
type Span struct {
	File       FileID
	Start, End uint32
}

// Valid reports whether the span has a valid file and ordered offsets.
func (s Span) Valid() bool {
	return s.File.Valid() && s.Start <= s.End
}

// Origin preserves the source location of an entity and, for synthesized
// entities, the origins from which it was derived.
type Origin struct {
	Primary     Span
	DerivedFrom []OriginID
}

// Valid reports whether the origin has a valid primary span and no invalid
// provenance IDs.
func (o Origin) Valid() bool {
	if !o.Primary.Valid() {
		return false
	}
	for _, id := range o.DerivedFrom {
		if !id.Valid() {
			return false
		}
	}
	return true
}
