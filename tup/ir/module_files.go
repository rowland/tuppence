package ir

import (
	pathpkg "path"
	"strings"
)

// ModuleFileName is the loader-neutral interpretation of a .tup source name.
type ModuleFileName struct {
	ModuleName  string
	Constraints []string
}

// ParseModuleFileName parses the module prefix and optional --OS/ARCH suffix.
func ParseModuleFileName(filename string) (ModuleFileName, bool) {
	base := pathpkg.Base(filename)
	if !strings.HasSuffix(base, ".tup") {
		return ModuleFileName{}, false
	}
	stem := strings.TrimSuffix(base, ".tup")
	constraintText := ""
	if marker := strings.Index(stem, "--"); marker >= 0 {
		constraintText = stem[marker+2:]
		stem = stem[:marker]
	}
	moduleName := stem
	for i, r := range stem {
		if !isIdentifierPart(r) {
			moduleName = stem[:i]
			break
		}
	}
	if moduleName == "" || !isIdentifierStart(rune(moduleName[0])) {
		return ModuleFileName{}, false
	}
	result := ModuleFileName{ModuleName: moduleName}
	if constraintText != "" {
		for _, constraint := range strings.Split(constraintText, "-") {
			if constraint == "" || !allIdentifierParts(constraint) {
				return ModuleFileName{}, false
			}
			result.Constraints = append(result.Constraints, constraint)
		}
	}
	return result, true
}

// Matches reports whether the file is eligible for the requested target.
func (f ModuleFileName) Matches(osName, architecture string) bool {
	for _, constraint := range f.Constraints {
		if constraint != osName && constraint != architecture {
			return false
		}
	}
	return true
}

func isIdentifierStart(r rune) bool {
	return r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || r >= '0' && r <= '9'
}

func allIdentifierParts(value string) bool {
	for _, r := range value {
		if !isIdentifierPart(r) {
			return false
		}
	}
	return value != ""
}
