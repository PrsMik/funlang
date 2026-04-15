package lsp

import (
	"fmt"
	"funlang/types"
	"runtime/debug"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

var (
	handler        protocol.Handler
	version        string = "0.1.0"
	documents             = make(map[string]string)
	documentStates        = make(map[string]*types.Info)
)

var semanticTokenTypes = []string{
	string(protocol.SemanticTokenTypeVariable),
	string(protocol.SemanticTokenTypeFunction),
	string(protocol.SemanticTokenTypeParameter),
	string(protocol.SemanticTokenTypeType),
}

var semanticTokenModifiers = []string{
	string(protocol.SemanticTokenModifierDeclaration),
	string(protocol.SemanticTokenModifierReadonly),
}

func StartServer(runTCP bool) {
	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		TextDocumentDidOpen:            textDocumentDidOpen,
		TextDocumentDidChange:          textDocumentDidChange,
		TextDocumentHover:              textDocumentHover,
		TextDocumentDefinition:         textDocumentDefinition,
		TextDocumentCompletion:         textDocumentCompletion,
		TextDocumentSignatureHelp:      textDocumentSignatureHelp,
		TextDocumentDocumentSymbol:     textDocumentDocumentSymbol,
		TextDocumentReferences:         textDocumentReferences,
		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
		TextDocumentRename:             textDocumentRename,
		TextDocumentPrepareRename:      textDocumentPrepareRename,
		TextDocumentFormatting:         textDocumentFormatting,
	}

	srv := server.NewServer(&handler, "funlang-lsp", false)

	if runTCP {
		srv.RunTCP("127.0.0.1:5007")
	} else {
		srv.RunStdio()
	}
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
			// " ", // ключевые слова
			".", // методы (?)
		}}

	capabilities.SignatureHelpProvider = &protocol.SignatureHelpOptions{
		TriggerCharacters:   []string{"(", ",", " "},
		RetriggerCharacters: []string{","},
	}

	fullSupport := true
	capabilities.SemanticTokensProvider = protocol.SemanticTokensRegistrationOptions{
		SemanticTokensOptions: protocol.SemanticTokensOptions{
			Legend: protocol.SemanticTokensLegend{
				TokenTypes:     semanticTokenTypes,
				TokenModifiers: semanticTokenModifiers,
			},
			Full: &fullSupport,
		},
	}

	prepareSupport := true
	capabilities.RenameProvider = protocol.RenameOptions{
		PrepareProvider: &prepareSupport,
	}

	documentFormattingProvider := true
	capabilities.DocumentFormattingProvider = &documentFormattingProvider

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
