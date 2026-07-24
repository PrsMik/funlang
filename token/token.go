package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	// ключевые слова (выражения)
	LET       // let
	FN        // fn
	RETURN    // return
	TRUE      // true
	FALSE     // false
	IF        // if
	ELSE      // else
	SWITCH    // switch
	CASE      // case
	DEFAULT   // default
	INTERFACE // interface
	IMPL      // impl
	TYPE      // type

	// идентификаторы и литералы

	COMMENT // // myComment
	IDENT   // myVar
	INT     // 123
	STRING  // "123"

	// ОПЕРАТОРЫ
	// присваивание

	ASSIGN // =
	LARROW // ->
	RARROW // <-

	// математика (+, -, *, /)

	PLUS              // +
	MINUS             // -
	ASTERISK          // *
	SLASH             // /
	COMMENT_SEPARATOR // //

	// сравнение (== и !=)

	EQUAL     // ==
	NOT_EQUAL // !=

	// сравнение (<, >, <=, >=)

	LESS             // <
	GREATER          // >
	LESS_OR_EQUAL    // <=
	GREATER_OR_EQUAL // >=

	// булева алгебра (&&, ||, !)

	TWO_AMPERSANDS // &&
	TWO_PIPES      // ||
	BANG           // !

	// типы (|, &, .)

	PIPE      // |
	AMPERSAND // &
	DOT       // .

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

	INT_TYPE    // int
	BOOL_TYPE   // bool
	STRING_TYPE // string
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
