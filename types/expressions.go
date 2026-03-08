package types

import (
	"fmt"
	"funlang/ast"
	"reflect"
)

func (chk *TypeChecker) checkExpression(expr ast.ExpressionNode) Type {
	exprTp := reflect.TypeOf(expr)
	if fn, ok := chk.expressionCheckFns[exprTp]; ok {
		return fn(expr)
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown expression type"))
	return &IllegalType{}
}

func (chk *TypeChecker) checkIntegerLiteral(expr ast.ExpressionNode) Type {
	return &IntType{}
}

func (chk *TypeChecker) checkBooleanLiteral(expr ast.ExpressionNode) Type {
	return &BoolType{}
}

func (chk *TypeChecker) checkIdentifier(expr ast.ExpressionNode) Type {
	identType, ok := chk.env.Get(expr.(*ast.Identifier).Value)
	if !ok {
		chk.errors = append(chk.errors, fmt.Errorf("unknown identifier: %s", expr.(*ast.Identifier).Value))
		return &IllegalType{}
	}
	return identType
}

func (chk *TypeChecker) checkPrefixExpression(expr ast.ExpressionNode) Type {
	op := expr.(*ast.PrefixExpression).Operator
	switch op {
	case "-":
		rightType := chk.checkExpression(expr.(*ast.PrefixExpression).Right)
		if Equals(rightType, &IntType{}) {
			return &IntType{}
		}
	case "!":
		rightType := chk.checkExpression(expr.(*ast.PrefixExpression).Right)
		if Equals(rightType, &BoolType{}) {
			return &BoolType{}
		}
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown operator: %s", op))
	return &IllegalType{}
}

func (chk *TypeChecker) checkInfixExpression(expr ast.ExpressionNode) Type {
	leftType := chk.checkExpression(expr.(*ast.InfixExpression).Left)
	rightType := chk.checkExpression(expr.(*ast.InfixExpression).Right)
	op := expr.(*ast.InfixExpression).Operator

	switch op {
	case "-", "+", "*", "/":
		if Equals(leftType, &IntType{}) && Equals(rightType, &IntType{}) {
			return &IntType{}
		}
	case "&&", "||":
		if Equals(leftType, &BoolType{}) && Equals(rightType, &BoolType{}) {
			return &BoolType{}
		}
	case "==", "!=":
		if (Equals(leftType, &BoolType{}) && Equals(rightType, &BoolType{})) ||
			(Equals(leftType, &IntType{}) && Equals(rightType, &IntType{})) {
			return &BoolType{}
		}
	case ">", "<", ">=", "<=":
		if Equals(leftType, &IntType{}) && Equals(rightType, &IntType{}) {
			return &BoolType{}
		}
	}

	chk.errors = append(chk.errors, fmt.Errorf("type mismatch between: %s %s; for operator %s",
		leftType.Signature(), rightType.Signature(), op))
	return &IllegalType{}
}

func (chk *TypeChecker) checkIfExpression(expr ast.ExpressionNode) Type {
	condType := chk.checkExpression(expr.(*ast.IfExpression).Condition)

	if !Equals(condType, &BoolType{}) {
		chk.errors = append(chk.errors, fmt.Errorf("wrong type %s for if condition", condType.Signature()))
		return &IllegalType{}
	}

	conseqType := chk.checkBlockStatement(expr.(*ast.IfExpression).Consequence)

	if expr.(*ast.IfExpression).Alternative != nil {
		alterType := chk.checkBlockStatement(expr.(*ast.IfExpression).Alternative)

		if !Equals(conseqType, alterType) {
			chk.errors = append(chk.errors, fmt.Errorf("type mismatch between %s & %s in if/else branches",
				conseqType.Signature(), alterType.Signature()))
			return &IllegalType{}
		}
	}

	return conseqType
}

func (chk *TypeChecker) checkFunctionLiteral(expr ast.ExpressionNode) Type {
	resFuncType := &FuncType{}

	expectedFuncType := chk.curExpectedType.(*FuncType)

	if len(expectedFuncType.ParamsTypes) != len(expr.(*ast.FunctionLiteral).Parameters) {
		chk.errors = append(chk.errors, fmt.Errorf("function literal has %d parameters, but expected %d",
			len(expr.(*ast.FunctionLiteral).Parameters), len(expectedFuncType.ParamsTypes)))
		return &IllegalType{}
	}

	chk.env = NewEnclosedTypeEviroment(chk.env)

	for i, param := range expr.(*ast.FunctionLiteral).Parameters {
		resFuncType.ParamsTypes = append(resFuncType.ParamsTypes, expectedFuncType.ParamsTypes[i])
		chk.env.Set(param.Value, expectedFuncType.ParamsTypes[i])
	}

	resFuncType.ReturnType = chk.checkBlockStatement(expr.(*ast.FunctionLiteral).Body)

	chk.env = chk.env.outer

	return resFuncType
}
