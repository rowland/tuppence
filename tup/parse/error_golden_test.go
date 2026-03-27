package parse

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

const updateErrorGoldensEnv = "UPDATE_ERROR_GOLDENS"

func TestErrorGoldenFixtures(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	baseDir := filepath.Join(filepath.Dir(thisFile), "testdata", "error")
	inputDir := filepath.Join(baseDir, "input")
	outputDir := filepath.Join(baseDir, "output")
	update := os.Getenv(updateErrorGoldensEnv) != ""

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
					got := parseErrorFixtureEntry(t, inputPath, entry)
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
						t.Fatalf("error message mismatch\n\nwant:\n%s\n\ngot:\n%s", expected.input, got)
					}
				})
			}

			if update {
				if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", filepath.Dir(outputPath), err)
				}
				if err := os.WriteFile(outputPath, renderErrorFixtureFile(gotEntries), 0o644); err != nil {
					t.Fatalf("WriteFile(%q): %v", outputPath, err)
				}
			}
		})
	}
}

// parseErrorFixtureEntry tokenizes and parses a broken code sample, expecting
// a parse error. It fails the test if parsing unexpectedly succeeds.
func parseErrorFixtureEntry(t *testing.T, fixturePath string, entry topLevelFixtureEntry) string {
	t.Helper()

	src := source.NewSource([]byte(entry.input), fixturePath)
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		// Tokenizer errors are acceptable — treat them as the reported error.
		return err.Error()
	}

	_, _, err = TopLevelItem(tokens)
	if err == nil {
		t.Fatalf("expected parse error for %q, but parsing succeeded", entry.name)
	}

	return err.Error()
}

func renderErrorFixtureFile(entries []topLevelFixtureEntry) []byte {
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
