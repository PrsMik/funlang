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

	if prs.peekTokenIs(token.LPAREN) {
		prs.expectPeek(token.LPAREN)
		expr.TargetType = prs.parseType()
	} else {
		expr.TargetType = nil
	}

	prs.expectPeek(token.LBRACE)

	expr.Body = prs.parseBlockStatement()

	return expr
}
