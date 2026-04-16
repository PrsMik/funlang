package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestIsValidIdentifier(t *testing.T) {
	if !isValidIdentifier("validName") {
		t.Errorf("expected validName to be valid")
	}
	if isValidIdentifier("1invalid") {
		t.Errorf("expected 1invalid to be invalid")
	}
	if isValidIdentifier("let") {
		t.Errorf("expected keyword 'let' to be invalid")
	}
}

func TestTextDocumentRename(t *testing.T) {
	input := `let abc: int = 1; let xyz: int = abc;`
	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.RenameParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 0, Character: 34},
		},
		NewName: "renamedVar",
	}

	res, err := textDocumentRename(nil, params)
	assertNoError(t, err)

	if res == nil {
		t.Fatalf("expected WorkspaceEdit, got nil")
	}

	changes := res.Changes[testURI]
	if len(changes) != 2 {
		t.Errorf("expected 2 text edits (1 decl + 1 usage), got %d", len(changes))
	}

	for _, edit := range changes {
		if edit.NewText != "renamedVar" {
			t.Errorf("expected new text 'renamedVar', got %q", edit.NewText)
		}
	}
}

func TestTextDocumentFormatting(t *testing.T) {
	input := "let    x:int   =1;"
	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.DocumentFormattingParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
		Options: protocol.FormattingOptions{
			protocol.FormattingOptionTabSize:      4,
			protocol.FormattingOptionInsertSpaces: true},
	}

	res, err := textDocumentFormatting(nil, params)
	assertNoError(t, err)

	if len(res) == 0 {
		t.Fatalf("expected text edits, got empty")
	}

	edit := res[0]
	expectedFormatted := "let x: int = 1;\n"
	if edit.NewText != expectedFormatted {
		t.Errorf("expected %q, got %q", expectedFormatted, edit.NewText)
	}
}
