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
	return nil
}

func (chk *TypeChecker) checkIntegerLiteral(expr ast.ExpressionNode) Type {
	return &IntType{}
}

func (chk *TypeChecker) checkBooleanLiteral(expr ast.ExpressionNode) Type {
	return &BoolType{}
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
	return nil
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
	case "==", "!=", "&&", "||":
		if Equals(leftType, &BoolType{}) && Equals(rightType, &BoolType{}) {
			return &BoolType{}
		}
	}
	chk.errors = append(chk.errors, fmt.Errorf("type mismatch between: %s %s; for operator %s",
		leftType.Signature(), rightType.Signature(), op))
	return nil
}
