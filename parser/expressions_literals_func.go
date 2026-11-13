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

func (prs *Parser) parseFunctionPrefix() ast.ExpressionNode {
	fnType := &ast.FunctionType{Token: prs.curToken}
	fnLiteral := &ast.FunctionLiteral{Token: prs.curToken}
	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	params, paramsTypes := prs.parseFunctionParameters()
	if params == nil && paramsTypes == nil {
		return nil
	}

	var retType ast.ExpressionNode
	if prs.peekTokenIs(token.RARROW) {
		prs.nextToken()
		prs.nextToken()
		retType = prs.parseType()
		if retType == nil {
			return nil
		}
	}

	if prs.typeDepth == 0 && prs.peekTokenIs(token.LBRACE) {
		prs.expectPeek(token.LBRACE)

		fnLiteral.Parameters = params
		fnLiteral.ParamsTypes = paramsTypes
		fnLiteral.ReturnType = retType
		fnLiteral.Body = prs.parseBlockStatement()
		return fnLiteral
	}

	if retType == nil {
		prs.peekError(token.RARROW)
		return nil
	}

	fnType.Parameters = params
	fnType.ParamsTypes = paramsTypes
	fnType.ReturnType = retType

	return fnType
}

func (prs *Parser) parseFunctionParameters() ([]*ast.Identifier, []ast.ExpressionNode) {
	params := []*ast.Identifier{}
	paramTypes := []ast.ExpressionNode{}

	if prs.peekTokenIs(token.RPAREN) {
		prs.nextToken()
		return params, paramTypes
	}

	prs.nextToken()

	for {
		name, pType := prs.parseFunctionParameter()
		if name == nil && pType == nil {
			return nil, nil
		}
		params = append(params, name)
		paramTypes = append(paramTypes, pType)

		if !prs.peekTokenIs(token.COMMA) {
			break
		}
		prs.nextToken()
		prs.nextToken()
	}

	if !prs.expectPeek(token.RPAREN) {
		return nil, nil
	}

	return params, paramTypes
}

func (prs *Parser) parseFunctionParameter() (*ast.Identifier, ast.ExpressionNode) {
	if prs.curTokenIs(token.IDENT) && prs.peekTokenIs(token.COLON) {
		param := &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}
		prs.nextToken()
		prs.nextToken()
		paramType := prs.parseType()
		if paramType == nil {
			return nil, nil
		}
		return param, paramType
	}

	if prs.typeDepth == 0 && prs.curTokenIs(token.IDENT) {
		return &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}, nil
	}

	paramType := prs.parseType()
	if paramType == nil {
		return nil, nil
	}
	return nil, paramType
}
