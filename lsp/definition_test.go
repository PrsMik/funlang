package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestTextDocumentDefinition(t *testing.T) {
	input := `let target: int = 10;
let another: int = target + 5;`

	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.DefinitionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 1, Character: 21},
		},
	}

	res, err := textDocumentDefinition(nil, params)
	assertNoError(t, err)

	locations, ok := res.([]protocol.Location)
	if !ok || len(locations) == 0 {
		t.Fatalf("expected definition locations, got %v", res)
	}

	loc := locations[0]
	if loc.URI != testURI {
		t.Errorf("expected URI %s, got %s", testURI, loc.URI)
	}

	if loc.Range.Start.Line != 0 || loc.Range.Start.Character != 4 {
		t.Errorf("wrong definition range start, got: %+v", loc.Range.Start)
	}
}
