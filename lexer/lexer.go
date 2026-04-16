package lexer

import (
	"fmt"
	"funlang/token"
	"unicode/utf8"
)

type Lexer struct {
	input   string
	curPos  int
	readPos int
	// curChar byte
	curChar rune

	curCol  int
	curLine int
}

func New(input string) *Lexer {
	lexer := Lexer{input: input, curCol: -1, curLine: 0}
	lexer.readChar()
	return &lexer
}

func (lxr *Lexer) NextToken() token.Token {
	var nextTok token.Token
	lxr.skipWhitespace()
	startPos := token.Position{Column: lxr.curCol, Line: lxr.curLine}

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
		var ok bool
		_, ok = newTwoCharToken(lxr, token.COMMENT_SEPARATOR)
		if !ok {
			nextTok = newToken(token.SLASH, '/')
		} else {
			lxr.readChar()
			nextTok = lxr.newCommentToken()
		}
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
	case '[':
		nextTok = newToken(token.LBRACKET, '[')
	case ']':
		nextTok = newToken(token.RBRACKET, ']')
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
		nextTok.Literal = fmt.Sprintf("%q", lxr.readString())
	case 0:
		nextTok = token.Token{Type: token.EOF, Literal: ""}
		nextTok.Start = startPos
		nextTok.End = token.Position{Column: lxr.curCol, Line: lxr.curLine}
		return nextTok
	default:
		if isLetter(lxr.curChar) {
			nextTok.Literal = lxr.readIdentifier()
			var ok bool
			nextTok.Type, ok = token.LookupType(nextTok.Literal)
			if !ok {
				nextTok.Type = token.LookupIdentifier(nextTok.Literal)
			}
		} else if isDigit(lxr.curChar) {

			nextTok.Literal = lxr.readNumber()

			nextTok.Type = token.INT
		} else {
			nextTok = newToken(token.ILLEGAL, lxr.curChar)
			lxr.readChar()
		}

		nextTok.Start = startPos
		nextTok.End = token.Position{Column: lxr.curCol, Line: lxr.curLine}
		return nextTok
	}

	nextTok.Start = startPos

	nextTok.End = token.Position{Column: lxr.curCol + 1, Line: lxr.curLine}

	lxr.readChar()

	return nextTok
}

func (lxr *Lexer) readChar() {
	var size int
	if lxr.readPos >= len(lxr.input) {
		lxr.curChar = 0
	} else {
		// lxr.curChar = lxr.input[lxr.readPos]
		tokenRune, tmp := utf8.DecodeRuneInString(lxr.input[lxr.readPos:])
		size = tmp
		lxr.curChar = tokenRune
		// lxr.readPos += size
	}

	if lxr.curChar == '\n' {
		lxr.curCol = -1
		lxr.curLine++
	} else {
		lxr.curCol++
	}

	lxr.curPos = lxr.readPos
	lxr.readPos += size
	// lxr.readPos++
}
