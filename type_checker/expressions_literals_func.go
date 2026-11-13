package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func (chk *TypeChecker) checkCallExpression(expr ast.ExpressionNode) types.Type {
	callExpr := expr.(*ast.CallExpression)
	rawCallType := chk.checkExpression(callExpr.Function)

	if _, err := rawCallType.(*types.IllegalType); err {
		return &types.IllegalType{}
	}

	// предотвращает баг с одним аргументом func = fn(int) -> int ... x = func(),
	// нет ошибки; проблема была в проверке BadExpression
	if len(callExpr.Arguments) == 1 {
		if _, isBad := callExpr.Arguments[0].(*ast.BadExpression); isBad {
			callExpr.Arguments = callExpr.Arguments[:0]
		}
	}

	switch callFuncType := rawCallType.(type) {
	case *types.FuncType:
		if len(callExpr.Arguments) != len(callFuncType.Params) {
			chk.typeError(fmt.Sprintf(`wrong number of arguments for "%s"`, callExpr.Function.String()), expr)
			return &types.IllegalType{}
		}

		for i, arg := range callExpr.Arguments {
			outerExpectedType := chk.curExpectedType
			chk.curExpectedType = callFuncType.Params[i].Type
			argType := chk.checkExpression(arg)
			chk.curExpectedType = outerExpectedType

			if !types.IsAssignable(callFuncType.Params[i].Type, argType) {
				chk.typeError(fmt.Sprintf("wrong type for argument %d: %s in func call \" %s \"; expected %s",
					i+1, argType.Signature(), callExpr.Function.String(), callFuncType.Params[i].Type.Signature()), expr)
				return &types.IllegalType{}
			}
		}

		return callFuncType.ReturnType
	case *types.BuiltinFunc:
		argTypes := make([]types.Type, len(callExpr.Arguments))
		for i, argExpr := range callExpr.Arguments {
			outerExpectedType := chk.curExpectedType
			chk.curExpectedType = nil
			argTypes[i] = chk.checkExpression(argExpr)
			chk.curExpectedType = outerExpectedType
		}

		returnType, err := callFuncType.CheckFunc(argTypes)

		if err != nil {
			// chk.errors = append(chk.errors, err)
			chk.typeError(err.Error(), expr)
			return &types.IllegalType{}
		}

		return returnType
	default:
		chk.typeError(fmt.Sprintf(`identifier "%s" is no a funciton`, callExpr.Function.String()), expr)
		return &types.IllegalType{}
	}

}

func (chk *TypeChecker) checkFunctionLiteral(expr ast.ExpressionNode) types.Type {
	resFuncType := &types.FuncType{}
	funLiteral := expr.(*ast.FunctionLiteral)

	var expectedFuncType types.FuncType

	outerEnv := chk.env
	outerComptimeEnv := chk.comptimeEnv
	outerExpectedType := chk.curExpectedType

	chk.env = types.NewEnclosedTypeEviroment(outerEnv)
	chk.recordScope(expr, chk.env)
	chk.comptimeEnv = object.NewEnclosedEnvironment(outerComptimeEnv)

	// откат окружений в конце
	defer func() {
		chk.env = outerEnv
		chk.comptimeEnv = outerComptimeEnv
		chk.curExpectedType = outerExpectedType
	}()

	// типы параметров описаны в самом литерале
	if len(funLiteral.ParamsTypes) != 0 && funLiteral.ParamsTypes[0] != nil || funLiteral.ReturnType != nil {

		if len(funLiteral.ParamsTypes) != 0 && funLiteral.ParamsTypes[0] != nil {
			// у литерала определены типы параметров
			for ind, paramTypeLit := range funLiteral.ParamsTypes {
				paramType := chk.resolveType(paramTypeLit)
				param := funLiteral.Parameters[ind]

				funcParam := types.FuncParam{Name: param.Value, Type: paramType}
				funcParam.TypeParam = chk.bindParam(param, paramType, nil)
				expectedFuncType.Params = append(expectedFuncType.Params, funcParam)
			}

		} else if expType, ok := chk.curExpectedType.(*types.FuncType); ok && len(expType.Params) != 0 {
			// у ожидаемого типа определены типы параметров
			expectedFuncType.Params = expType.Params

		} else if len(funLiteral.Parameters) != 0 {
			chk.typeError("function literal has parameters, but no type specified for them", expr)
			return &types.IllegalType{}
		}

		if funLiteral.ReturnType != nil {
			expectedFuncType.ReturnType = chk.resolveType(funLiteral.ReturnType)
		} else {
			chk.typeError("function literal has typed parameters, but no return type specified", expr)
			return &types.IllegalType{}
		}

	} else {
		tempType, ok := chk.curExpectedType.(*types.FuncType)
		if !ok {
			chk.typeError("function literal has no specified types, but no func type is expected", expr)
			return &types.IllegalType{}
		}
		expectedFuncType = *tempType
	}

	if len(expectedFuncType.Params) != len(funLiteral.Parameters) {
		chk.typeError(fmt.Sprintf("function literal has %d parameters, but expected %d",
			len(funLiteral.Parameters), len(expectedFuncType.Params)), expr)
		return &types.IllegalType{}
	}

	// спуск типов из ожидаемого типа функции на переменные в литерале
	for i, param := range funLiteral.Parameters {
		paramType := expectedFuncType.Params[i].Type
		tparam := chk.bindParam(param, paramType, expectedFuncType.Params[i].TypeParam)

		resFuncType.Params = append(resFuncType.Params,
			types.FuncParam{Name: param.Value, Type: paramType, TypeParam: tparam})

		chk.recordType(param, paramType)
	}

	chk.curExpectedType = expectedFuncType.ReturnType

	resFuncType.ReturnType = chk.checkBlockStatement(funLiteral.Body)

	if !types.IsAssignable(expectedFuncType.ReturnType, resFuncType.ReturnType) {
		chk.typeError(fmt.Sprintf("function literal has return type %s, but expected %s",
			resFuncType.ReturnType.Signature(), expectedFuncType.ReturnType.Signature()), expr)
		return &types.IllegalType{}
	}

	return resFuncType
}
