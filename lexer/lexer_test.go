package lexer_test

import (
	"fmt"
	"funlang/lexer"
	"funlang/token"
	"os"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	data, err := os.ReadFile("D:/Repo/funlang/main2.txt")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	content := string(data)
	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FN, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.INT_TYPE, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.COLON, ":"},
		{token.INT_TYPE, "int"},
		{token.RPAREN, ")"},
		{token.RARROW, "->"},
		{token.INT_TYPE, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.LARROW, "<-"},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	lexer := lexer.New(content)
	for _, tt := range tests {
		t.Run("TEST", func(t *testing.T) {
			got := lexer.NextToken()
			if got.Type != tt.wantType {
				t.Errorf("NextToken() Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Literal != tt.wantLiteral {
				t.Errorf("NextToken() Literal = %v, want %v", got.Literal, tt.wantLiteral)
			}
		})
	}
}
