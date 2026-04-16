package ast_test

import (
	"funlang/ast"
	"funlang/token"
	"testing"
)

func TestStatements(t *testing.T) {
	tests := []struct {
		name                 string
		node                 ast.StatementNode
		expectedTokenLiteral string
		expectedString       string
		expectedStart        token.Position
		expectedEnd          token.Position
	}{
		{
			name: "LetStatement (with value)",
			node: &ast.LetStatement{
				Token:     createToken(token.LET, "let", 1, 1),
				Name:      &ast.Identifier{Value: "x"},
				Type:      &ast.SimpleType{Value: "int"},
				Value:     &ast.IntegerLiteral{Token: createToken(token.INT, "5", 1, 13), Value: 5},
				SemiToken: createToken(token.SEMICOLON, ";", 1, 14),
			},
			expectedTokenLiteral: "let",
			expectedString:       "let x: int = 5;",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 15),
		},
		{
			name: "LetStatement (without value)",
			node: &ast.LetStatement{
				Token:     createToken(token.LET, "let", 1, 1),
				Name:      &ast.Identifier{Value: "x"},
				Type:      &ast.SimpleType{Value: "int"},
				Value:     nil,
				SemiToken: createToken(token.SEMICOLON, ";", 1, 10),
			},
			expectedTokenLiteral: "let",
			expectedString:       "let x: int = ;",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 11),
		},
		{
			name: "ReturnStatement (with value)",
			node: &ast.ReturnStatement{
				Token:     createToken(token.RETURN, "return", 1, 1),
				Value:     &ast.IntegerLiteral{Token: createToken(token.INT, "5", 1, 8), Value: 5},
				SemiToken: createToken(token.SEMICOLON, ";", 1, 9),
			},
			expectedTokenLiteral: "return",
			expectedString:       "return 5;",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 10),
		},
		{
			name: "ReturnStatement (without value)",
			node: &ast.ReturnStatement{
				Token:     createToken(token.RETURN, "return", 1, 1),
				Value:     nil,
				SemiToken: createToken(token.SEMICOLON, ";", 1, 7),
			},
			expectedTokenLiteral: "return",
			expectedString:       "return ;",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 8),
		},
		{
			name: "BlockStatement",
			node: &ast.BlockStatement{
				Token: createToken(token.LBRACE, "{", 1, 1),
				Statements: []ast.StatementNode{
					&ast.ReturnStatement{
						Token: createToken(token.RETURN, "return", 1, 3),
					},
				},
				SemiToken: createToken(token.RBRACE, "}", 1, 10),
			},
			expectedTokenLiteral: "{",
			expectedString:       "return ;",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 11),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.node.statementNode()

			if got := tt.node.TokenLiteral(); got != tt.expectedTokenLiteral {
				t.Errorf("TokenLiteral() wrong. expected=%q, got=%q", tt.expectedTokenLiteral, got)
			}
			if got := tt.node.String(); got != tt.expectedString {
				t.Errorf("String() wrong. expected=%q, got=%q", tt.expectedString, got)
			}
			if got := tt.node.Start(); got != tt.expectedStart {
				t.Errorf("Start() wrong. expected=%v, got=%v", tt.expectedStart, got)
			}
			if got := tt.node.End(); got != tt.expectedEnd {
				t.Errorf("End() wrong. expected=%v, got=%v", tt.expectedEnd, got)
			}
		})
	}
}
