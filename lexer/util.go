package lexer

import (
	"unicode"
	"unicode/utf8"
)

func (lxr *Lexer) peekChar() rune {
	if lxr.readPos >= len(lxr.input) {
		return 0
	} else {
		tknRune, _ := utf8.DecodeRuneInString(lxr.input[lxr.readPos:])
		return tknRune
	}
}

func (lxr *Lexer) readIdentifier() string {
	startPos := lxr.curPos
	for isLetter(lxr.curChar) || isDigit(lxr.curChar) {
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

func (lxr *Lexer) readNumber() string {
	startPos := lxr.curPos
	for isDigit(lxr.curChar) {
		lxr.readChar()
	}
	return lxr.input[startPos:lxr.curPos]
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (lxr *Lexer) skipWhitespace() {
	for lxr.curChar == ' ' || lxr.curChar == '\t' || lxr.curChar == '\n' || lxr.curChar == '\r' {
		lxr.readChar()
	}
}

func (lxr *Lexer) skipCommentSpaces() {
	for lxr.curChar == ' ' || lxr.curChar == '\t' || lxr.curChar == '\r' {
		lxr.readChar()
	}
}
