package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseType() ast.TypeNode {
	switch prs.curToken.Type {
	case token.INT_TYPE, token.BOOL_TYPE, token.STRING_TYPE:
		return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
	case token.LBRACKET:
		return prs.parseArrayType()
	case token.LBRACE:
		return prs.parseHashMapType()
	case token.FN:
		return prs.parseFunctionType()
	default:
		prs.typeError()
		return nil
	}
}

func (prs *Parser) parseArrayType() ast.TypeNode {
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

func (prs *Parser) parseHashMapType() ast.TypeNode {
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

func (prs *Parser) parseFunctionType() ast.TypeNode {
	fnType := &ast.FunctionType{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	fnType.ParamsTypes = prs.parseFunctionParamsTypes()

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

func (prs *Parser) parseFunctionParamsTypes() []ast.TypeNode {
	params := []ast.TypeNode{}

	if prs.peekTokenIs(token.RPAREN) {
		prs.nextToken()
		return params
	}

	prs.nextToken()

	params = append(params, prs.parseType())

	for prs.peekTokenIs(token.COMMA) {
		prs.nextToken()
		prs.nextToken()
		params = append(params, prs.parseType())
	}

	if !prs.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}
