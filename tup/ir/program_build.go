package ir

// BuildInitialScopes creates one explicit initial scope for each module.
func BuildInitialScopes(modules ModuleTable) ScopeTable {
	scopes := make(ScopeTable, 0, len(modules))
	for _, module := range modules {
		scopes = append(scopes, Scope{ID: ScopeID(module.ID), Module: module.ID, Kind: ScopeModule, Initial: true, Origin: module.Origin})
	}
	return scopes
}

// BuildIndexedProgram constructs the initial semantic snapshot from one target
// configuration, keeping source selection and import resolution consistent.
func BuildIndexedProgram(filenames []string, target CompilationTarget, prolog PrologConfig) (IndexedProgram, []string) {
	if target.PrimaryPath == "" {
		return IndexedProgram{}, []string{"primary module path is required"}
	}
	files, modules, invalid := BuildModuleIndexForTarget(filenames, target)
	program := IndexedProgram{
		Modules:          modules,
		SourceFiles:      files,
		Scopes:           BuildInitialScopes(modules),
		ModuleResolution: target.ModuleRoots,
		Prolog:           prolog,
	}
	primaryFound := false
	for _, module := range modules {
		if module.ResolutionPath() == target.PrimaryPath {
			program.PrimaryModule = module.ID
			primaryFound = true
			break
		}
	}
	if !primaryFound {
		invalid = append(invalid, target.PrimaryPath)
	}
	if prolog.Enabled && len(target.PrologPaths) > 0 {
		program.Prolog.Modules = nil
		program.Prolog.Architecture = target.Architecture
		roots := append([]string(nil), target.ModuleRoots.SearchPaths...)
		if target.ModuleRoots.ProgramRoot != "" {
			roots = append([]string{target.ModuleRoots.ProgramRoot}, roots...)
		}
		for _, path := range target.PrologPaths {
			found := false
			for _, module := range modules {
				if module.ResolutionPath() == path {
					program.Prolog.Modules = append(program.Prolog.Modules, module.ID)
					found = true
					break
				}
			}
			if !found {
				invalid = append(invalid, path)
				if len(files) > 0 {
					program.Diagnostics = append(program.Diagnostics, NewUnresolvedPrologDiagnostic(Origin{Primary: Span{File: files[0].ID}}, path, roots))
				}
			}
		}
	}
	return program, invalid
}
