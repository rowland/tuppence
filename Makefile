# Makefile

# Variables for file names.
TUP = tup.ebnf
HTML = tup.ebnf.html
FMT = tools/fmt_ebnf.go
FMT_HTML = tools/fmt_ebnf.html
FMT_JS = tools/fmt_ebnf.js

# Language server variables
LS_DIR = tools/tupls
LS_BIN = $(LS_DIR)/tupls
INSTALL_DIR = $(HOME)/bin

# VS Code extension variables
VSCODE_DIR = $(LS_DIR)/vscode
VSCODE_EXT_DIR = $(HOME)/.vscode/extensions/rowland.tuppence-language-server-0.0.1

# Default target.
all: $(HTML)

# Generate tup.ebnf.html when tup.ebnf (or the Go formatter) is newer.
$(HTML): $(TUP) $(FMT) $(FMT_HTML) $(FMT_JS)
	@echo "Updating $(HTML) from $(TUP)..."
	go run $(FMT) -i $(TUP) -o $(HTML)

# Optional target to update the grammar explicitly.
grammar: $(HTML)

# Build the language server
ls-build:
	@echo "Building language server..."
	cd $(LS_DIR) && go build -o tupls ./cmd/tupls

# Install the language server
ls-install: ls-build
	@echo "Installing language server to $(INSTALL_DIR)/$(LS_BIN)..."
	mkdir -p $(INSTALL_DIR)
	cp $(LS_BIN) $(INSTALL_DIR)/

# Build the VS Code extension
vscode-build:
	@echo "Building VS Code extension..."
	cd $(VSCODE_DIR) && npm install && npm run compile

# Install the VS Code extension
vscode-install: vscode-build
	@echo "Installing VS Code extension..."
	rm -rf $(VSCODE_EXT_DIR)
	mkdir -p $(VSCODE_EXT_DIR)
	cp -r $(VSCODE_DIR)/out/* $(VSCODE_EXT_DIR)/
	cp -r $(VSCODE_DIR)/syntaxes $(VSCODE_EXT_DIR)/
	cp -r $(VSCODE_DIR)/language-configuration.json $(VSCODE_EXT_DIR)/
	cp -r $(VSCODE_DIR)/package.json $(VSCODE_EXT_DIR)/
	cp -r $(VSCODE_DIR)/icons $(VSCODE_EXT_DIR)/

# Build and install everything
dev-setup: ls-install vscode-install
	@echo "Development environment setup complete."
	@echo "Please restart VS Code to activate the extension."

# Clean up the generated files
clean: clean-ls clean-vscode
	rm -f $(HTML)

# Clean language server build artifacts
clean-ls:
	rm -f $(LS_BIN)

# Clean VS Code extension build artifacts
clean-vscode:
	rm -rf $(VSCODE_DIR)/out
	rm -rf $(VSCODE_DIR)/node_modules

test:
	cd tup && go test ./...

.PHONY: all grammar clean clean-ls clean-vscode ls-build ls-install vscode-build vscode-install dev-setup test
