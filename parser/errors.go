package parser

import (
	"fmt"
	"funlang/token"
)

func (prs *Parser) Errors() []ParseError {
	return prs.errors
}

func (prs *Parser) integerLiteralParseError() {
	msg := fmt.Sprintf("could not parse %q as integer", prs.curToken.Literal)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) tokenError(tknGot token.Token, tknWant token.Token, tknPos token.Token) {
	got, _ := token.LookupString(tknGot.Type)
	want, _ := token.LookupString(tknWant.Type)
	msg := fmt.Sprintf("expected next token to be %s, got %s instead; with value: %s", want, got, tknGot.Literal)
	err := ParseError{Msg: "parse error: " + msg, Token: tknPos}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) peekError(tknType token.TokenType) {
	want, _ := token.LookupString(tknType)
	got, _ := token.LookupString(prs.peekToken.Type)
	msg := fmt.Sprintf("expected next token to be %s, got %s instead; with value: %s", want, got, prs.peekToken.Literal)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.peekToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) typeError() {
	got, _ := token.LookupString(prs.curToken.Type)
	msg := fmt.Sprintf("expected type definition got: %s instead; with value: %s", got, prs.curToken.Literal)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) stmtError() {
	got, _ := token.LookupString(prs.curToken.Type)
	msg := fmt.Sprintf("expected some statement got: %s instead; with value: %s", got, prs.curToken.Literal)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}

func (prs *Parser) noPrefixParseFnError(tknType token.TokenType) {
	tknStr, _ := token.LookupString(tknType)
	msg := fmt.Sprintf("no prefix parse function for %s found", tknStr)
	err := ParseError{Msg: "parse error: " + msg, Token: prs.curToken}
	prs.errors = append(prs.errors, err)
}
