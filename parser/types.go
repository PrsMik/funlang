package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseType() ast.ExpressionNode {
	switch prs.curToken.Type {
	case token.TYPE:
		return &ast.TypeType{Token: prs.curToken, Value: prs.curToken.Literal}
	case token.INT_TYPE, token.BOOL_TYPE, token.STRING_TYPE:
		return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
	case token.LBRACKET:
		return prs.parseArrayType()
	case token.LBRACE:
		return prs.parseHashMapType()
	case token.FN:
		return prs.parseFunctionType()
	case token.INT, token.STRING, token.TRUE, token.FALSE, token.BANG, token.MINUS:
		prs.typeError()
		return nil
	default:
		return prs.parseExpression(LOWEST)
	}
}

func (prs *Parser) parseArrayType() ast.ExpressionNode {
	arrType := &ast.ArrayType{Token: prs.curToken}

	if prs.peekTokenIs(token.RBRACKET) {
		prs.nextToken()
		prs.typeError()
		return nil
	} else {
		prs.nextToken()
	}

	arrType.ElementsType = prs.parseType()

	prs.nextToken()

	return arrType
}

func (prs *Parser) parseHashMapType() ast.ExpressionNode {
	hashMapType := &ast.HashMapType{Token: prs.curToken}

	if prs.peekTokenIs(token.RBRACE) {
		prs.nextToken()
		prs.typeError()
		return nil
	} else {
		prs.nextToken()
	}

	hashMapType.KeyType = prs.parseType()

	if !prs.expectPeek(token.COLON) {
		return nil
	}

	prs.nextToken()

	hashMapType.ElementType = prs.parseType()

	prs.nextToken()

	return hashMapType
}

func (prs *Parser) parseFunctionType() ast.ExpressionNode {
	fnType := &ast.FunctionType{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	fnType.Parameters, fnType.ParamsTypes = prs.parseFunctionParamsTypes()

	if !prs.expectPeek(token.RARROW) {
		return nil
	}

	prs.nextToken()

	fnType.ReturnType = prs.parseType()
	if fnType.ReturnType == nil {
		return nil
	}

	return fnType
}

func (prs *Parser) parseFunctionParamsTypes() ([]*ast.Identifier, []ast.ExpressionNode) {
	params := []*ast.Identifier{}
	paramTypes := []ast.ExpressionNode{}

	if prs.peekTokenIs(token.RPAREN) {
		prs.nextToken()
		return params, paramTypes
	}

	prs.nextToken()

	if !prs.expectParseParamAndTypePair(&params, &paramTypes) {
		return nil, nil
	}

	for prs.peekTokenIs(token.COMMA) {
		prs.nextToken()
		prs.nextToken()
		if !prs.expectParseParamAndTypePair(&params, &paramTypes) {
			return nil, nil
		}
	}

	if !prs.expectPeek(token.RPAREN) {
		return nil, nil
	}

	return params, paramTypes
}

func (prs *Parser) expectParseParamAndTypePair(params *[]*ast.Identifier, paramTypes *[]ast.ExpressionNode) bool {
	firstLiteral := prs.parseExpression(LOWEST)

	if prs.peekTokenIs(token.COLON) {
		firstParam, ok := firstLiteral.(*ast.Identifier)
		if !ok {
			prs.typeError()
			return ok
		}
		*params = append(*params, firstParam)
		prs.nextToken()
		prs.nextToken()
		*paramTypes = append(*paramTypes, prs.parseType())
	} else {
		*paramTypes = append(*paramTypes, firstLiteral)
	}

	return true
}
