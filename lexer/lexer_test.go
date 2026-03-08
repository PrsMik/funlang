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
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.BOOL_TYPE, "bool"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
	}
	lexer := lexer.New(content)
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
