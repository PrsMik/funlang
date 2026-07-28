package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"
)

func (chk *TypeChecker) checkIfExpression(expr ast.ExpressionNode) types.Type {
	// --- РАЗМЕТКА УСЛОВИЯ ---
	ifExpr := expr.(*ast.IfExpression)

	startPos := ifExpr.Start()
	var endPos token.Position

	if ifExpr.Consequence != nil {
		endPos = ifExpr.Consequence.Start()
	} else {
		endPos = ifExpr.End()
	}

	condArea := &ast.BadExpression{
		From: startPos,
		To:   endPos,
	}

	chk.recordExpectedType(condArea, &types.BoolType{})

	// --- ПРОВЕРКА ТИПА ---
	oldType := chk.curExpectedType

	chk.curExpectedType = &types.BoolType{}

	condType := chk.checkExpression(expr.(*ast.IfExpression).Condition)

	chk.curExpectedType = oldType

	if !types.Equals(condType, &types.BoolType{}) {
		chk.typeError(fmt.Sprintf("wrong type %s for if condition", condType.Signature()), expr)
		return &types.IllegalType{}
	}

	conseqType := chk.checkBlockStatement(expr.(*ast.IfExpression).Consequence)

	if expr.(*ast.IfExpression).Alternative != nil {
		alterType := chk.checkBlockStatement(expr.(*ast.IfExpression).Alternative)

		if !types.Equals(conseqType, alterType) {
			chk.typeError(fmt.Sprintf("type mismatch between %s & %s in if/else branches",
				conseqType.Signature(), alterType.Signature()), expr)
			return &types.IllegalType{}
		}
	}

	return conseqType
}

func (chk *TypeChecker) checkSwitchExpression(expr ast.ExpressionNode) types.Type {
	chk.typeError(fmt.Sprintf("unknown expression type %T", expr), expr)
	return &types.IllegalType{}
}
