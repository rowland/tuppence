package ir

// SourceFile records the loader's stable identity and parsed module-file metadata.
type SourceFile struct {
	ID     FileID
	Path   string
	Module ModuleFileName
}

// SourceFileTable is an ordered collection of indexed source files.
type SourceFileTable []SourceFile

// Valid reports whether file identities and paths are unique and metadata is valid.
func (t SourceFileTable) Valid() bool {
	seenIDs := make(map[FileID]struct{}, len(t))
	seenPaths := make(map[string]struct{}, len(t))
	for _, file := range t {
		if !file.ID.Valid() || file.Path == "" || file.Module.ModuleName == "" {
			return false
		}
		if _, exists := seenIDs[file.ID]; exists {
			return false
		}
		if _, exists := seenPaths[file.Path]; exists {
			return false
		}
		seenIDs[file.ID] = struct{}{}
		seenPaths[file.Path] = struct{}{}
	}
	return true
}

// IndexSourceFiles selects and indexes source files for a target.
func IndexSourceFiles(filenames []string, osName, architecture string) (SourceFileTable, []string) {
	selected, invalid := SelectModuleFiles(filenames, osName, architecture)
	var table SourceFileTable
	var next FileID = 1
	for _, group := range selected {
		for _, filename := range group.Files {
			parsed, _ := ParseModuleFileName(filename)
			table = append(table, SourceFile{ID: next, Path: filename, Module: parsed})
			next++
		}
	}
	return table, invalid
}
