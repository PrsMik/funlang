package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseStatement() ast.StatementNode {
	switch prs.curToken.Type {
	case token.LET:
		stmt := prs.parseLetStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	case token.RETURN:
		stmt := prs.parseReturnStatement()
		if stmt == nil {
			return nil
		}
		return stmt
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

	reserveCurToken := prs.curToken

	prs.nextToken()

	statement.Value = prs.parseExpression(LOWEST)
	if statement.Value == nil {
		// return nil
		statement.Value = &ast.BadExpression{From: reserveCurToken.End, To: prs.curToken.Start}
	}

	reserveCurToken = prs.curToken

	if prs.peekTokenIs(token.SEMICOLON) {
		// return nil
		prs.nextToken()
	} else {
		if !prs.curTokenIs(token.SEMICOLON) {
			prs.tokenError(prs.curToken, token.Token{Type: token.SEMICOLON}, reserveCurToken)
		}
	}

	statement.SemiToken = prs.curToken

	return statement
}

func (prs *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: prs.curToken}

	reserveCurToken := prs.curToken

	prs.nextToken()

	statement.Value = prs.parseExpression(LOWEST)
	if statement.Value == nil {
		// return nil
		statement.Value = &ast.BadExpression{From: reserveCurToken.End, To: prs.curToken.Start}
	}

	reserveCurToken = prs.curToken

	if prs.peekTokenIs(token.SEMICOLON) {
		// return nil
		prs.nextToken()
	} else {
		if !prs.curTokenIs(token.SEMICOLON) {
			prs.tokenError(prs.curToken, token.Token{Type: token.SEMICOLON}, reserveCurToken)
		}
	}

	statement.SemiToken = prs.curToken

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

	block.SemiToken = prs.curToken

	return block
}
