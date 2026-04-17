package token_test

import (
	"funlang/token"
	"testing"
)

func TestLookupIdentifier(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		expected   token.TokenType
	}{
		{"Keyword let", "let", token.LET},
		{"Keyword fn", "fn", token.FN},
		{"Keyword return", "return", token.RETURN},
		{"Keyword if", "if", token.IF},
		{"Keyword else", "else", token.ELSE},
		{"Keyword true", "true", token.TRUE},
		{"Keyword false", "false", token.FALSE},
		{"User defined identifier", "myVariable", token.IDENT},
		{"Single letter identifier", "x", token.IDENT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := token.LookupIdentifier(tt.identifier)
			if actual != tt.expected {
				t.Errorf("LookupIdentifier(%q) = %v; want %v", tt.identifier, actual, tt.expected)
			}
		})
	}
}

func TestLookupType(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		expected   token.TokenType
		expectedOk bool
	}{
		{"Int type", "int", token.INT_TYPE, true},
		{"Bool type", "bool", token.BOOL_TYPE, true},
		{"String type", "string", token.STRING_TYPE, true},
		{"Not a type keyword", "let", token.IDENT, false},
		{"Random identifier", "myType", token.IDENT, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := token.LookupType(tt.identifier)
			if actual != tt.expected || ok != tt.expectedOk {
				t.Errorf("LookupType(%q) = (%v, %v); want (%v, %v)",
					tt.identifier, actual, ok, tt.expected, tt.expectedOk)
			}
		})
	}
}

func TestLookupOperator(t *testing.T) {
	tests := []struct {
		name       string
		literal    string
		expected   token.TokenType
		expectedOk bool
	}{
		{"Assign operator", "=", token.ASSIGN, true},
		{"Plus operator", "+", token.PLUS, true},
		{"Equal operator", "==", token.EQUAL, true},
		{"Not equal operator", "!=", token.NOT_EQUAL, true},
		{"Comment separator", "//", token.COMMENT_SEPARATOR, true},
		{"Not an operator", "foo", token.ILLEGAL, false},
		{"Empty string", "", token.ILLEGAL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := token.LookupOperator(tt.literal)
			if actual != tt.expected || ok != tt.expectedOk {
				t.Errorf("LookupOperator(%q) = (%v, %v); want (%v, %v)",
					tt.literal, actual, ok, tt.expected, tt.expectedOk)
			}
		})
	}
}
