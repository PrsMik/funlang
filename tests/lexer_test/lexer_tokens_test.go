package lexer_test

import (
	"funlang/lexer"
	"funlang/token"
	"testing"
)

type expectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func testTokenStream(t *testing.T, input string, tests []expectedToken) {
	t.Helper()
	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d (literal: %q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_OperatorsAndDelimiters(t *testing.T) {
	input := `= + - * / ! < > == != && || : ; , ( ) { } [ ] ->`

	tests := []expectedToken{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.BANG, "!"},
		{token.LESS, "<"},
		{token.GREATER, ">"},
		{token.EQUAL, "=="},
		{token.NOT_EQUAL, "!="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.COLON, ":"},
		{token.SEMICOLON, ";"},
		{token.COMMA, ","},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.RARROW, "->"},
		{token.EOF, ""},
	}

	testTokenStream(t, input, tests)
}

func TestLexer_KeywordsAndTypes(t *testing.T) {
	input := `let fn return if else true false int bool string`

	tests := []expectedToken{
		{token.LET, "let"},
		{token.FN, "fn"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.INT_TYPE, "int"},
		{token.BOOL_TYPE, "bool"},
		{token.STRING_TYPE, "string"},
		{token.EOF, ""},
	}

	testTokenStream(t, input, tests)
}

func TestLexer_LiteralsAndIdentifiers(t *testing.T) {
	input := `myVar 123 "hello world" "foo"`

	tests := []expectedToken{
		{token.IDENT, "myVar"},
		{token.INT, "123"},
		{token.STRING, "\"hello world\""},
		{token.STRING, "\"foo\""},
		{token.EOF, ""},
	}

	testTokenStream(t, input, tests)
}
