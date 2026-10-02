package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseIfExpression() ast.ExpressionNode {
	expr := &ast.IfExpression{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	reserveCurToken := prs.curToken

	prs.nextToken()

	reserveEndToken := prs.curToken

	expr.Condition = prs.parseExpression(LOWEST)

	if expr.Condition == nil {
		expr.Condition = &ast.BadExpression{From: reserveCurToken.End, To: reserveEndToken.Start}
	}

	// if !prs.expectPeek(token.RPAREN) {
	// 	return nil
	// }

	if !prs.curTokenIs(token.RPAREN) {
		prs.expectPeek(token.RPAREN)
	}

	// if !prs.expectPeek(token.LBRACE) {
	// 	return nil
	// }

	if !prs.curTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)
	}

	expr.Consequence = prs.parseBlockStatement()

	// if !prs.expectPeek(token.ELSE) {
	// 	return nil
	// }

	prs.expectPeek(token.ELSE)

	// if !prs.expectPeek(token.LBRACE) {
	// 	return nil
	// }

	prs.expectPeek(token.LBRACE)

	expr.Alternative = prs.parseBlockStatement()

	return expr
}

func (prs *Parser) parseSwitchExpression() ast.ExpressionNode {
	expr := &ast.SwitchExpression{Token: prs.curToken}

	prs.nextToken()

	expr.Value = prs.parseExpression(LOWEST)

	if !prs.curTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)
	}

	for prs.peekTokenIs(token.CASE) {
		prs.expectPeek(token.CASE)

		caseBlock := &ast.CaseBlock{Token: prs.curToken}

		prs.nextToken()

		caseBlock.Pattern = prs.parseExpression(LOWEST)

		if !prs.curTokenIs(token.LBRACE) {
			prs.expectPeek(token.LBRACE)
		}

		caseBlock.Body = prs.parseBlockStatement()

		expr.Cases = append(expr.Cases, caseBlock)
	}

	prs.expectPeek(token.DEFAULT)

	prs.expectPeek(token.LBRACE)

	expr.Default = prs.parseBlockStatement()

	prs.expectPeek(token.RBRACE)

	return expr
}
