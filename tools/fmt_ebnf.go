// fmt_ebnf.go
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"html/template"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Rule represents a single grammar rule.
type Rule struct {
	Name    string
	Content string
}

//go:embed fmt_ebnf.html
var tmplHTML string

//go:embed fmt_ebnf.js
var scriptJS string

func main() {
	inputPath := flag.String("i", "", "Input EBNF file")
	outputPath := flag.String("o", "", "Output HTML file")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: go run fmt_ebnf.go -i input.ebnf -o output.html")
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}
	content := string(data)

	// Split input into rule blocks (blank lines separate rules).
	ruleBlocks := splitRules(content)

	// Extract each rule (assumes rule name appears before the '=' in the first line).
	var rules []Rule
	ruleNameRegex := regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
	for _, block := range ruleBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		m := ruleNameRegex.FindStringSubmatch(block)
		if m == nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not extract rule name from block:\n%s\n", block)
			continue
		}
		name := m[1]
		rules = append(rules, Rule{Name: name, Content: block})
	}

	// Sort rules alphabetically.
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Name < rules[j].Name
	})

	// Build a set of all rule names.
	ruleSet := make(map[string]bool)
	for _, r := range rules {
		ruleSet[r.Name] = true
	}

	// Compute reverse dependencies: for each rule, list all rules that reference it.
	dependents := make(map[string]map[string]bool)
	for _, r := range rules {
		dependents[r.Name] = make(map[string]bool)
	}
	for _, r := range rules {
		refs := extractReferences(r.Content, ruleSet)
		// Exclude self-references.
		delete(refs, r.Name)
		for ref := range refs {
			dependents[ref][r.Name] = true
		}
	}

	// Convert dependency maps to sorted slices.
	depSlices := make(map[string][]string)
	for rule, depSet := range dependents {
		var deps []string
		for d := range depSet {
			deps = append(deps, d)
		}
		sort.Strings(deps)
		depSlices[rule] = deps
	}

	depJSON, err := json.Marshal(depSlices)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling dependency mapping: %v\n", err)
		os.Exit(1)
	}

	// Build HTML for the rule blocks.
	var ruleBlocksHTML strings.Builder
	for _, r := range rules {
		ruleBlocksHTML.WriteString(fmt.Sprintf(`<div class="rule" id="%s">`, r.Name))
		ruleBlocksHTML.WriteString("\n  <code>")
		escapedContent := html.EscapeString(r.Content)
		ruleBlocksHTML.WriteString(escapedContent)
		ruleBlocksHTML.WriteString("</code>\n</div>\n\n")
	}

	// Build a map from rule name to its (escaped) content for the preview tooltip.
	ruleContentsMap := make(map[string]string)
	for _, r := range rules {
		ruleContentsMap[r.Name] = html.EscapeString(r.Content)
	}
	contentsJSON, err := json.Marshal(ruleContentsMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling rule contents: %v\n", err)
		os.Exit(1)
	}

	// Prepare the data for the template.
	dataTmpl := struct {
		RuleBlocks     template.HTML
		DependentsJSON template.JS
		ContentsJSON   template.JS
		ScriptJS       template.JS
	}{
		RuleBlocks:     template.HTML(ruleBlocksHTML.String()),
		DependentsJSON: template.JS(depJSON),
		ContentsJSON:   template.JS(contentsJSON),
		ScriptJS:       template.JS(scriptJS),
	}

	tmpl, err := template.New("page").Parse(tmplHTML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, dataTmpl); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
}

// splitRules splits the input text into rule blocks using one or more blank lines.
func splitRules(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	re := regexp.MustCompile(`\n\s*\n`)
	return re.Split(text, -1)
}

// extractReferences scans ruleContent (ignoring text inside quotes) and returns a set of tokens found in ruleSet.
func extractReferences(ruleContent string, ruleSet map[string]bool) map[string]bool {
	refs := make(map[string]bool)
	inQuote := false
	var token strings.Builder
	for i, r := range ruleContent {
		if r == '"' {
			inQuote = !inQuote
			if token.Len() > 0 && !inQuote {
				word := token.String()
				if ruleSet[word] {
					refs[word] = true
				}
				token.Reset()
			}
			continue
		}
		if inQuote {
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			token.WriteRune(r)
		} else {
			if token.Len() > 0 {
				word := token.String()
				if ruleSet[word] {
					refs[word] = true
				}
				token.Reset()
			}
		}
		if i == len(ruleContent)-1 && token.Len() > 0 {
			word := token.String()
			if ruleSet[word] {
				refs[word] = true
			}
		}
	}
	return refs
}
