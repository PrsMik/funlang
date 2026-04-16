package parser_test

import (
	"funlang/lexer"
	"funlang/parser"
	"strings"
	"testing"
)

func TestParserErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Missing identifier in let",
			"let : int = 5;",
			"expected next token to be IDENT",
		},
		{
			"Missing type in let",
			"let x = 5;",
			"expected next token to be COLON",
		},
		{
			"Invalid type definition",
			"let x: 123 = 5;",
			"expected type definition got: INT",
		},
		{
			"Missing semicolon",
			"let x: int = 5",
			"expected next token to be SEMICOLON",
		},
		{
			"Mismatched parentheses in expression",
			"let x: int = (5 + 5;",
			"expected next token to be RPAREN",
		},
		{
			"Invalid prefix operator",
			"let x: int = *5;",
			"no prefix parse function for ASTERISK found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			p.ParseProgram()

			errors := p.Errors()
			if len(errors) == 0 {
				t.Fatalf("Parser expected to have errors for input: %q, but got none", tt.input)
			}

			found := false
			for _, err := range errors {
				if strings.Contains(err.Msg, tt.expected) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Error message %q not found in parser errors: %v", tt.expected, errors)
			}
		})
	}
}
