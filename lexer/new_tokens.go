package lexer

import (
	"funlang/token"
	"unicode/utf8"
)

func newToken(tokenType token.TokenType, tokenRune rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenRune)}
}

// совершает поглощение следующего символа и проверяет является ли он
// корректным двухсимвольным токеном, если да - возвращает его, сохраняя позицию,
// если нет - возвращает ILLEGAL токен и false, откатывая позицию
func newTwoCharToken(lexer *Lexer, tokenType token.TokenType) (token.Token, bool) {
	char := lexer.curChar
	pos := lexer.curPos
	startCol := lexer.curCol

	lexer.readChar()

	literal := string(char) + string(lexer.curChar)

	wantType, ok := token.LookupOperator(literal)
	if !ok || wantType != tokenType {
		lexer.readPos = lexer.curPos
		lexer.curPos = pos
		lexer.curCol = startCol
		return token.Token{Type: token.ILLEGAL, Literal: literal}, false
	}

	token := token.Token{Type: tokenType, Literal: literal}
	return token, true
}

func (lxr *Lexer) newCommentToken() token.Token {
	var nextTok token.Token
	var startPos int

	lxr.skipCommentSpaces()

	if lxr.curChar != '\n' && lxr.curChar != 0 {
		startPos = lxr.curPos
		nextTok.Type = token.COMMENT
		nextTok.Start = token.Position{Column: lxr.curCol, Line: lxr.curLine}

		for lxr.peekChar() != '\n' && lxr.peekChar() != 0 {
			lxr.readChar()
		}

		nextTok.Literal = lxr.input[startPos : lxr.curPos+utf8.RuneLen(lxr.curChar)]

		nextTok.End = token.Position{Column: lxr.curCol, Line: lxr.curLine}
	} else {
		startPos = lxr.curPos
		nextTok.Type = token.COMMENT
		nextTok.Start = token.Position{Column: lxr.curCol, Line: lxr.curLine}
		nextTok.Literal = ""
		nextTok.End = token.Position{Column: lxr.curCol, Line: lxr.curLine}
	}

	return nextTok
}
