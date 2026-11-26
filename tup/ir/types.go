package ir

import "strings"

// TypeKind identifies semantic type categories. Storage characteristics are
// represented separately by StorageKind.
type TypeKind uint8

const (
	TypeNominal TypeKind = iota + 1
	TypeArray
	TypeTuple
	TypeUnion
)

// Valid reports whether the type kind is recognized.
func (k TypeKind) Valid() bool {
	return k >= TypeNominal && k <= TypeUnion
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
	StorageBytes
)

// Valid reports whether the storage kind is recognized.
func (k StorageKind) Valid() bool {
	return k >= StorageI8 && k <= StorageBytes
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
	Elements    []TypeID
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
		return strings.TrimSpace(t.Name) == "" && t.Element.Valid() && t.Elements == nil && (!t.FixedLength && t.Length == 0 || t.FixedLength) && !t.Declaration.Valid() && !t.HomeModule.Valid()
	}
	if t.Kind == TypeTuple {
		return strings.TrimSpace(t.Name) == "" && t.Elements != nil && !t.Element.Valid() && !t.FixedLength && t.Length == 0 && !t.Declaration.Valid() && !t.HomeModule.Valid()
	}
	if t.Kind == TypeUnion {
		return strings.TrimSpace(t.Name) == "" && len(t.Elements) > 0 && !t.Element.Valid() && !t.FixedLength && t.Length == 0 && !t.Declaration.Valid() && !t.HomeModule.Valid()
	}
	return strings.TrimSpace(t.Name) != "" && t.Element == 0 && t.Elements == nil && !t.FixedLength && t.Length == 0 && t.Declaration.Valid() && t.HomeModule.Valid()
}
