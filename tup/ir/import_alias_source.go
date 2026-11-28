package ir

// ImportAliasSource identifies how an imported member alias entered scope.
type ImportAliasSource uint8

const (
	// ImportAliasSourceUser is an alias written by source code.
	ImportAliasSourceUser ImportAliasSource = iota
	// ImportAliasSourceProlog is an alias introduced by the compiler-provided core prolog.
	ImportAliasSourceProlog
)

// Valid reports whether the source kind is known.
func (s ImportAliasSource) Valid() bool {
	return s == ImportAliasSourceUser || s == ImportAliasSourceProlog
}
