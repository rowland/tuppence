package ir

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

// Function records source-level callable metadata. Parameter schemas and
// specialization keys are added separately.
type Function struct {
	ID          FunctionID
	Declaration DeclID
	Effect      EffectKind
	Generic     bool
	Origin      Origin
}

// Valid reports whether the function has the required identity and provenance.
func (f Function) Valid() bool {
	return f.ID.Valid() && f.Declaration.Valid() && f.Effect.Valid() && f.Origin.Valid()
}
