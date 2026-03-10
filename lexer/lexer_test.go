package lexer_test

import (
	"funlang/lexer"
	"funlang/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	input := `let five: int = 5;
			  let ten = 10;
			  let add = fn(x, y) {
				return x + y;
				};
				"foobar"
				"foo bar"`
	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.COLON, ":"},
		{token.INT_TYPE, "int"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FN, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
	}
	lexer := lexer.New(input)
	for _, tt := range tests {
		t.Run("TEST", func(t *testing.T) {
			got := lexer.NextToken()
			if got.Type != tt.wantType {
				t.Errorf("NextToken() Type = %v and want %v", got.Type, tt.wantType)
			}
			if got.Literal != tt.wantLiteral {
				t.Errorf("NextToken() Literal = %v and want %v", got.Literal, tt.wantLiteral)
			}
		})
	}
}
