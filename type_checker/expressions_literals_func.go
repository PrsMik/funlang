package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkCallExpression(expr ast.ExpressionNode) types.Type {
	callExpr := expr.(*ast.CallExpression)
	rawCallType := chk.checkExpression(callExpr.Function)

	if _, err := rawCallType.(*types.IllegalType); err {
		return &types.IllegalType{}
	}

	switch callFuncType := rawCallType.(type) {
	case *types.FuncType:
		if len(callExpr.Arguments) != len(callFuncType.Params) {
			if len(callFuncType.Params) != 0 {
				chk.typeError(fmt.Sprintf(`wrong number of arguments for "%s"`, callExpr.Function.String()), expr)
				return &types.IllegalType{}
			} else {
				if _, ok := callExpr.Arguments[0].(*ast.BadExpression); ok {
					callExpr.Arguments = callExpr.Arguments[1:]
				}
			}
		}

		for i, arg := range callExpr.Arguments {
			argType := chk.checkExpression(arg)

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
			argTypes[i] = chk.checkExpression(argExpr)
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

	// типы параметров описаны в самом литерале
	if len(funLiteral.ParamTypes) != 0 && funLiteral.ParamTypes[0] != nil || funLiteral.ReturnType != nil {

		if len(funLiteral.ParamTypes) != 0 && funLiteral.ParamTypes[0] != nil {
			// у литерала определены типы параметров
			for ind, param := range funLiteral.ParamTypes {
				paramType := chk.resolveType(param)
				funcParam := &types.FuncParam{Name: funLiteral.Parameters[ind].Value, Type: paramType}
				expectedFuncType.Params = append(expectedFuncType.Params, *funcParam)
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

	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.recordScope(expr, chk.env)

	// спуск типов из ожидаемого типа функции на переменные в литерале
	for i, param := range funLiteral.Parameters {
		parameter := types.FuncParam{Name: param.Value, Type: expectedFuncType.Params[i].Type}
		resFuncType.Params = append(resFuncType.Params, parameter)

		chk.recordType(param, expectedFuncType.Params[i].Type)
		chk.env.Set(param.Value, expectedFuncType.Params[i].Type, param)
	}

	oldExp := chk.curExpectedType

	chk.curExpectedType = expectedFuncType.ReturnType

	resFuncType.ReturnType = chk.checkBlockStatement(funLiteral.Body)

	if !types.IsAssignable(expectedFuncType.ReturnType, resFuncType.ReturnType) {
		chk.typeError(fmt.Sprintf("function literal has return type %s, but expected %s",
			resFuncType.ReturnType.Signature(), expectedFuncType.ReturnType.Signature()), expr)
		return &types.IllegalType{}
	}

	chk.curExpectedType = oldExp

	chk.env = chk.env.Outer

	return resFuncType
}
