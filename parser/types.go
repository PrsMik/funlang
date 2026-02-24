package parser

import (
	"funlang/ast"
	"funlang/token"
)

func (prs *Parser) parseType() ast.TypeNode {
	switch prs.curToken.Type {
	case token.INT_TYPE, token.BOOL_TYPE:
		return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
		// TODO: fn type
	}
	return nil
}
