package parser

import (
	"funlang/ast"
	"funlang/token"
)

type (
	prefixParseFn func() ast.ExpressionNode
	infixParseFn  func(ast.ExpressionNode) ast.ExpressionNode
)

func (prs *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	prs.prefixParseFns[tokenType] = fn
}
func (prs *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	prs.infixParseFns[tokenType] = fn
}

func (prs *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefix := prs.prefixParseFns[prs.curToken.Type]
	var leftExp ast.ExpressionNode
	if prefix == nil {
		prs.noPrefixParseFnError(prs.curToken.Type)
		return nil
	}

	leftExp = prefix()

	for !prs.peekTokenIs(token.SEMICOLON) && precedence < prs.peekTokenPrecedence() {
		infix := prs.infixParseFns[prs.peekToken.Type]

		if infix == nil {
			return nil
		}

		prs.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (prs *Parser) parsePrefixExpression() ast.ExpressionNode {
	expression := &ast.PrefixExpression{Token: prs.curToken, Operator: prs.curToken.Literal}
	operatorEnd := prs.curToken.End

	prs.nextToken()

	expression.Right = prs.parseExpression(PREFIX)

	if expression.Right == nil {
		expression.Right = &ast.BadExpression{
			From: operatorEnd,
			To:   prs.curToken.Start,
		}
	}

	return expression
}

func (prs *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.InfixExpression{Token: prs.curToken, Operator: prs.curToken.Literal, Left: left}

	operatorEnd := prs.curToken.End
	precedence := prs.curTokenPrecedence()

	prs.nextToken()

	expression.Right = prs.parseExpression(precedence)

	if expression.Right == nil {
		expression.Right = &ast.BadExpression{
			From: operatorEnd,
			To:   prs.curToken.Start,
		}
	}

	return expression
}

func (prs *Parser) parseGroupedExpression() ast.ExpressionNode {
	prs.nextToken()

	expr := prs.parseExpression(LOWEST)

	if !prs.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (prs *Parser) parseComptimeExpression() ast.ExpressionNode {
	expr := &ast.ComptimeExpression{Token: prs.curToken}
	prs.nextToken()
	expr.Value = prs.parseExpression(LOWEST)
	return expr
}
