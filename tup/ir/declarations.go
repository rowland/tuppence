package ir

import "strings"

// DeclKind identifies the declaration category visible to semantic analysis.
type DeclKind uint8

const (
	DeclValue DeclKind = iota + 1
	DeclFunction
	DeclType
	DeclParameter
)

// Valid reports whether the declaration kind is recognized.
func (k DeclKind) Valid() bool {
	return k >= DeclValue && k <= DeclParameter
}

// Declaration records the source identity and lexical owner of one declaration.
type Declaration struct {
	ID        DeclID
	Scope     ScopeID
	Name      string
	Kind      DeclKind
	OwnerType DeclID
	Origin    Origin
}

// Valid reports whether the declaration has the required identity and source data.
func (d Declaration) Valid() bool {
	return d.ID.Valid() && d.Scope.Valid() && strings.TrimSpace(d.Name) != "" && d.Kind.Valid() && d.Origin.Valid()
}
