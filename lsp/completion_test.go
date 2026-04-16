package lsp

import (
	"funlang/token"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestIsExpectedTypeContext(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []token.Token
		expected bool
	}{
		{
			name: "Type after colon in let",
			tokens: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COLON, Literal: ":"},
			},
			expected: true,
		},
		{
			name: "Assignment value context",
			tokens: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COLON, Literal: ":"},
				{Type: token.INT_TYPE, Literal: "int"},
				{Type: token.ASSIGN, Literal: "="},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExpectedTypeContext(tt.tokens)
			if got != tt.expected {
				t.Errorf("isExpectedTypeContext() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTextDocumentCompletion(t *testing.T) {
	input := `let alpha: int = 1;
let beta: `

	setupTestDocument(input)
	defer clearTestState()

	params := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 1, Character: 10},
		},
	}

	res, err := textDocumentCompletion(nil, params)
	assertNoError(t, err)

	items, ok := res.([]protocol.CompletionItem)
	if !ok {
		t.Fatalf("expected completion items, got %T", res)
	}

	if len(items) == 0 {
		t.Errorf("expected completions, got empty list")
	}

	foundIntType := false
	for _, item := range items {
		if item.Label == "int" {
			foundIntType = true
			break
		}
	}

	if !foundIntType {
		t.Errorf("expected 'int' type in completion items")
	}
}
