package lexer_test

import (
	"funlang/lexer"
	"funlang/token"
	"testing"
)

func TestLexer_CompleteProgram(t *testing.T) {
	input := `let five: int = 5;
let ten = 10; //  абоба
let сумм = fn(x, y) {
	return x + y;
};
"foobar"
"foo bar"
[1, 2];`

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
		{token.COMMENT, "абоба"},

		{token.LET, "let"},
		{token.IDENT, "сумм"},
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

		{token.STRING, "\"foobar\""},
		{token.STRING, "\"foo bar\""},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		got := l.NextToken()
		if got.Type != tt.wantType {
			t.Errorf("Test[%d]: Type wrong. expected=%d, got=%d", i, tt.wantType, got.Type)
		}
		if got.Literal != tt.wantLiteral {
			t.Errorf("Test[%d]: Literal wrong. expected=%q, got=%q", i, tt.wantLiteral, got.Literal)
		}
	}
}
