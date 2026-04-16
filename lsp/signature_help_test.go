package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestTextDocumentSignatureHelp(t *testing.T) {
	input := `let add: fn(int, int) -> int = fn(x, y) { return x + y; };
let z: int = add(1, );`
	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.SignatureHelpParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 1, Character: 19},
		},
	}

	res, err := textDocumentSignatureHelp(nil, params)
	assertNoError(t, err)

	if res == nil {
		t.Fatalf("expected signature help, got nil")
	}

	if len(res.Signatures) == 0 {
		t.Fatalf("expected at least one signature")
	}

	sig := res.Signatures[0]
	if len(sig.Parameters) != 2 {
		t.Errorf("expected 2 parameters in signature, got %d", len(sig.Parameters))
	}

	if res.ActiveParameter != nil && *res.ActiveParameter != 1 {
		t.Errorf("expected active parameter 1, got %v", *res.ActiveParameter)
	}
}
