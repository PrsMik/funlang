package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseCallExpression(function ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.CallExpression{Token: prs.curToken, Function: function}
	expr.Arguments = prs.parseExpressionList(token.RPAREN)
	expr.SemiToken = prs.curToken
	return expr
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

func (prs *Parser) parseFunctionParameters() ([]*ast.Identifier, []ast.ExpressionNode) {
	literals := []*ast.Identifier{}
	paramTypes := []ast.ExpressionNode{}

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
