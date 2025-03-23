package main

import (
	"net/url"
	"path/filepath"

	"github.com/rowland/tuppence/tup/tok"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const languageName = "tuppence"

var version = "0.0.1"

var handler protocol.Handler

func main() {
	// Configure logging
	commonlog.Configure(1, nil)

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
	// For now, just validate the document on open
	return validateDocument(ctx, params.TextDocument.URI, params.TextDocument.Text)
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	// For full sync, we get the full content
	if len(params.ContentChanges) > 0 {
		if textChange, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent); ok {
			return validateDocument(ctx, params.TextDocument.URI, textChange.Text)
		}
	}
	return nil
}

func textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
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
		return err
	}
	filename := filepath.Base(u.Path)

	// Tokenize the content
	tokens, err := tok.Tokenize([]byte(content), filename)
	if err != nil {
		return err
	}

	// Convert invalid tokens to diagnostics
	var diagnostics []protocol.Diagnostic
	for _, token := range tokens {
		if token.Invalid {
			// Convert token position to line and character
			startLine, startChar := token.File.Position(int(token.Offset))
			endLine, endChar := token.File.Position(int(token.Offset + token.Length))

			// Create diagnostic for invalid token
			severity := protocol.DiagnosticSeverityError
			source := languageName
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(startLine), // Source.Position already returns 0-based line numbers
						Character: protocol.UInteger(startChar), // Source.Position already returns 0-based column numbers
					},
					End: protocol.Position{
						Line:      protocol.UInteger(endLine),
						Character: protocol.UInteger(endChar),
					},
				},
				Severity: &severity,
				Source:   &source,
				Message:  "Invalid token",
			})
		}
	}

	// Publish diagnostics
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
	return nil
}
