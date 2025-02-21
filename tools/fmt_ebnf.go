// fmt_ebnf.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Rule struct {
	Name    string
	Content string
}

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

	// Split input into rule blocks (assuming one or more blank lines separate rules)
	ruleBlocks := splitRules(content)

	// Extract rule name from each block (a valid identifier before the '=')
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

	// Sort rules alphabetically by name.
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Name < rules[j].Name
	})

	// Build a set of all rule names.
	ruleSet := make(map[string]bool)
	for _, r := range rules {
		ruleSet[r.Name] = true
	}

	// Compute dependencies: for each rule, find which other rules it references.
	// Then build a reverse mapping: for each rule, list all rules that reference it.
	dependents := make(map[string]map[string]bool)
	for _, r := range rules {
		dependents[r.Name] = make(map[string]bool)
	}
	for _, r := range rules {
		refs := extractReferences(r.Content, ruleSet)
		// Exclude self-reference.
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

	// Marshal the dependency mapping to JSON for JavaScript.
	depJSON, err := json.Marshal(depSlices)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling dependency mapping: %v\n", err)
		os.Exit(1)
	}

	// Build the HTML.
	var sb strings.Builder

	// Reduced margins/padding, plus absolute-position popup that scrolls if needed:
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>EBNF Grammar</title>
  <style>
    body { 
      margin: 0;
      padding: 0;
      background-color: #f4f4f4; 
      font-family: Arial, sans-serif; 
      line-height: 1.6; 
    }
    h1 { 
      text-align: center; 
      color: #2c3e50;
      margin: 0;
      padding: 10px 0;
    }
    .grammar { 
      max-width: 900px; 
      margin: 0 auto; 
      padding: 10px; 
    }
    .rule { 
      background: #fff;
      border: 1px solid #ddd; 
      border-radius: 4px; 
      margin-bottom: 10px; 
      padding: 10px;
      box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .rule code {
      display: block;
      font-family: "Courier New", monospace; 
      white-space: pre-wrap;
      word-break: break-word;
      margin: 0;
      padding: 0;
    }
    .rule a { 
      color: #3498db; 
      text-decoration: none; 
    }
    .rule a:hover { 
      text-decoration: underline; 
    }
    .rule-name { 
      color: #3498db; 
      cursor: pointer; 
    }
    .rule-name:hover { 
      text-decoration: underline; 
    }
    #dependencyPopup {
      display: none;
      position: absolute;  /* changed from fixed to absolute */
      background: #fff;
      border: 2px solid #3498db;
      padding: 10px;
      border-radius: 4px;
      box-shadow: 0 2px 6px rgba(0,0,0,0.3);
      z-index: 1000;
      max-width: 400px;
      max-height: 300px;
      overflow-y: auto; /* allow scrolling if the list is long */
      word-wrap: break-word;
    }
    #dependencyPopup ul {
      list-style-type: none;
      padding-left: 0;
      margin: 0;
    }
    #dependencyPopup li {
      margin: 5px 0;
    }
    #dependencyPopup a {
      color: #3498db;
    }
  </style>
</head>
<body>
  <h1>EBNF Grammar</h1>
  <div class="grammar">
`)

	// Output each rule in a <div>, with <code> for its content
	for _, rule := range rules {
		sb.WriteString(fmt.Sprintf(`    <div class="rule" id="%s">`, rule.Name))
		sb.WriteString("\n      <code>")
		escapedContent := html.EscapeString(rule.Content)
		sb.WriteString(escapedContent)
		sb.WriteString("</code>\n    </div>\n\n")
	}

	sb.WriteString(`  </div>
  <div id="dependencyPopup"></div>
