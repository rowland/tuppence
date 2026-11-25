package ir

import "strings"

// EffectKind identifies whether a function may perform effects.
type EffectKind uint8

const (
	EffectPure EffectKind = iota + 1
	Effectful
)

// Valid reports whether the effect kind is recognized.
func (e EffectKind) Valid() bool {
	return e >= EffectPure && e <= Effectful
}

// FunctionParameter describes one named parameter in declaration order.
type FunctionParameter struct {
	Name string
	Type TypeID
}

// Valid reports whether the parameter has a name and type.
func (p FunctionParameter) Valid() bool {
	return strings.TrimSpace(p.Name) != "" && p.Type.Valid()
}

// Function records source-level callable metadata. Parameter schemas and
// specialization keys are added separately.
type Function struct {
	ID          FunctionID
	Declaration DeclID
	Parameters  []FunctionParameter
	Result      TypeID
	Effect      EffectKind
	Generic     bool
	Origin      Origin
}

// Valid reports whether the function has the required identity and provenance.
func (f Function) Valid() bool {
	return f.ID.Valid() && f.Declaration.Valid() && f.Effect.Valid() && f.Origin.Valid()
}
