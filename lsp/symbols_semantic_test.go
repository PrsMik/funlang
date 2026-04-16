package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestTextDocumentDocumentSymbol(t *testing.T) {
	input := `let a: int = 1;
let myFunc: fn() -> int = fn() { return 1; };`
	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
	}

	res, err := textDocumentDocumentSymbol(nil, params)
	assertNoError(t, err)

	symbols, ok := res.([]protocol.DocumentSymbol)
	if !ok || len(symbols) != 2 {
		t.Fatalf("expected 2 symbols, got %v", res)
	}

	if symbols[0].Name != "a" || symbols[0].Kind != protocol.SymbolKindVariable {
		t.Errorf("expected first symbol to be variable 'a'")
	}
	if symbols[1].Name != "myFunc" || symbols[1].Kind != protocol.SymbolKindFunction {
		t.Errorf("expected second symbol to be function 'myFunc'")
	}
}

func TestTextDocumentSemanticTokensFull(t *testing.T) {
	input := `let myVar: int = 10;`
	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
	}

	res, err := textDocumentSemanticTokensFull(nil, params)
	assertNoError(t, err)

	if res == nil || len(res.Data) == 0 {
		t.Fatalf("expected semantic tokens data, got empty")
	}
}
