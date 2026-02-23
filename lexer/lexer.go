package lexer

import "funlang/token"

type Lexer struct {
	input   string
	curPos  int
	readPos int
	curChar byte
}

func New(input string) *Lexer {
	lexer := Lexer{input: input}
	lexer.readChar()
	return &lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var nextTok token.Token
	lexer.skipWhitespace()
	switch lexer.curChar {
	case '=':
		nextTok = newToken(token.ASSIGN, '=')
	case '+':
		nextTok = newToken(token.PLUS, '+')
	case '-':
		nextTok = newToken(token.MINUS, '-')
	case '*':
		nextTok = newToken(token.MUL, '*')
	case '/':
		nextTok = newToken(token.DIV, '/')
	case ',':
		nextTok = newToken(token.COMMA, ',')
	case ':':
		nextTok = newToken(token.COLON, ':')
	case ';':
		nextTok = newToken(token.SEMICOLON, ';')
	case '(':
		nextTok = newToken(token.LPAREN, '(')
	case ')':
		nextTok = newToken(token.RPAREN, ')')
	case '{':
		nextTok = newToken(token.LBRACE, '{')
	case '}':
		nextTok = newToken(token.RBRACE, '}')
	case 0:
		nextTok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(lexer.curChar) {
			nextTok.Literal = lexer.readIdentifier()
			var ok bool
			nextTok.Type, ok = token.LookupType(nextTok.Literal)
			if !ok {
				nextTok.Type = token.LookupIdentifier(nextTok.Literal)
			}
			return nextTok
		} else if isDigit(lexer.curChar) {
			nextTok.Literal = lexer.readNumber()
			nextTok.Type = token.INT
			return nextTok
		} else {
			nextTok = newToken(token.ILLEGAL, lexer.curChar)
		}
	}
	lexer.readChar()
	return nextTok
}

func (lexer *Lexer) readChar() {
	if lexer.readPos >= len(lexer.input) {
		lexer.curChar = 0
	} else {
		lexer.curChar = lexer.input[lexer.readPos]
	}
	lexer.curPos = lexer.readPos
	lexer.readPos++
}

func newToken(tokenType token.TokenType, tokenChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenChar)}
}

func (lexer *Lexer) readIdentifier() string {
	startPos := lexer.curPos
	for isLetter(lexer.curChar) {
		lexer.readChar()
	}
	return lexer.input[startPos:lexer.curPos]
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.curChar == ' ' || lexer.curChar == '\t' || lexer.curChar == '\n' || lexer.curChar == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readNumber() string {
	startPos := lexer.curPos
	for isDigit(lexer.curChar) {
		lexer.readChar()
	}
	return lexer.input[startPos:lexer.curPos]
}
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
