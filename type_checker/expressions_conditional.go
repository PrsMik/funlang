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
	switchExpr := expr.(*ast.SwitchExpression)

	valType := chk.checkExpression(switchExpr.Value)
	if _, isIllegal := valType.(*types.IllegalType); isIllegal {
		return valType
	}

	switchReturnType := chk.checkCaseBlocks(switchExpr, valType)

	// проверка default
	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.recordScope(switchExpr.Default, chk.env)

	defaultBodyType := chk.checkBlockStatement(switchExpr.Default)

	chk.env = chk.env.Outer

	if switchReturnType == nil {
		// если Illegal - пусто внутри блока
		if _, ok := defaultBodyType.(*types.IllegalType); !ok {
			switchReturnType = defaultBodyType
		}
	} else if !types.Equals(switchReturnType, defaultBodyType) {
		chk.typeError(fmt.Sprintf("switch default returns %s, but cases return %s",
			defaultBodyType.Signature(), switchReturnType.Signature()), switchExpr.Default)
	}

	if switchReturnType == nil {
		chk.typeError("switch must return some value in case blocks", switchExpr)
		return &types.IllegalType{}
	}

	return switchReturnType
}

func (chk *TypeChecker) checkCaseBlocks(switchExpr *ast.SwitchExpression, valType types.Type) types.Type {
	// проверки и изменения для type-switch
	isTypeSwitch := false
	var narrowedVar *ast.Identifier

	if _, ok := valType.(*types.TypeType); ok {
		isTypeSwitch = true

		if typeAccess, isAccess := switchExpr.Value.(*ast.TypeAccessExpression); isAccess {
			if ident, isIdent := typeAccess.Left.(*ast.Identifier); isIdent {
				narrowedVar = ident
			}
		}
	}

	var switchReturnType types.Type
	for _, caseBlock := range switchExpr.Cases {
		var patternType types.Type

		if isTypeSwitch {
			patternType = chk.resolveType(caseBlock.Pattern)
			if _, isIllegal := patternType.(*types.IllegalType); isIllegal {
				continue
			}
		} else {
			oldExp := chk.curExpectedType
			chk.curExpectedType = valType
			patternType = chk.checkExpression(caseBlock.Pattern)
			chk.curExpectedType = oldExp

			if !types.Equals(patternType, valType) {
				chk.typeError(fmt.Sprintf("case pattern type %s doesn't match switch value type %s",
					patternType.Signature(), valType.Signature()), caseBlock.Pattern)
			}
		}

		chk.env = types.NewEnclosedTypeEviroment(chk.env)
		chk.recordScope(caseBlock.Body, chk.env)

		// сужение типа для case
		if isTypeSwitch && narrowedVar != nil {
			chk.env.Set(narrowedVar.Value, patternType, narrowedVar)
		}

		caseBodyType := chk.checkBlockStatement(caseBlock.Body)

		chk.env = chk.env.Outer

		// ветки switch должны возвращать один и тот же тип
		if switchReturnType == nil {
			// если Illegal - пусто внутри блока
			if _, ok := caseBodyType.(*types.IllegalType); !ok {
				switchReturnType = caseBodyType
			}
		} else if !types.Equals(switchReturnType, caseBodyType) {
			chk.typeError(fmt.Sprintf("mismatch in switch cases return values: expected %s, got %s",
				switchReturnType.Signature(), caseBodyType.Signature()), caseBlock.Body)
		}
	}

	return switchReturnType
}
