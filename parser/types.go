package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseType() ast.TypeNode {
	switch prs.curToken.Type {
	case token.INT_TYPE, token.BOOL_TYPE:
		return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
	case token.FN:
		return prs.parseFunctionType()
	default:
		prs.typeError()
		return nil
	}
}

func (prs *Parser) parseFunctionType() ast.TypeNode {
	fnType := &ast.FunctionType{Token: prs.curToken}

	if !prs.expectPeek(token.LPAREN) {
		return nil
	}

	fnType.ParamsTypes = prs.parseFunctionParamsTypes()

	return fnType
}

func (prs *Parser) parseFunctionParamsTypes() []ast.TypeNode {
	return []ast.TypeNode{}
}
