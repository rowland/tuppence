package ir

import "strings"

// TypeKind identifies the initial semantic type categories.
type TypeKind uint8

const (
	TypePrimitive TypeKind = iota + 1
	TypeNominal
)

// Valid reports whether the type kind is recognized.
func (k TypeKind) Valid() bool {
	return k >= TypePrimitive && k <= TypeNominal
}

// Type is the minimal named semantic type record. Structural types and type
// variables will be added separately.
type Type struct {
	ID     TypeID
	Name   string
	Kind   TypeKind
	Origin Origin
}

// Valid reports whether the type has the required identity and source data.
func (t Type) Valid() bool {
	return t.ID.Valid() && strings.TrimSpace(t.Name) != "" && t.Kind.Valid() && t.Origin.Valid()
}
