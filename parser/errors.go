package parser

import (
	"fmt"
	"funlang/token"
)

func (prs *Parser) Errors() []string {
	return prs.errors
}

func (prs *Parser) peekError(tknType token.TokenType) {
	want, _ := token.LookupString(tknType)
	got, _ := token.LookupString(prs.peekToken.Type)
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", want, got)
	prs.errors = append(prs.errors, msg)
}

func (prs *Parser) typeError() {
	got, _ := token.LookupString(prs.curToken.Type)
	msg := fmt.Sprintf("expected type definition got: %s instead", got)
	prs.errors = append(prs.errors, msg)
}
