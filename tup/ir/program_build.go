package ir

// BuildIndexedProgram constructs the initial semantic snapshot from one target
// configuration, keeping source selection and import resolution consistent.
func BuildIndexedProgram(filenames []string, target CompilationTarget, prolog PrologConfig) (IndexedProgram, []string) {
	files, modules, invalid := BuildModuleIndexForTarget(filenames, target)
	program := IndexedProgram{
		Modules:          modules,
		SourceFiles:      files,
		ModuleResolution: target.ModuleRoots,
		Prolog:           prolog,
	}
	if prolog.Enabled && len(target.PrologPaths) > 0 {
		program.Prolog.Modules = nil
		program.Prolog.Architecture = target.Architecture
		for _, path := range target.PrologPaths {
			for _, module := range modules {
				if module.ResolutionPath() == path {
					program.Prolog.Modules = append(program.Prolog.Modules, module.ID)
					break
				}
			}
		}
	}
	return program, invalid
}
