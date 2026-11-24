package ir

import (
	"fmt"
	"strings"
)

// Dump returns a stable, human-readable representation of the indexed tables.
// Tables and entries are emitted in their stored order.
func (p IndexedProgram) Dump() string {
	var out strings.Builder

	out.WriteString("modules:\n")
	for _, module := range p.Modules {
		fmt.Fprintf(&out, "  %d name=%s\n", module.ID, module.Name)
	}

	out.WriteString("types:\n")
	for _, typ := range p.Types {
		fmt.Fprintf(&out, "  %d kind=%d name=%s declaration=%d home_module=%d\n", typ.ID, typ.Kind, typ.Name, typ.Declaration, typ.HomeModule)
	}

	out.WriteString("scopes:\n")
	for _, scope := range p.Scopes {
		fmt.Fprintf(&out, "  %d module=%d parent=%d kind=%d entries=[", scope.ID, scope.Module, scope.Parent, scope.Kind)
		for i, entry := range scope.Entries {
			if i > 0 {
				out.WriteString(",")
			}
			fmt.Fprintf(&out, "%s:%d", entry.Name, entry.Declaration)
		}
		out.WriteString("]\n")
	}

	out.WriteString("declarations:\n")
	for _, declaration := range p.Declarations {
		fmt.Fprintf(&out, "  %d scope=%d kind=%d name=%s\n", declaration.ID, declaration.Scope, declaration.Kind, declaration.Name)
	}

	out.WriteString("bindings:\n")
	for _, binding := range p.Bindings {
		fmt.Fprintf(&out, "  %d scope=%d previous=%d name=%s type=%d value=%d mutable=%t\n", binding.ID, binding.Scope, binding.Previous, binding.Name, binding.TypeSlot, binding.Value, binding.Mutable)
	}

	out.WriteString("functions:\n")
	for _, function := range p.Functions {
		fmt.Fprintf(&out, "  %d declaration=%d effect=%d generic=%t\n", function.ID, function.Declaration, function.Effect, function.Generic)
	}

	out.WriteString("aliases:\n")
	for _, alias := range p.Aliases {
		fmt.Fprintf(&out, "  declaration=%d target=%d\n", alias.Declaration, alias.Target)
	}

	out.WriteString("diagnostics:\n")
	for _, diagnostic := range p.Diagnostics {
		fmt.Fprintf(&out, "  severity=%d code=%s message=%s\n", diagnostic.Severity, diagnostic.Code, diagnostic.Message)
	}

	return out.String()
}
