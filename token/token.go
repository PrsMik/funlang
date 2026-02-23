package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	// ключевые слова (выражения)
	LET
	FN
	RETURN

	// идентификаторы и литералы
	IDENT
	INT

	// операторы
	// присваивание
	ASSIGN
	// математика (+, -, *, /)
	PLUS
	MINUS
	MUL
	DIV
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
	NOT

	// спецсимволы
	COMMA
	COLON
	SEMICOLON
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// типы
	INT_TYPE
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	// ключевые слова (выражения)
	"let":    LET,
	"fn":     FN,
	"return": RETURN,
}

var operators = map[string]TokenType{
	// операторы
	// присваивание
	"=": ASSIGN,
	// математика (+, -, *, /)
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	// сравнение (== и !=)
	"==": EQUAL,
	"!=": NOT_EQUAL,
	// сравнение (<, >, <=, >=)
	"<":  LESS,
	">":  GREATER,
	"<=": LESS_OR_EQUAL,
	">=": GREATER_OR_EQUAL,
	// булева алгебра (&&, ||, !)
	"&&": AND,
	"||": OR,
	"!":  NOT,
}

var symbols = map[string]TokenType{
	// спецсимволы
	",": COMMA,
	":": COLON,
	";": SEMICOLON,
	"(": LPAREN,
	")": RPAREN,
	"{": LBRACE,
	"}": RBRACE,
}

var types = map[string]TokenType{
	// типы
	"int": INT_TYPE,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENT
}

func LookupType(identifier string) (TokenType, bool) {
	if tokenType, ok := types[identifier]; ok {
		return tokenType, true
	}
	return IDENT, false
}
