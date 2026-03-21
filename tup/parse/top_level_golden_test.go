package parse

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

const updateTopLevelGoldensEnv = "UPDATE_TOP_LEVEL_GOLDENS"

type topLevelFixtureEntry struct {
	name  string
	input string
}

func TestTopLevelGoldenFixtures(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	baseDir := filepath.Join(filepath.Dir(thisFile), "testdata", "top_level")
	inputDir := filepath.Join(baseDir, "input")
	outputDir := filepath.Join(baseDir, "output")
	update := os.Getenv(updateTopLevelGoldensEnv) != ""

	files, err := os.ReadDir(inputDir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", inputDir, err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".tup" {
			continue
		}

		file := file
		t.Run(strings.TrimSuffix(file.Name(), ".tup"), func(t *testing.T) {
			inputPath := filepath.Join(inputDir, file.Name())
			outputPath := filepath.Join(outputDir, file.Name())

			inputEntries, err := readTopLevelFixtureFile(inputPath)
			if err != nil {
				t.Fatalf("readTopLevelFixtureFile(%q): %v", inputPath, err)
			}

			var expectedEntries []topLevelFixtureEntry
			if !update {
				expectedEntries, err = readTopLevelFixtureFile(outputPath)
				if err != nil {
					t.Fatalf("readTopLevelFixtureFile(%q): %v", outputPath, err)
				}
				if len(expectedEntries) != len(inputEntries) {
					t.Fatalf("entry count mismatch: input=%d output=%d", len(inputEntries), len(expectedEntries))
				}
			}

			var gotEntries []topLevelFixtureEntry
			for i, entry := range inputEntries {
				t.Run(entry.name, func(t *testing.T) {
					item := parseTopLevelFixtureEntry(t, inputPath, entry)
					got := strings.TrimSuffix(item.String(), "\n")
					gotEntries = append(gotEntries, topLevelFixtureEntry{
						name:  entry.name,
						input: got,
					})

					if update {
						return
					}

					expected := expectedEntries[i]
					if entry.name != expected.name {
						t.Fatalf("entry name mismatch: input=%q output=%q", entry.name, expected.name)
					}
					if got != expected.input {
						t.Fatalf("serialized output mismatch\n\nwant (%d):\n%q\n\ngot (%d):\n%q", len(expected.input), expected.input, len(got), got)
					}
				})
			}

			if update {
				if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", filepath.Dir(outputPath), err)
				}
				if err := os.WriteFile(outputPath, renderTopLevelFixtureFile(gotEntries), 0o644); err != nil {
					t.Fatalf("WriteFile(%q): %v", outputPath, err)
				}
			}
		})
	}
}

func parseTopLevelFixtureEntry(t *testing.T, fixturePath string, entry topLevelFixtureEntry) ast.TopLevelItem {
	t.Helper()

	src := source.NewSource([]byte(entry.input), fixturePath)
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		t.Fatalf("Tokenize(%q): %v", entry.name, err)
	}

	item, remainder, err := TopLevelItem(tokens)
	if err != nil {
		t.Fatalf("TopLevelItem(%q): %v", entry.name, err)
	}

	remainder = skipTrivia(remainder)
	if len(remainder) == 1 && remainder[0].Type == tok.TokEOF {
		return item
	}
	if len(remainder) > 0 && remainder[0].Type != tok.TokEOF {
		t.Fatalf("TopLevelItem(%q): unexpected trailing token %q", entry.name, remainder[0].Value())
	}

	return item
}

func readTopLevelFixtureFile(path string) ([]topLevelFixtureEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	var entries []topLevelFixtureEntry
	for i := 0; i < len(lines); {
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}

		if !strings.HasPrefix(lines[i], "# ") {
			return nil, fmt.Errorf("line %d: expected comment starting with '# '", i+1)
		}
		name := strings.TrimPrefix(lines[i], "# ")
		i++

		start := i
		for i < len(lines) && strings.TrimSpace(lines[i]) != "" {
			i++
		}

		if start == i {
			return nil, fmt.Errorf("entry %q: missing body", name)
		}

		body := strings.Join(lines[start:i], "\n")
		entries = append(entries, topLevelFixtureEntry{name: name, input: body})
	}

	return entries, nil
}

func renderTopLevelFixtureFile(entries []topLevelFixtureEntry) []byte {
	var buf bytes.Buffer
	for i, entry := range entries {
		if i > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString("# ")
		buf.WriteString(entry.name)
		buf.WriteString("\n")
		buf.WriteString(entry.input)
	}
	buf.WriteString("\n")
	return buf.Bytes()
}
