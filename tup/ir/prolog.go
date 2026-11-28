package ir

import "strings"

// PrologConfig records the compiler-provided core include configuration for a
// semantic snapshot. The selected exports and aliases occupy a synthetic
// lexical scope immediately outside each module's initial scope; the prolog is
// not a general blanket import.
type PrologConfig struct {
	Enabled      bool
	Modules      []ModuleID
	Architecture string
}

// Valid reports whether the configuration is internally consistent.
func (c PrologConfig) Valid() bool {
	if !c.Enabled {
		return len(c.Modules) == 0 && strings.TrimSpace(c.Architecture) == ""
	}
	if len(c.Modules) == 0 || strings.TrimSpace(c.Architecture) == "" {
		return false
	}
	seen := make(map[ModuleID]struct{}, len(c.Modules))
	for _, module := range c.Modules {
		if !module.Valid() {
			return false
		}
		if _, exists := seen[module]; exists {
			return false
		}
		seen[module] = struct{}{}
	}
	return true
}

func (c PrologConfig) Includes(module ModuleID) bool {
	if !c.Enabled {
		return false
	}
	for _, candidate := range c.Modules {
		if candidate == module {
			return true
		}
	}
	return false
}
