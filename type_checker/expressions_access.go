package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkMemberAccessExpression(expr ast.ExpressionNode) types.Type {
	chk.typeError(fmt.Sprintf("unknown expression type %T", expr), expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkTypeAccessExpression(expr ast.ExpressionNode) types.Type {
	chk.typeError(fmt.Sprintf("unknown expression type %T", expr), expr)
	return &types.IllegalType{}
}
