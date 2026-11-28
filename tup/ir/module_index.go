package ir

// BuildModuleIndex performs the shared loader-independent source indexing path.
func BuildModuleIndex(filenames []string, osName, architecture string) (SourceFileTable, ModuleTable, []string) {
	files, invalid := IndexSourceFiles(filenames, osName, architecture)
	return files, BuildModulesFromSourceFiles(files), invalid
}
