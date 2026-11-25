package ir

import "strings"

// TypeKind identifies semantic type categories. Storage characteristics are
// represented separately by StorageKind.
type TypeKind uint8

const (
	TypeNominal TypeKind = iota + 1
	TypeArray
)

// Valid reports whether the type kind is recognized.
func (k TypeKind) Valid() bool {
	return k >= TypeNominal && k <= TypeArray
}

// StorageKind identifies backend-relevant representation characteristics, not
// language-level behavior or standard-library operations.
type StorageKind uint8

const (
	StorageI8 StorageKind = iota + 1
	StorageI16
	StorageI32
	StorageI64
	StorageU8
	StorageU16
	StorageU32
	StorageU64
	StorageF16
	StorageF32
	StorageF64
	StorageV128
	StoragePTR
)

// Valid reports whether the storage kind is recognized.
func (k StorageKind) Valid() bool {
	return k >= StorageI8 && k <= StoragePTR
}

// StorageType describes representation only. Arithmetic, comparison, and
// other behavior belong to Tuppence declarations and library functions.
type StorageType struct {
	Kind StorageKind
}

// Valid reports whether the storage descriptor has a recognized kind.
func (s StorageType) Valid() bool {
	return s.Kind.Valid()
}

// Type is the minimal semantic type record. A nominal type gets its identity
// from a declaration and home module; structural types will be added later.
type Type struct {
	ID          TypeID
	Name        string
	Kind        TypeKind
	Declaration DeclID
	HomeModule  ModuleID
	Element     TypeID
	FixedLength bool
	Length      uint64
	Origin      Origin
}

// Valid reports whether the type has the required identity and source data.
func (t Type) Valid() bool {
	if !t.ID.Valid() || !t.Kind.Valid() || !t.Origin.Valid() {
		return false
	}
	if t.Kind == TypeArray {
		return t.Element.Valid() && !t.Declaration.Valid() && !t.HomeModule.Valid()
	}
	return strings.TrimSpace(t.Name) != "" && t.Declaration.Valid() && t.HomeModule.Valid()
}
