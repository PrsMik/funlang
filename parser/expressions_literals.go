package parser

import (
	"funlang/ast"
	"funlang/token"
	"strconv"
)

func (prs *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}
}

func (prs *Parser) parseIntegerLiteral() ast.ExpressionNode {
	intLiteral := &ast.IntegerLiteral{Token: prs.curToken}

	value, err := strconv.ParseInt(prs.curToken.Literal, 0, 0)
	if err != nil {
		prs.integerLiteralParseError()
		return nil
	}

	intLiteral.Value = int(value)
	return intLiteral
}

func (prs *Parser) parseStringLiteral() ast.ExpressionNode {
	return &ast.StringLiteral{Token: prs.curToken, Value: prs.curToken.Literal[1 : len(prs.curToken.Literal)-1]}
}

func (prs *Parser) parseBooleanLiteral() ast.ExpressionNode {
	return &ast.BooleanLiteral{Token: prs.curToken, Value: prs.curTokenIs(token.TRUE)}
}

func (prs *Parser) parseArrayPrefix() ast.ExpressionNode {
	array := &ast.ArrayLiteral{Token: prs.curToken}
	// lbracket := prs.curToken
	// if !prs.expectPeek(token.RBRACKET) {
	// 	return nil
	// }

	// prs.nextToken()

	// elemType := prs.parseType()
	// if elemType == nil {
	// 	return nil
	// }

	arrType, ok := prs.parseType().(*ast.ArrayType)
	if !ok {
		return nil
	}

	if prs.peekTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)
		array.Elements = prs.parseExpressionList(token.RBRACE)
		array.TypeExpr = arrType.ElementsType
		array.SemiToken = prs.curToken
		return array
	}

	// return &ast.ArrayType{Token: lbracket, ElementsType: elemType}
	return arrType
}

func (prs *Parser) parseExpressionList(end token.TokenType) []ast.ExpressionNode {
	list := []ast.ExpressionNode{}

	if prs.peekTokenIs(end) {
		reserveCurToken := prs.curToken
		prs.nextToken()
		expr := &ast.BadExpression{From: reserveCurToken.End, To: prs.curToken.Start}
		list = append(list, expr)
		return list
	}

	prs.nextToken()

	expr := prs.parseExpression(LOWEST)
	list = append(list, expr)

	for prs.peekTokenIs(token.COMMA) || prs.curTokenIs(token.COMMA) {

		if !prs.curTokenIs(token.COMMA) && prs.peekTokenIs(token.COMMA) {
			prs.nextToken()
		}

		// prs.nextToken()
		reservePrevToken := prs.curToken

		prs.nextToken()
		reserveCurToken := prs.curToken

		expr = prs.parseExpression(LOWEST)
		if expr == nil {
			expr = &ast.BadExpression{From: reservePrevToken.End, To: reserveCurToken.Start}
		}

		list = append(list, expr)
	}

	prs.expectPeek(end)

	return list
}

func (prs *Parser) parseHashMapPrefix() ast.ExpressionNode {
	hashMapLiteral := &ast.HashMapLiteral{Token: prs.curToken}
	hashMapLiteral.Pairs = make(map[ast.ExpressionNode]ast.ExpressionNode)

	hashMapType, ok := prs.parseType().(*ast.HashMapType)
	if !ok {
		return nil
	}

	if prs.peekTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)

		// for !prs.peekTokenIs(token.RBRACE) {
		// 	reserveCurToken := prs.curToken

		// 	prs.nextToken()
		// 	key := prs.parseExpression(LOWEST)

		// 	prs.expectPeek(token.COLON)

		// 	if key == nil {
		// 		key = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		// 	}

		// 	reserveCurToken = prs.curToken

		// 	prs.nextToken()

		// 	value := prs.parseExpression(LOWEST)

		// 	if value == nil {
		// 		hashMapLiteral.Pairs[key] = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		// 	} else {
		// 		hashMapLiteral.Pairs[key] = value
		// 	}

		// 	if !prs.peekTokenIs(token.RBRACE) && !prs.expectPeek(token.COMMA) {
		// 		break
		// 	}
		// }

		// if prs.curTokenIs(token.COMMA) {
		// 	resTailCurToken := prs.curToken

		// 	prs.expectPeek(token.RBRACE)

		// 	tailVirtualExpr := &ast.BadExpression{From: resTailCurToken.End, To: prs.curToken.Start}

		// 	hashMapLiteral.Pairs[tailVirtualExpr] = &ast.BadExpression{
		// 		From: token.Position{Line: -1, Column: -1},
		// 		To:   token.Position{Line: -1, Column: -1},
		// 	}
		// } else {
		// 	prs.expectPeek(token.RBRACE)
		// }

		list := prs.parsePairedExpressionList(token.COLON, token.RBRACE)
		hashMapLiteral.Pairs = list

		hashMapLiteral.SemiToken = prs.curToken
		hashMapLiteral.TypeExpr = hashMapType

		return hashMapLiteral
	}

	return hashMapType
}

// предполагает, что находится на открывающем токене
// начала списка, аналогичном закрывающему
// пример: если список {a : b, c : d}
// начальный токен - первая {
func (prs *Parser) parsePairedExpressionList(
	delimeter,
	end token.TokenType) map[ast.ExpressionNode]ast.ExpressionNode {

	res := make(map[ast.ExpressionNode]ast.ExpressionNode)

	for !prs.peekTokenIs(end) {
		reserveCurToken := prs.curToken

		prs.nextToken()
		key := prs.parseExpression(LOWEST)

		prs.expectPeek(delimeter)

		if key == nil {
			key = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		}

		reserveCurToken = prs.curToken

		prs.nextToken()

		value := prs.parseExpression(LOWEST)

		if value == nil {
			res[key] = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		} else {
			res[key] = value
		}

		if !prs.peekTokenIs(end) && !prs.expectPeek(token.COMMA) {
			break
		}
	}

	if prs.curTokenIs(token.COMMA) {
		resTailCurToken := prs.curToken

		prs.expectPeek(end)

		tailVirtualExpr := &ast.BadExpression{From: resTailCurToken.End, To: prs.curToken.Start}

		res[tailVirtualExpr] = &ast.BadExpression{
			From: token.Position{Line: -1, Column: -1},
			To:   token.Position{Line: -1, Column: -1},
		}
	} else {
		prs.expectPeek(end)
	}

	return res
}
