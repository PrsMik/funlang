package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseIndexExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.IndexExpression{Token: prs.curToken, Left: left}

	prs.nextToken()

	reserveCurToken := prs.curToken

	expression.Index = prs.parseExpression(LOWEST)

	if expression.Index == nil {
		expression.Index = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.End}
	}

	// if !prs.expectPeek(token.RBRACKET) {
	// 	return nil
	// }
	prs.expectPeek(token.RBRACKET)

	expression.SemiToken = prs.curToken

	return expression
}

func (prs *Parser) parseDotOperatorExpression(left ast.ExpressionNode) ast.ExpressionNode {
	if prs.peekTokenIs(token.LPAREN) {
		return prs.parseTypeAccessExpression(left)
	} else if prs.peekTokenIs(token.LBRACE) {
		return prs.parseImplInstantiationExpression(left)
	} else {
		return prs.parseMemberAccessExpression(left)
	}
}

func (prs *Parser) parseTypeAccessExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.TypeAccessExpression{Token: prs.curToken, Left: left}

	prs.expectPeek(token.LPAREN)

	if !prs.expectPeek(token.TYPE) {
		return nil
	}

	prs.expectPeek(token.RPAREN)

	expr.SemiToken = prs.curToken

	return expr
}

func (prs *Parser) parseMemberAccessExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.MemberAccessExpression{Token: prs.curToken, Left: left}

	if !prs.expectPeek(token.IDENT) {
		return nil
	}
	expr.Property = &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}

	return expr
}

func (prs *Parser) parseImplInstantiationExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.ImplInstantiationExpression{Token: prs.curToken, Left: left}

	prs.nextToken()

	fields := prs.parseHashMapLiteral()

	if fields == nil {
		return nil
	}

	expr.Fields = fields

	return expr
}
