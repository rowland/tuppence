package ir

import "strings"

// PrologConfig records the compiler-provided core include configuration for a
// semantic snapshot. The selected exports and aliases behave as declarations in
// the module's initial scope; the prolog is not a general blanket import.
type PrologConfig struct {
	Enabled      bool
	Module       ModuleID
	Architecture string
}

// Valid reports whether the configuration is internally consistent.
func (c PrologConfig) Valid() bool {
	if !c.Enabled {
		return !c.Module.Valid() && strings.TrimSpace(c.Architecture) == ""
	}
	return c.Module.Valid() && strings.TrimSpace(c.Architecture) != ""
}
