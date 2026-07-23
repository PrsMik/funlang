package parser

import (
	"fmt"
	"funlang/token"
)

func (prs *Parser) Errors() []ParseError {
	return prs.errors
}

func (prs *Parser) integerLiteralParseError() {
	msg := fmt.Sprintf("could not parse %q as integer; on pos from %v to %v",
		prs.curToken.Literal, prs.curToken.Start, prs.curToken.End)

	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) tokenError(tknGot token.Token, tknWant token.Token, tknPos token.Token) {
	got, _ := token.LookupString(tknGot.Type)
	want, _ := token.LookupString(tknWant.Type)

	msg := fmt.Sprintf("expected next token to be %s, got %s instead; with value: %s ; on pos from %v to %v",
		want, got, tknGot.Literal, tknPos.Start, tknPos.End)

	err := ParseError{Msg: "parse error: " + msg, Token: tknPos}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) peekError(tknType token.TokenType) {
	got, _ := token.LookupString(prs.peekToken.Type)
	want, _ := token.LookupString(tknType)

	msg := fmt.Sprintf("expected next token to be %s, got %s instead; with value: %s ; on pos from %v to %v",
		want, got, prs.peekToken.Literal, prs.peekToken.Start, prs.peekToken.End)

	err := ParseError{Msg: "parse error: " + msg, Token: prs.peekToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) typeError() {
	got, _ := token.LookupString(prs.curToken.Type)
	msg := fmt.Sprintf("expected type definition got: %s instead; with value: %s ; on pos from %v to %v",
		got, prs.curToken.Literal, prs.curToken.Start, prs.curToken.End)

	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) stmtError() {
	got, _ := token.LookupString(prs.curToken.Type)
	msg := fmt.Sprintf("expected some statement got: %s instead; with value: %s ; on pos from %v to %v",
		got, prs.curToken.Literal, prs.curToken.Start, prs.curToken.End)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) noPrefixParseFnError(tknType token.TokenType) {
	tknStr, _ := token.LookupString(tknType)
	msg := fmt.Sprintf("no prefix parse function for %s found; on pos from %v to %v",
		tknStr, prs.curToken.Start, prs.curToken.End)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}
