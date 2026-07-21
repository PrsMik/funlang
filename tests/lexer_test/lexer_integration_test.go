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

func TestLexer_TypesIntegration(t *testing.T) {
	input := `let Result: fn(T: type, E: type) -> type = fn(T: type, E: type) {
    return OkType(T) | ErrType(E); 
};

let error_text: string = switch res.(type) {
    case ErrType(string) {
        return res.error; 
    }
	default {
		return res.error;
	}
};`

	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		// let Result: fn(T: type, E: type) -> type = fn(T: type, E: type) {
		{token.LET, "let"},
		{token.IDENT, "Result"},
		{token.COLON, ":"},
		{token.FN, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "T"},
		{token.COLON, ":"},
		{token.TYPE, "type"},
		{token.COMMA, ","},
		{token.IDENT, "E"},
		{token.COLON, ":"},
		{token.TYPE, "type"},
		{token.RPAREN, ")"},
		{token.RARROW, "->"},
		{token.TYPE, "type"},
		{token.ASSIGN, "="},
		{token.FN, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "T"},
		{token.COLON, ":"},
		{token.TYPE, "type"},
		{token.COMMA, ","},
		{token.IDENT, "E"},
		{token.COLON, ":"},
		{token.TYPE, "type"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		// return OkType(T) | ErrType(E);
		{token.RETURN, "return"},
		{token.IDENT, "OkType"},
		{token.LPAREN, "("},
		{token.IDENT, "T"},
		{token.RPAREN, ")"},
		{token.PIPE, "|"}, // Union Type
		{token.IDENT, "ErrType"},
		{token.LPAREN, "("},
		{token.IDENT, "E"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		// };
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		// let error_text: string = switch res.(type) {
		{token.LET, "let"},
		{token.IDENT, "error_text"},
		{token.COLON, ":"},
		{token.STRING_TYPE, "string"},
		{token.ASSIGN, "="},
		{token.SWITCH, "switch"},
		{token.IDENT, "res"},
		{token.DOT, "."},
		{token.LPAREN, "("},
		{token.TYPE, "type"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		// case ErrType(string) {
		{token.CASE, "case"},
		{token.IDENT, "ErrType"},
		{token.LPAREN, "("},
		{token.STRING_TYPE, "string"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		// return res.error;
		{token.RETURN, "return"},
		{token.IDENT, "res"},
		{token.DOT, "."}, // Member Access
		{token.IDENT, "error"},
		{token.SEMICOLON, ";"},

		// }
		{token.RBRACE, "}"},

		// default {
		{token.DEFAULT, "default"},
		{token.LBRACE, "{"},

		// return res.error;
		{token.RETURN, "return"},
		{token.IDENT, "res"},
		{token.DOT, "."}, // Member Access
		{token.IDENT, "error"},
		{token.SEMICOLON, ";"},

		// }
		{token.RBRACE, "}"},

		// };
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		got := l.NextToken()
		if got.Type != tt.wantType {
			t.Errorf("Test[%d]: Type wrong. expected=%v, got=%v (literal: %q)", i, tt.wantType, got.Type, got.Literal)
		}
		if got.Literal != tt.wantLiteral {
			t.Errorf("Test[%d]: Literal wrong. expected=%q, got=%q", i, tt.wantLiteral, got.Literal)
		}
	}
}
