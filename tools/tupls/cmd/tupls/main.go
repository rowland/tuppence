package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rowland/tuppence/tup/tok"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const languageName = "tuppence"

var version = "0.0.1"

var handler protocol.Handler
var logger *log.Logger

func main() {
	// Configure logging to a file
	logFile, err := os.OpenFile("/tmp/tupls.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger = log.New(logFile, "", log.LstdFlags|log.Lmicroseconds)

	handler = protocol.Handler{
		Initialize:            initialize,
		Initialized:           initialized,
		Shutdown:              shutdown,
		SetTrace:              setTrace,
		TextDocumentDidOpen:   textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,
		TextDocumentDidClose:  textDocumentDidClose,
	}

	server := server.NewServer(&handler, languageName, false)
	server.RunStdio()
}

func initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	// Configure initial capabilities
	openClose := true
	change := protocol.TextDocumentSyncKindFull
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: &openClose,
		Change:    &change,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    languageName,
			Version: &version,
		},
	}, nil
}

func initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(ctx *glsp.Context) error {
	return nil
}

func setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	logger.Printf("Document opened: %s", params.TextDocument.URI)
	// Validate the document on open
	return validateDocument(ctx, params.TextDocument.URI, params.TextDocument.Text)
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	logger.Printf("Document changed: %s", params.TextDocument.URI)
	logger.Printf("Number of content changes: %d", len(params.ContentChanges))
	// For full sync, we get the full content
	if len(params.ContentChanges) > 0 {
		// Handle both TextDocumentContentChangeEvent and TextDocumentContentChangeEventWhole
		switch change := params.ContentChanges[0].(type) {
		case protocol.TextDocumentContentChangeEvent:
			logger.Printf("Content length: %d", len(change.Text))
			// Clear old diagnostics first
			ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
				URI:         params.TextDocument.URI,
				Diagnostics: []protocol.Diagnostic{},
			})
			// Then validate the new content
			return validateDocument(ctx, params.TextDocument.URI, change.Text)
		case protocol.TextDocumentContentChangeEventWhole:
			logger.Printf("Content length: %d", len(change.Text))
			// Clear old diagnostics first
			ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
				URI:         params.TextDocument.URI,
				Diagnostics: []protocol.Diagnostic{},
			})
			// Then validate the new content
			return validateDocument(ctx, params.TextDocument.URI, change.Text)
		default:
			logger.Printf("Unexpected content change type: %T", params.ContentChanges[0])
		}
	}
	return nil
}

func textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	logger.Printf("Document closed: %s", params.TextDocument.URI)
	// Clear any diagnostics when document is closed
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: []protocol.Diagnostic{},
	})
	return nil
}

func validateDocument(ctx *glsp.Context, uri protocol.DocumentUri, content string) error {
	// Extract filename from URI
	u, err := url.Parse(string(uri))
	if err != nil {
		logger.Printf("Error parsing URI: %v", err)
		return err
	}
	filename := filepath.Base(u.Path)
	logger.Printf("Validating document: %s", filename)
	logger.Printf("Content length: %d", len(content))

	// Tokenize the content
	tokens, err := tok.Tokenize([]byte(content), filename)
	if err != nil {
		logger.Printf("Error tokenizing: %v", err)
		return err
	}

	logger.Printf("Found %d tokens", len(tokens))

	// Convert invalid tokens to diagnostics
	var diagnostics []protocol.Diagnostic
	for i, token := range tokens {
		if token.Invalid {
			logger.Printf("Found invalid token at index %d: %s", i, token.Value())
			logger.Printf("Token type: %v", token.Type)
			logger.Printf("Token offset: %d, length: %d, error offset: %d", token.Offset, token.Length, token.ErrorOffset)
			// Convert token position to line and character
			startLine, startChar := token.File.Position(int(token.Offset))
			endLine, endChar := token.File.Position(int(token.Offset + token.Length))
			logger.Printf("Token position: line %d, char %d to line %d, char %d", startLine, startChar, endLine, endChar)

			// Create diagnostic for invalid token
			severity := protocol.DiagnosticSeverityError
			source := languageName
			message := "Invalid token"
			if token.ErrorOffset > 0 {
				// If we have a specific error offset, use it to provide more precise error location
				_, errorChar := token.File.Position(int(token.ErrorOffset))
				message = fmt.Sprintf("Invalid token at position %d", errorChar)
			}

			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(startLine),
						Character: protocol.UInteger(startChar),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(endLine),
						Character: protocol.UInteger(endChar),
					},
				},
				Severity: &severity,
				Source:   &source,
				Message:  message,
			})
		}
	}

	logger.Printf("Publishing %d diagnostics", len(diagnostics))
	// Publish diagnostics
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
	return nil
}
