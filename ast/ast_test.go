package ast

import (
	"funlang/token"
	"testing"
)

func TestProgram(t *testing.T) {
	mockStatement := &LetStatement{
		Token: createToken(token.LET, "let", 1, 1),
		Name: &Identifier{
			Token: createToken(token.IDENT, "myVar", 1, 5),
			Value: "myVar",
		},
		Type: &SimpleType{
			Token: createToken(token.INT_TYPE, "int", 1, 12),
			Value: "int",
		},
		Value: &Identifier{
			Token: createToken(token.IDENT, "anotherVar", 1, 18),
			Value: "anotherVar",
		},
		SemiToken: createToken(token.SEMICOLON, ";", 1, 28),
	}

	tests := []struct {
		name                 string
		program              *Program
		expectedString       string
		expectedTokenLiteral string
	}{
		{
			name: "Non-empty program",
			program: &Program{
				Statements: []StatementNode{mockStatement},
			},
			expectedString:       "let myVar: int = anotherVar;",
			expectedTokenLiteral: "let",
		},
		{
			name: "Empty program",
			program: &Program{
				Statements: []StatementNode{},
			},
			expectedString:       "",
			expectedTokenLiteral: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.program.String(); got != tt.expectedString {
				t.Errorf("Program.String() wrong. expected=%q, got=%q", tt.expectedString, got)
			}
			if got := tt.program.TokenLiteral(); got != tt.expectedTokenLiteral {
				t.Errorf("Program.TokenLiteral() wrong. expected=%q, got=%q", tt.expectedTokenLiteral, got)
			}

			if len(tt.program.Statements) > 0 {
				expectedStart := pos(1, 1)
				expectedEnd := pos(1, 29) // 28 + len(";") = 29
				if got := tt.program.Start(); got != expectedStart {
					t.Errorf("Program.Start() wrong. expected=%v, got=%v", expectedStart, got)
				}
				if got := tt.program.End(); got != expectedEnd {
					t.Errorf("Program.End() wrong. expected=%v, got=%v", expectedEnd, got)
				}
			}
		})
	}
}
