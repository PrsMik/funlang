package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseInterfaceLiteral() ast.ExpressionNode {
	expr := &ast.InterfaceLiteral{Token: prs.curToken}

	prs.expectPeek(token.LBRACE)

	expr.Body = prs.parseBlockStatement()

	return expr
}

func (prs *Parser) parseImplLiteral() ast.ExpressionNode {
	expr := &ast.ImplLiteral{Token: prs.curToken}

	if !prs.peekTokenIs(token.LBRACE) {
		prs.nextToken()
		expr.TargetType = prs.parseType(LOWEST)
	}

	if !prs.curTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)
	}

	expr.Body = prs.parseBlockStatement()

	return expr
}
