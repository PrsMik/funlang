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
	COMMA
	COLON
	SEMICOLON
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// типы
	INT_TYPE
	BOOL_TYPE
	STRING_TYPE
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
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

var operators = map[string]TokenType{
	// операторы
	// присваивание
	"=":  ASSIGN,
	"<-": LARROW,
	"->": RARROW,

	// математика (+, -, *, /)
	"+": PLUS,
	"-": MINUS,
	"*": ASTERISK,
	"/": SLASH,
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
	"!":  BANG,
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
	"int":    INT_TYPE,
	"bool":   BOOL_TYPE,
	"string": STRING_TYPE,
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

func LookupOperator(literal string) (TokenType, bool) {
	if tokenType, ok := operators[literal]; ok {
		return tokenType, true
	}
	return ILLEGAL, false
}

func LookupString(tokenType TokenType) (string, bool) {
	if tokenString, ok := tokenStrings[tokenType]; ok {
		return tokenString, true
	}
	return "UNKNOWN_TOKEN", false
}

var tokenStrings = map[TokenType]string{
	ILLEGAL:          "ILLEGAL",
	EOF:              "EOF",
	LET:              "LET",
	FN:               "FN",
	RETURN:           "RETURN",
	TRUE:             "TRUE",
	FALSE:            "FALSE",
	IF:               "IF",
	ELSE:             "ELSE",
	IDENT:            "IDENT",
	INT:              "INT",
	BOOL:             "BOOL",
	STRING:           "STRING",
	ASSIGN:           "ASSIGN",
	LARROW:           "LARROW",
	RARROW:           "RARROW",
	PLUS:             "PLUS",
	MINUS:            "MINUS",
	ASTERISK:         "ASTERISK",
	SLASH:            "SLASH",
	EQUAL:            "EQUAL",
	NOT_EQUAL:        "NOTEQAL",
	LESS:             "LESS",
	GREATER:          "GREATER",
	LESS_OR_EQUAL:    "LESS_OR_EQUAL",
	GREATER_OR_EQUAL: "GREATER_OR_EQUAL",
	AND:              "AND",
	OR:               "OR",
	BANG:             "BANG",
	COMMA:            "COMMA",
	COLON:            "COLON",
	SEMICOLON:        "SEMICOLON",
	LPAREN:           "LPAREN",
	RPAREN:           "RPAREN",
	LBRACE:           "LBRACE",
	RBRACE:           "RBRACE",
	INT_TYPE:         "INT_TYPE",
	BOOL_TYPE:        "BOOL_TYPE",
	STRING_TYPE:      "STRING_TYPE",
}