`)

	// JavaScript code:
	// 1) Instead of inline onclick with ruleName, store data-rule in a span.
	// 2) On click, position the popup absolutely near the clicked span.
	// 3) Remove "Click here to close" text from the popup.
	jsBlock := "\n  <script>\n" +
		"    var dependents = " + string(depJSON) + ";\n\n" +
		"    function processTextNodes(node, ruleID, pattern) {\n" +
		"      if (node.nodeType === Node.TEXT_NODE) {\n" +
		"        var replaced = node.textContent.replace(pattern, function(match) {\n" +
		"          if (match === ruleID) {\n" +
		"            // self-reference -> clickable span with data-rule\n" +
		"            return `<span class=\"rule-name\" data-rule=\"${ruleID}\" onclick=\"showDeps(event)\">${match}</span>`;\n" +
		"          } else {\n" +
		"            return `<a href=\"#${match}\">${match}</a>`;\n" +
		"          }\n" +
		"        });\n" +
		"        var span = document.createElement(\"span\");\n" +
		"        span.innerHTML = replaced;\n" +
		"        node.parentNode.replaceChild(span, node);\n" +
		"      } else if (node.nodeType === Node.ELEMENT_NODE && !node.classList.contains(\"rule-name\")) {\n" +
		"        Array.from(node.childNodes).forEach(function(child) {\n" +
		"          processTextNodes(child, ruleID, pattern);\n" +
		"        });\n" +
		"      }\n" +
		"    }\n\n" +
		"    document.addEventListener(\"DOMContentLoaded\", function() {\n" +
		"      var ruleNames = Object.keys(dependents).sort();\n" +
		"      var pattern = new RegExp(\"\\\\b(\" + ruleNames.join(\"|\") + \")\\\\b\", \"g\");\n" +
		"      document.querySelectorAll(\"div.rule > code\").forEach(function(codeElem) {\n" +
		"        var ruleID = codeElem.parentElement.id;\n" +
		"        processTextNodes(codeElem, ruleID, pattern);\n" +
		"      });\n" +
		"    });\n\n" +
		"    function showDeps(e) {\n" +
		"      // e.target is the clicked span with data-rule\n" +
		"      var ruleName = e.target.dataset.rule;\n" +
		"      var deps = dependents[ruleName] || [];\n" +
		"      var message = \"\";\n" +
		"      if (deps.length === 0) {\n" +
		"        message = `<strong>${ruleName}</strong> is not referenced by any other rule.`;\n" +
		"      } else {\n" +
		"        message = `<strong>Rules that depend on ${ruleName}:</strong><br><ul>`;\n" +
		"        deps.forEach(function(dep) {\n" +
		"          message += `<li><a href=\"#${dep}\">${dep}</a></li>`;\n" +
		"        });\n" +
		"        message += \"</ul>\";\n" +
		"      }\n" +
		"      var popup = document.getElementById(\"dependencyPopup\");\n" +
		"      popup.innerHTML = message;\n" +
		"\n" +
		"      // Position the popup near the clicked span.\n" +
		"      var rect = e.target.getBoundingClientRect();\n" +
		"      popup.style.display = \"block\";\n" +
		"      popup.style.top = (window.scrollY + rect.bottom + 5) + \"px\";\n" +
		"      popup.style.left = (window.scrollX + rect.left) + \"px\";\n" +
		"      e.stopPropagation();\n" +
		"    }\n\n" +
		"    document.addEventListener(\"click\", function(e) {\n" +
		"      // If we click outside the popup or the rule-name span, hide the popup.\n" +
		"      var popup = document.getElementById(\"dependencyPopup\");\n" +
		"      if (popup.style.display === \"block\") {\n" +
		"        var isPopup = popup.contains(e.target);\n" +
		"        var isRuleName = e.target.classList && e.target.classList.contains(\"rule-name\");\n" +
		"        if (!isPopup && !isRuleName) {\n" +
		"          popup.style.display = \"none\";\n" +
		"        }\n" +
		"      }\n" +
		"    });\n" +
		"  </script>\n" +
		"</body>\n" +
		"</html>\n"

	sb.WriteString(jsBlock)

	err = os.WriteFile(*outputPath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}
}

// splitRules splits the input text into rule blocks.
// It assumes that one or more blank lines separate individual rules.
func splitRules(text string) []string {
	// Normalize line endings.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Split on one or more blank lines.
	re := regexp.MustCompile(`\n\s*\n`)
	return re.Split(text, -1)
}

// extractReferences scans ruleContent (ignoring text inside quoted strings)
// and returns a set of tokens that match names in ruleSet.
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
