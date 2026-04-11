package lsp

import (
	"fmt"
	"funlang/type_checker"
	"runtime/debug"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

// TODO
// 4. Семантическая подсветка
// 5. Переименование
// 6. Автоформатирование (и комментарии?)

var (
	handler        protocol.Handler
	version        string = "0.1.0"
	documents             = make(map[string]string)
	documentStates        = make(map[string]*type_checker.TypeChecker)
)

func StartServer() {
	handler = protocol.Handler{
		Initialize:                 initialize,
		Initialized:                initialized,
		TextDocumentDidOpen:        textDocumentDidOpen,
		TextDocumentDidChange:      textDocumentDidChange,
		TextDocumentHover:          textDocumentHover,
		TextDocumentDefinition:     textDocumentDefinition,
		TextDocumentCompletion:     textDocumentCompletion,
		TextDocumentSignatureHelp:  textDocumentSignatureHelp,
		TextDocumentDocumentSymbol: textDocumentDocumentSymbol,
		TextDocumentReferences:     textDocumentReferences,
	}

	srv := server.NewServer(&handler, "funlang-lsp", false)

	// srv.RunStdio()
	srv.RunTCP("127.0.0.1:5007")
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	syncKind := protocol.TextDocumentSyncKindFull
	capabilities.TextDocumentSync = syncKind

	hoverProvider := true
	capabilities.HoverProvider = &hoverProvider

	referencesProvider := true
	capabilities.ReferencesProvider = &referencesProvider

	definitionProvider := true
	capabilities.DefinitionProvider = &definitionProvider

	documentSymbolProvider := true
	capabilities.DocumentSymbolProvider = &documentSymbolProvider

	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{
			":", // типы
			"=", // значения
			"(", // аргументы
			",", // следующие аргументы
			" ", // ключевые слова
			".", // методы (?)
		}}

	capabilities.SignatureHelpProvider = &protocol.SignatureHelpOptions{
		TriggerCharacters:   []string{"(", ",", " "},
		RetriggerCharacters: []string{","},
	}

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

func handlePanic(context *glsp.Context) {
	if r := recover(); r != nil {
		stack := string(debug.Stack())

		errorMessage := fmt.Sprintf("LSP Panic recovered: %v\nStack trace:\n%s", r, stack)

		fmt.Println(errorMessage)

		context.Notify(protocol.ServerWindowLogMessage, protocol.LogMessageParams{
			Type:    protocol.MessageTypeError,
			Message: errorMessage,
		})
	}
}
