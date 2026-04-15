package parser

import (
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

func (prs *Parser) parseGroupedExpression() ast.ExpressionNode {
	prs.nextToken()

	expr := prs.parseExpression(LOWEST)

	if !prs.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (prs *Parser) parseCallExpression(function ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.CallExpression{Token: prs.curToken, Function: function}
	expr.Arguments = prs.parseExpressionList(token.RPAREN)
	expr.SemiToken = prs.curToken
	return expr
}

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

func (prs *Parser) parseBoolean() ast.ExpressionNode {
	return &ast.BooleanLiteral{Token: prs.curToken, Value: prs.curTokenIs(token.TRUE)}
}

func (prs *Parser) parseArrayLiteral() ast.ExpressionNode {
	array := &ast.ArrayLiteral{Token: prs.curToken}
	array.Elements = prs.parseExpressionList(token.RBRACKET)
	array.SemiToken = prs.curToken
	return array
}

func (prs *Parser) parseHashMapLiteral() ast.ExpressionNode {
	hashMapLiteral := &ast.HashMapLiteral{Token: prs.curToken}
	hashMapLiteral.Pairs = make(map[ast.ExpressionNode]ast.ExpressionNode)

	for !prs.peekTokenIs(token.RBRACE) {
		reserveCurToken := prs.curToken

		prs.nextToken()
		key := prs.parseExpression(LOWEST)

		// if !prs.expectPeek(token.COLON) {
		// 	return nil
		// }
		prs.expectPeek(token.COLON)

		if key == nil {
			key = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		}

		reserveCurToken = prs.curToken

		prs.nextToken()

		value := prs.parseExpression(LOWEST)

		if value == nil {
			hashMapLiteral.Pairs[key] = &ast.BadExpression{From: reserveCurToken.Start, To: prs.curToken.Start}
		} else {
			hashMapLiteral.Pairs[key] = value
		}

		if !prs.peekTokenIs(token.RBRACE) && !prs.expectPeek(token.COMMA) {
			// return nil
			break
		}
	}

	// if !prs.expectPeek(token.RBRACE) {
	// 	return nil
	// }

	if prs.curTokenIs(token.COMMA) {
		resTailCurToken := prs.curToken

		prs.expectPeek(token.RBRACE)

		tailVirtualExpr := &ast.BadExpression{From: resTailCurToken.End, To: prs.curToken.Start}

		hashMapLiteral.Pairs[tailVirtualExpr] = &ast.BadExpression{
			From: token.Position{Line: -1, Column: -1},
			To:   token.Position{Line: -1, Column: -1},
		}
	} else {
		prs.expectPeek(token.RBRACE)
	}

	hashMapLiteral.SemiToken = prs.curToken

	return hashMapLiteral
}

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

	// if prs.curTokenIs(end) {
	// 	return list
	// }

	// if !prs.expectPeek(end) {
	// 	return nil
	// }
	prs.expectPeek(end)

	return list
}

func (prs *Parser) parseFunctionLiteral() ast.ExpressionNode {
	fnLiteral := &ast.FunctionLiteral{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	fnLiteral.Parameters, fnLiteral.ParamTypes = prs.parseFunctionParameters()

	if prs.peekTokenIs(token.RARROW) {
		prs.nextToken()
		prs.nextToken()
		fnLiteral.ReturnType = prs.parseType()
	}

	if !prs.expectPeek(token.LBRACE) {
		return nil
	}

	fnLiteral.Body = prs.parseBlockStatement()

	return fnLiteral
}

func (prs *Parser) parseFunctionParameters() ([]*ast.Identifier, []ast.TypeNode) {
	literals := []*ast.Identifier{}
	paramTypes := []ast.TypeNode{}

	if prs.peekTokenIs(token.RPAREN) {
		prs.nextToken()
		return literals, paramTypes
	}

	prs.nextToken()

	ident := &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}
	literals = append(literals, ident)

	if prs.peekTokenIs(token.COLON) {
		prs.nextToken()
		prs.nextToken()
		paramTypes = append(paramTypes, prs.parseType())
	} else {
		paramTypes = append(paramTypes, nil)
	}

	for prs.peekTokenIs(token.COMMA) {
		prs.nextToken()
		prs.nextToken()
		ident := &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}
		literals = append(literals, ident)
		if prs.peekTokenIs(token.COLON) {
			prs.nextToken()
			prs.nextToken()
			paramTypes = append(paramTypes, prs.parseType())
		} else {
			paramTypes = append(paramTypes, nil)
		}
	}

	if !prs.expectPeek(token.RPAREN) {
		return nil, nil
	}

	return literals, paramTypes
}
