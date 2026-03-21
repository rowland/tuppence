package ast

import "strings"

const prettyIndent = "    "

func indentString(value string) string {
	if value == "" {
		return ""
	}

	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = prettyIndent + line
	}
	return strings.Join(lines, "\n")
}
