package lexer_test

import (
	"funlang/lexer"
	"funlang/token"
	"testing"
)

func TestLexer_TokenBoundaries(t *testing.T) {
	input := `let x: int = 52;
let x: fn(bool) -> bool = fn(y) {
  return y && true;
}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedStart   token.Position
		expectedEnd     token.Position
	}{
		// строка 0
		{token.LET, "let", token.Position{Line: 0, Column: 0}, token.Position{Line: 0, Column: 3}},
		{token.IDENT, "x", token.Position{Line: 0, Column: 4}, token.Position{Line: 0, Column: 5}},
		{token.COLON, ":", token.Position{Line: 0, Column: 5}, token.Position{Line: 0, Column: 6}},
		{token.INT_TYPE, "int", token.Position{Line: 0, Column: 7}, token.Position{Line: 0, Column: 10}},
		{token.ASSIGN, "=", token.Position{Line: 0, Column: 11}, token.Position{Line: 0, Column: 12}},
		{token.INT, "52", token.Position{Line: 0, Column: 13}, token.Position{Line: 0, Column: 15}},
		{token.SEMICOLON, ";", token.Position{Line: 0, Column: 15}, token.Position{Line: 0, Column: 16}},

		// строка 1
		{token.LET, "let", token.Position{Line: 1, Column: 0}, token.Position{Line: 1, Column: 3}},
		{token.IDENT, "x", token.Position{Line: 1, Column: 4}, token.Position{Line: 1, Column: 5}},
		{token.COLON, ":", token.Position{Line: 1, Column: 5}, token.Position{Line: 1, Column: 6}},
		{token.FN, "fn", token.Position{Line: 1, Column: 7}, token.Position{Line: 1, Column: 9}},
		{token.LPAREN, "(", token.Position{Line: 1, Column: 9}, token.Position{Line: 1, Column: 10}},
		{token.BOOL_TYPE, "bool", token.Position{Line: 1, Column: 10}, token.Position{Line: 1, Column: 14}},
		{token.RPAREN, ")", token.Position{Line: 1, Column: 14}, token.Position{Line: 1, Column: 15}},
		{token.RARROW, "->", token.Position{Line: 1, Column: 16}, token.Position{Line: 1, Column: 18}},
		{token.BOOL_TYPE, "bool", token.Position{Line: 1, Column: 19}, token.Position{Line: 1, Column: 23}},
		{token.ASSIGN, "=", token.Position{Line: 1, Column: 24}, token.Position{Line: 1, Column: 25}},
		{token.FN, "fn", token.Position{Line: 1, Column: 26}, token.Position{Line: 1, Column: 28}},
		{token.LPAREN, "(", token.Position{Line: 1, Column: 28}, token.Position{Line: 1, Column: 29}},
		{token.IDENT, "y", token.Position{Line: 1, Column: 29}, token.Position{Line: 1, Column: 30}},
		{token.RPAREN, ")", token.Position{Line: 1, Column: 30}, token.Position{Line: 1, Column: 31}},
		{token.LBRACE, "{", token.Position{Line: 1, Column: 32}, token.Position{Line: 1, Column: 33}},

		// строка 2
		{token.RETURN, "return", token.Position{Line: 2, Column: 2}, token.Position{Line: 2, Column: 8}},
		{token.IDENT, "y", token.Position{Line: 2, Column: 9}, token.Position{Line: 2, Column: 10}},
		{token.AND, "&&", token.Position{Line: 2, Column: 11}, token.Position{Line: 2, Column: 13}},
		{token.TRUE, "true", token.Position{Line: 2, Column: 14}, token.Position{Line: 2, Column: 18}},
		{token.SEMICOLON, ";", token.Position{Line: 2, Column: 18}, token.Position{Line: 2, Column: 19}},

		// строка 3
		{token.RBRACE, "}", token.Position{Line: 3, Column: 0}, token.Position{Line: 3, Column: 1}},
		{token.EOF, "", token.Position{Line: 3, Column: 1}, token.Position{Line: 3, Column: 1}},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Start != tt.expectedStart {
			t.Errorf("tests[%d] - START position wrong. expected={Line: %d, Column: %d}, got={Line: %d, Column: %d}",
				i, tt.expectedStart.Line, tt.expectedStart.Column, tok.Start.Line, tok.Start.Column)
		}

		if tok.End != tt.expectedEnd {
			t.Errorf("tests[%d] - END position wrong. expected={Line: %d, Column: %d}, got={Line: %d, Column: %d}",
				i, tt.expectedEnd.Line, tt.expectedEnd.Column, tok.End.Line, tok.End.Column)
		}
	}
}
