# Makefile

# Variables for file names.
TUP = tup.ebnf
HTML = tup.ebnf.html
FMT = tools/fmt_ebnf.go
FMT_HTML = tools/fmt_ebnf.html
FMT_JS = tools/fmt_ebnf.js

# Default target.
all: $(HTML)

# Generate tup.ebnf.html when tup.ebnf (or the Go formatter) is newer.
$(HTML): $(TUP) $(FMT) $(FMT_HTML) $(FMT_JS)
	@echo "Updating $(HTML) from $(TUP)..."
	go run $(FMT) -i $(TUP) -o $(HTML)

# Optional target to update the grammar explicitly.
grammar: $(HTML)

# Clean up the generated file.
clean:
	rm -f $(HTML)

.PHONY: all grammar clean
