package lsp

import (
	"funlang/token"
	"strings"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestTextDocumentHover(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		position      protocol.Position
		expectedFound bool
		expectedText  string
	}{
		{
			name:          "Hover over primitive let binding",
			input:         "let myVar: int = 5;",
			position:      protocol.Position{Line: 0, Character: 5}, // наведение на 'myVar'
			expectedFound: true,
			expectedText:  "myVar: <int>",
		},
		{
			name:          "Hover over whitespace",
			input:         "let myVar: int = 5;  ",
			position:      protocol.Position{Line: 0, Character: 20}, // пустое пространство
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestDocument(tt.input)
			defer clearTestState()

			params := &protocol.HoverParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
					Position:     tt.position,
				},
			}

			result, err := textDocumentHover(nil, params)
			assertNoError(t, err)

			if !tt.expectedFound {
				if result != nil {
					t.Errorf("expected no hover info, got: %v", result)
				}
				return
			}

			if result == nil {
				t.Fatalf("expected hover info, got nil")
			}

			markup := result.Contents.(protocol.MarkupContent)
			if markup.Kind != protocol.MarkupKindMarkdown {
				t.Errorf("expected markdown kind, got %s", markup.Kind)
			}

			if !strings.Contains(markup.Value, tt.expectedText) {
				t.Errorf("expected hover to contain %q, got %q", tt.expectedText, markup.Value)
			}
		})
	}
}

func TestIsPosInside(t *testing.T) {
	pos := protocol.Position{Line: 1, Character: 5}
	start := token.Position{Line: 1, Column: 0}
	end := token.Position{Line: 1, Column: 10}

	if !isPosInside(pos, start, end) {
		t.Errorf("expected pos to be inside")
	} else {
		t.Logf("pos is inside")
	}
}
