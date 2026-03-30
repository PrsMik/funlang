package lsp

import (
	"funlang/type_checker"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

var (
	handler        protocol.Handler
	version        string = "0.1.0"
	documents             = make(map[string]string)
	documentStates        = make(map[string]*type_checker.TypeChecker)
)

func StartServer() {
	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		TextDocumentDidOpen:    textDocumentDidOpen,
		TextDocumentDidChange:  textDocumentDidChange,
		TextDocumentHover:      textDocumentHover,
		TextDocumentDefinition: textDocumentDefinition,
		TextDocumentCompletion: textDocumentCompletion,
	}

	srv := server.NewServer(&handler, "funlang-lsp", false)

	srv.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	syncKind := protocol.TextDocumentSyncKindFull
	capabilities.TextDocumentSync = syncKind

	hoverProvider := true
	capabilities.HoverProvider = &hoverProvider

	definitionProvider := true
	capabilities.DefinitionProvider = &definitionProvider

	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{
			":", // типы
			"=", // значения
			"(", // аргументы
			",", // следующие аргументы
			" ", // ключевые слова
			".", // методы (?)
		}}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "funlang-lsp",
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	documents[params.TextDocument.URI] = params.TextDocument.Text
	validateDocument(context, params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) > 0 {
		var newText string

		switch change := params.ContentChanges[0].(type) {
		case protocol.TextDocumentContentChangeEvent:
			newText = change.Text
		case protocol.TextDocumentContentChangeEventWhole:
			newText = change.Text
		}

		if newText != "" {
			documents[params.TextDocument.URI] = newText
			validateDocument(context, params.TextDocument.URI, newText)
		}
	}
	return nil
}
