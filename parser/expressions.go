package parser

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"strconv"
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
	if prefix == nil {
		prs.noPrefixParseFnError(prs.curToken.Type)
		return nil
	}
	leftExp := prefix()

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

	prs.nextToken()

	expression.Right = prs.parseExpression(PREFIX)

	return expression
}

func (prs *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.InfixExpression{Token: prs.curToken, Operator: prs.curToken.Literal, Left: left}

	precedence := prs.curTokenPrecedence()
	prs.nextToken()
	expression.Right = prs.parseExpression(precedence)

	return expression
}

func (prs *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}
}

func (prs *Parser) parseIntegerLiteral() ast.ExpressionNode {
	intLiteral := &ast.IntegerLiteral{Token: prs.curToken}

	value, err := strconv.ParseInt(prs.curToken.Literal, 0, 0)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", prs.curToken.Literal)
		prs.errors = append(prs.errors, msg)
		return nil
	}

	intLiteral.Value = int(value)
	return intLiteral
}

func (prs *Parser) noPrefixParseFnError(tknType token.TokenType) {
	tknStr, _ := token.LookupString(tknType)
	msg := fmt.Sprintf("no prefix parse function for %s found", tknStr)
	prs.errors = append(prs.errors, msg)
}
