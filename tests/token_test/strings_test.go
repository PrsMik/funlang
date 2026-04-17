package token_test

import (
	"funlang/token"
	"testing"
)

func TestLookupString(t *testing.T) {
	tests := []struct {
		name       string
		tokenType  token.TokenType
		expected   string
		expectedOk bool
	}{
		{"Let token", token.LET, "LET", true},
		{"Ident token", token.IDENT, "IDENT", true},
		{"Assign token", token.ASSIGN, "ASSIGN", true},
		{"Int type token", token.INT_TYPE, "INT_TYPE", true},
		{"Illegal token", token.ILLEGAL, "ILLEGAL", true},
		{"EOF token", token.EOF, "EOF", true},
		{"Unknown token", token.TokenType(99999), "UNKNOWN_TOKEN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := token.LookupString(tt.tokenType)
			if actual != tt.expected || ok != tt.expectedOk {
				t.Errorf("LookupString(%v) = (%q, %v); want (%q, %v)",
					tt.tokenType, actual, ok, tt.expected, tt.expectedOk)
			}
		})
	}
}

func TestTokenMapsCompleteness(t *testing.T) {
	mapsToTest := []struct {
		mapName  string
		tokenMap map[string]token.TokenType
	}{
		{"Keywords", token.Keywords},
		{"Operators", token.Operators},
		{"Symbols", token.Symbols},
		{"Types", token.Types},
	}

	for _, mt := range mapsToTest {
		t.Run(mt.mapName, func(t *testing.T) {
			for literal, tokType := range mt.tokenMap {
				str, ok := token.LookupString(tokType)

				if !ok {
					t.Errorf("tokenType %v (literal %q) is missing in tokenStrings map", tokType, literal)
				}

				if str == "UNKNOWN_TOKEN" || str == "" {
					t.Errorf("tokenType %v (literal %q) has invalid string representation: %q", tokType, literal, str)
				}
			}
		})
	}
}
