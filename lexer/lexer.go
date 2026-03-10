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

func (lxr *Lexer) NextToken() token.Token {
	var nextTok token.Token
	lxr.skipWhitespace()
	switch lxr.curChar {
	case '=':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.EQUAL)
		if !ok {
			nextTok = newToken(token.ASSIGN, '=')
		}
	case '+':
		nextTok = newToken(token.PLUS, '+')
	case '-':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.RARROW)
		if !ok {
			nextTok = newToken(token.MINUS, '-')
		}
	case '*':
		nextTok = newToken(token.ASTERISK, '*')
	case '/':
		nextTok = newToken(token.SLASH, '/')
	case '<':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.LESS_OR_EQUAL)
		if !ok {
			if nextTok.Literal == "<-" {
				nextTok.Type = token.LARROW
			} else {
				nextTok = newToken(token.LESS, '<')
			}
		}
	case '>':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.GREATER_OR_EQUAL)
		if !ok {
			nextTok = newToken(token.GREATER, '>')
		}
	case '!':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.NOT_EQUAL)
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
	case '&':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.AND)
		if !ok {
			nextTok = newToken(token.ILLEGAL, lxr.curChar)
		}
	case '|':
		var ok bool
		nextTok, ok = newTwoCharToken(lxr, token.OR)
		if !ok {
			nextTok = newToken(token.ILLEGAL, lxr.curChar)
		}
	case '"':
		nextTok.Type = token.STRING
		nextTok.Literal = lxr.readString()
	case 0:
		nextTok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(lxr.curChar) {
			nextTok.Literal = lxr.readIdentifier()
			var ok bool
			nextTok.Type, ok = token.LookupType(nextTok.Literal)
			if !ok {
				nextTok.Type = token.LookupIdentifier(nextTok.Literal)
			}
			return nextTok
		} else if isDigit(lxr.curChar) {
			nextTok.Literal = lxr.readNumber()
			nextTok.Type = token.INT
			return nextTok
		} else {
			nextTok = newToken(token.ILLEGAL, lxr.curChar)
		}
	}
	lxr.readChar()
	return nextTok
}

func (lxr *Lexer) readChar() {
	if lxr.readPos >= len(lxr.input) {
		lxr.curChar = 0
	} else {
		lxr.curChar = lxr.input[lxr.readPos]
	}
	lxr.curPos = lxr.readPos
	lxr.readPos++
}

func (lxr *Lexer) peekChar() byte {
	if lxr.readPos >= len(lxr.input) {
		return 0
	} else {
		return lxr.input[lxr.readPos]
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

func (lxr *Lexer) readIdentifier() string {
	startPos := lxr.curPos
	for isLetter(lxr.curChar) {
		lxr.readChar()
	}
	return lxr.input[startPos:lxr.curPos]
}

func (lxr *Lexer) readString() string {
	position := lxr.curPos + 1
	for {
		lxr.readChar()
		if lxr.curChar == '"' || lxr.curChar == 0 {
			break
		}
	}
	return lxr.input[position:lxr.curPos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (lxr *Lexer) skipWhitespace() {
	for lxr.curChar == ' ' || lxr.curChar == '\t' || lxr.curChar == '\n' || lxr.curChar == '\r' {
		lxr.readChar()
	}
}

func (lxr *Lexer) readNumber() string {
	startPos := lxr.curPos
	for isDigit(lxr.curChar) {
		lxr.readChar()
	}
	return lxr.input[startPos:lxr.curPos]
}
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
