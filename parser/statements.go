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
		prs.stmtError()
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
		return nil
	}

	if !prs.expectPeek(token.ASSIGN) {
		return nil
	}

	prs.nextToken()

	statement.Value = prs.parseExpression(LOWEST)
	if statement.Value == nil {
		return nil
	}

	if !prs.expectPeek(token.SEMICOLON) {
		// prs.nextToken()
		return nil
	}

	return statement
}

func (prs *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: prs.curToken}

	prs.nextToken()

	statement.Value = prs.parseExpression(LOWEST)
	if statement.Value == nil {
		return nil
	}

	if prs.peekTokenIs(token.SEMICOLON) {
		prs.nextToken()
	}

	return statement
}

func (prs *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: prs.curToken}
	block.Statements = []ast.StatementNode{}

	prs.nextToken()

	for !prs.curTokenIs(token.RBRACE) && !prs.curTokenIs(token.EOF) {
		stmt := prs.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		prs.nextToken()
	}

	return block
}
