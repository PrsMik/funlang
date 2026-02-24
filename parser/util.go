package parser

import "funlang/token"

func (prs *Parser) nextToken() {
	prs.curToken = prs.peekToken
	prs.peekToken = prs.lxr.NextToken()
}

func (prs *Parser) curTokenIs(tknType token.TokenType) bool {
	return prs.curToken.Type == tknType
}

func (prs *Parser) peekTokenIs(tknType token.TokenType) bool {
	return prs.peekToken.Type == tknType
}

func (prs *Parser) expectPeek(tknType token.TokenType) bool {
	if prs.peekTokenIs(tknType) {
		prs.nextToken()
		return true
	} else {
		prs.peekError(tknType)
		return false
	}
}
