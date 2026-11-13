package parser

import (
	"funlang/ast"
	"funlang/token"
	"maps"
)

type parsExpr = func(int) ast.ExpressionNode

func (prs *Parser) switchToTypePrefixes() {
	prs.registerPrefix(token.COMPTIME, prs.parseComptimeType)
}

func (prs *Parser) convertToTypeContext(fn parsExpr) parsExpr {
	return func(precedence int) ast.ExpressionNode {
		savedPrefixFns := prs.prefixParseFns

		typePrefixFns := make(map[token.TokenType]prefixParseFn, len(savedPrefixFns))
		maps.Copy(typePrefixFns, savedPrefixFns)

		prs.prefixParseFns = typePrefixFns
		prs.switchToTypePrefixes()

		defer func() { prs.prefixParseFns = savedPrefixFns }()

		res := fn(precedence)

		switch res.(type) {
		case *ast.IntegerLiteral,
			*ast.BooleanLiteral,
			*ast.StringLiteral,
			*ast.ArrayLiteral,
			*ast.HashMapLiteral,
			*ast.FunctionLiteral,
			*ast.PrefixExpression:

			prs.typeError()
			return nil
		default:
			return res
		}
	}
}

func (prs *Parser) parseType() ast.ExpressionNode {
	prs.typeDepth++
	defer func() { prs.typeDepth-- }()
	return prs.convertToTypeContext(prs.parseExpression)(LOWEST)
}

func (prs *Parser) parseTypeType() ast.ExpressionNode {
	return &ast.TypeType{Token: prs.curToken, Value: prs.curToken.Literal}
}

func (prs *Parser) parseComptimeType() ast.ExpressionNode {
	expr := &ast.ComptimeExpression{Token: prs.curToken}
	prs.nextToken()
	expr.Value = prs.parseType()
	if expr.Value == nil {
		return nil
	}
	return expr
}

func (prs *Parser) parseSimpleType() ast.ExpressionNode {
	return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
}
