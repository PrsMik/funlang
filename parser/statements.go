package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseStatement() ast.StatementNode {
	switch prs.curToken.Type {
	case token.LET:
		return prs.parseLetStatement()
	case token.RETURN:
		return prs.parseReturnStatement()
	default:
		return nil
	}
}

func (prs *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: prs.curToken}

	if !prs.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}

	if !prs.expectPeek(token.COLON) {
		return nil
	}

	prs.nextToken()

	statement.Type = prs.parseType()
	if statement.Type == nil {
		prs.typeError()
		return nil
	}

	if !prs.expectPeek(token.ASSIGN) {
		return nil
	}

	prs.nextToken()

	statement.Value = prs.parseExpression(LOWEST)

	for !prs.curTokenIs(token.SEMICOLON) {
		prs.nextToken()
	}

	return statement
}

func (prs *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: prs.curToken}

	prs.nextToken()

	statement.Value = prs.parseExpression(int(token.IDENT))

	for !prs.curTokenIs(token.SEMICOLON) {
		prs.nextToken()
	}

	return statement
}
