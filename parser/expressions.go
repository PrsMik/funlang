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

func (prs *Parser) parseIfExpression() ast.ExpressionNode {
	expr := &ast.IfExpression{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	prs.nextToken()
	expr.Condition = prs.parseExpression(LOWEST)

	if !prs.expectPeek(token.RPAREN) {
		return nil
	}

	if !prs.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = prs.parseBlockStatement()

	if !prs.expectPeek(token.ELSE) {
		return nil
	}

	if !prs.expectPeek(token.LBRACE) {
		return nil
	}

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
	expr.Arguments = prs.parseFunctionArguments()
	return expr
}

func (prs *Parser) parseFunctionArguments() []ast.ExpressionNode {
	args := []ast.ExpressionNode{}

	if prs.peekTokenIs(token.RPAREN) {
		prs.nextToken()
		return args
	}

	prs.nextToken()
	args = append(args, prs.parseExpression(LOWEST))

	for prs.peekTokenIs(token.COMMA) {
		prs.nextToken()
		prs.nextToken()
		args = append(args, prs.parseExpression(LOWEST))
	}

	if !prs.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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

func (prs *Parser) parseBoolean() ast.ExpressionNode {
	return &ast.BooleanLiteral{Token: prs.curToken, Value: prs.curTokenIs(token.TRUE)}
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

func (prs *Parser) noPrefixParseFnError(tknType token.TokenType) {
	tknStr, _ := token.LookupString(tknType)
	msg := fmt.Sprintf("no prefix parse function for %s found", tknStr)
	prs.errors = append(prs.errors, msg)
}
