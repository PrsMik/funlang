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
		var ok bool
		nextTok, ok = newTwoCharToken(lexer, token.EQUAL)
		if !ok {
			nextTok = newToken(token.ASSIGN, '=')
		}
	case '+':
		nextTok = newToken(token.PLUS, '+')
	case '-':
		var ok bool
		nextTok, ok = newTwoCharToken(lexer, token.RARROW)
		if !ok {
			nextTok = newToken(token.MINUS, '-')
		}
	case '*':
		nextTok = newToken(token.ASTERISK, '*')
	case '/':
		nextTok = newToken(token.SLASH, '/')
	case '<':
		var ok bool
		nextTok, ok = newTwoCharToken(lexer, token.LESS_OR_EQUAL)
		if !ok {
			if nextTok.Literal == "<-" {
				nextTok.Type = token.LARROW
			} else {
				nextTok = newToken(token.LESS, '<')
			}
		}
	case '>':
		var ok bool
		nextTok, ok = newTwoCharToken(lexer, token.GREATER_OR_EQUAL)
		if !ok {
			nextTok = newToken(token.GREATER, '>')
		}
	case '!':
		var ok bool
		nextTok, ok = newTwoCharToken(lexer, token.NOT_EQUAL)
		if !ok {
			nextTok = newToken(token.BANG, '!')
		}
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

func (lexer *Lexer) peekChar() byte {
	if lexer.readPos >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPos]
	}
}

func newToken(tokenType token.TokenType, tokenChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenChar)}
}

func newTwoCharToken(lexer *Lexer, tokenType token.TokenType) (token.Token, bool) {
	char := lexer.curChar
	pos := lexer.curPos
	lexer.readChar()
	literal := string(char) + string(lexer.curChar)
	wantType, ok := token.LookupOperator(literal)
	if !ok || wantType != tokenType {
		lexer.readPos = lexer.curPos
		lexer.curPos = pos
		return token.Token{Type: token.ILLEGAL, Literal: literal}, false
	}
	token := token.Token{Type: tokenType, Literal: literal}
	return token, true
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
