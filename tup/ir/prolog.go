package ir

import "strings"

// PrologConfig records the compiler-provided core import configuration for a
// semantic snapshot. The prolog itself is represented by ordinary module imports
// and local aliases.
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
