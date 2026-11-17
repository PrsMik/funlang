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
		return prs.parseTypeCastAccessExpression(left)
	} else if prs.peekTokenIs(token.LBRACE) {
		return prs.parseImplInstantiationExpression(left)
	} else {
		return prs.parseMemberAccessExpression(left)
	}
}

func (prs *Parser) parseTypeCastAccessExpression(left ast.ExpressionNode) ast.ExpressionNode {
	dotToken := prs.curToken
	prs.expectPeek(token.LPAREN)

	if prs.peekTokenIs(token.TYPE) {
		expr := &ast.TypeAccessExpression{Token: prs.curToken, Left: left}
		prs.nextToken()
		prs.expectPeek(token.RPAREN)
		expr.SemiToken = prs.curToken
		return expr
	}

	expr := &ast.SafeTypeCastExpression{Token: dotToken, Left: left}
	prs.nextToken()

	expr.TargetType = prs.parseType(LOWEST)

	if expr.TargetType == nil {
		return nil
	}

	prs.expectPeek(token.RPAREN)

	expr.SemiToken = prs.curToken
	return expr
}

func (prs *Parser) parseUnsafeTypeCastExpression() ast.ExpressionNode {
	expr := &ast.UnsafeTypeCastExpression{Token: prs.curToken}

	prs.nextToken()

	parsedType := prs.parseType(LOWEST)

	if parsedType == nil {
		return nil
	}

	if callExpr, isCall := parsedType.(*ast.CallExpression); isCall {
		expr.TargetType = callExpr.Function
		if len(callExpr.Arguments) != 1 {
			prs.peekError(token.RPAREN)
			return nil
		}
		expr.Value = callExpr.Arguments[0]
		expr.SemiToken = callExpr.SemiToken
	} else {
		prs.peekError(token.LPAREN)
		return nil
	}
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

	fields := &ast.HashMapLiteral{Token: prs.curToken}

	list := prs.parsePairedExpressionList(token.COLON, token.RBRACE)

	if list == nil {
		return nil
	}

	fields.Pairs = list
	fields.SemiToken = prs.curToken

	expr.Fields = fields

	return expr
}
