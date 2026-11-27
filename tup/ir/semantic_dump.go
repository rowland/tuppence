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
		fmt.Fprintf(&out, "  %d name=%s path=%s files=%v imports=[", module.ID, module.Name, module.Path, module.Files)
		for i, importBinding := range module.Imports {
			if i > 0 {
				out.WriteString(",")
			}
			fmt.Fprintf(&out, "%s:%s:%d", importBinding.LocalName, importBinding.Path, importBinding.Target)
		}
		out.WriteString("]\n")
	}

	out.WriteString("source_files:\n")
	for _, file := range p.SourceFiles {
		fmt.Fprintf(&out, "  %d path=%s module=%s constraints=%v\n", file.ID, file.Path, file.Module.ModuleName, file.Module.Constraints)
	}

	out.WriteString("primary_module:\n")
	fmt.Fprintf(&out, "  id=%d\n", p.PrimaryModule)

	out.WriteString("types:\n")
	for _, typ := range p.Types {
		fmt.Fprintf(&out, "  %d kind=%d name=%s declaration=%d home_module=%d element=%d elements=%v fixed_length=%t length=%d\n", typ.ID, typ.Kind, typ.Name, typ.Declaration, typ.HomeModule, typ.Element, typ.Elements, typ.FixedLength, typ.Length)
	}

	out.WriteString("constructors:\n")
	for _, constructor := range p.Constructors {
		fmt.Fprintf(&out, "  type=%d storage=%d\n", constructor.Type, constructor.Storage)
	}

	out.WriteString("scopes:\n")
	for _, scope := range p.Scopes {
		fmt.Fprintf(&out, "  %d module=%d parent=%d kind=%d initial=%t origin_file=%d entries=[", scope.ID, scope.Module, scope.Parent, scope.Kind, scope.Initial, scope.Origin.Primary.File)
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
		fmt.Fprintf(&out, "  %d scope=%d kind=%d name=%s owner_type=%d exported=%t\n", declaration.ID, declaration.Scope, declaration.Kind, declaration.Name, declaration.OwnerType, declaration.Exported)
	}

	out.WriteString("bindings:\n")
	for _, binding := range p.Bindings {
		fmt.Fprintf(&out, "  %d scope=%d previous=%d name=%s type=%d value=%d mutable=%t\n", binding.ID, binding.Scope, binding.Previous, binding.Name, binding.TypeSlot, binding.Value, binding.Mutable)
	}

	out.WriteString("functions:\n")
	for _, function := range p.Functions {
		fmt.Fprintf(&out, "  %d declaration=%d parameters=%v result=%d effect=%d generic=%t any=%t cast=%t\n", function.ID, function.Declaration, function.Parameters, function.Result, function.Effect, function.Generic, function.Any, function.Cast)
	}

	out.WriteString("aliases:\n")
	for _, alias := range p.Aliases {
		fmt.Fprintf(&out, "  declaration=%d target=%d\n", alias.Declaration, alias.Target)
	}

	out.WriteString("expressions:\n")
	for _, expression := range p.Expressions {
		fmt.Fprintf(&out, "  %d scope=%d declaration=%d function=%d value=%d\n", expression.ID, expression.Scope, expression.Declaration, expression.Function, expression.Value)
	}

	out.WriteString("calls:\n")
	for _, call := range p.Calls {
		fmt.Fprintf(&out, "  %d scope=%d callee=%d qualifier=%s name=%s curry=%t arguments=[", call.ID, call.Scope, call.Callee, call.Qualifier, call.Name, call.Curry)
		for i, argument := range call.Arguments {
			if i > 0 {
				out.WriteString(",")
			}
			fmt.Fprintf(&out, "%s:%d:literal=%d:template=%t", argument.Name, argument.Type, argument.Literal, argument.LiteralTemplate)
		}
		out.WriteString("]\n")
	}

	out.WriteString("import_aliases:\n")
	for _, alias := range p.ImportAliases {
		fmt.Fprintf(&out, "  scope=%d name=%s module=%d declaration=%d source=%d\n", alias.Scope, alias.LocalName, alias.Module, alias.Declaration, alias.Source)
	}

	out.WriteString("literal_defaults:\n")
	fmt.Fprintf(&out, "  int=%d float=%d\n", p.LiteralDefaults.Int, p.LiteralDefaults.Float)

	out.WriteString("module_resolution:\n")
	fmt.Fprintf(&out, "  program_root=%s search_paths=%v\n", p.ModuleResolution.ProgramRoot, p.ModuleResolution.SearchPaths)

	out.WriteString("prolog:\n")
	fmt.Fprintf(&out, "  enabled=%t modules=%v architecture=%s\n", p.Prolog.Enabled, p.Prolog.Modules, p.Prolog.Architecture)

	out.WriteString("diagnostics:\n")
	for _, diagnostic := range p.Diagnostics {
		fmt.Fprintf(&out, "  severity=%d code=%s message=%s canonical_module=%d canonical_declaration=%d import_path=%s import_candidates=%v import_search_roots=%v import_resolved_path=%s\n", diagnostic.Severity, diagnostic.Code, diagnostic.Message, diagnostic.CanonicalModule, diagnostic.CanonicalDeclaration, diagnostic.ImportPath, diagnostic.ImportCandidates, diagnostic.ImportSearchRoots, diagnostic.ImportResolvedPath)
	}

	return out.String()
}
