package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	// ключевые слова (выражения)
	LET
	FN
	RETURN
	TRUE
	FALSE
	IF
	ELSE

	// идентификаторы и литералы
	COMMENT
	IDENT
	INT
	BOOL
	STRING

	// операторы
	// присваивание
	ASSIGN
	LARROW
	RARROW
	// математика (+, -, *, /)
	PLUS
	MINUS
	ASTERISK
	SLASH
	COMMENT_SEPARATOR // "//"
	// сравнение (== и !=)
	EQUAL
	NOT_EQUAL
	// сравнение (<, >, <=, >=)
	LESS
	GREATER
	LESS_OR_EQUAL
	GREATER_OR_EQUAL
	// булева алгебра (&&, ||, !)
	AND
	OR
	BANG

	// спецсимволы
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	// типы
	INT_TYPE
	BOOL_TYPE
	STRING_TYPE
)

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type    TokenType
	Literal string
	Start   Position
	End     Position
}
