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
	return program, invalid
}
