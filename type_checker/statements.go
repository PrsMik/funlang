package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkStatement(stmt ast.StatementNode) types.Type {
	switch curStmt := stmt.(type) {
	case *ast.LetStatement:
		return chk.checkLetStatement(curStmt)
	case *ast.ReturnStatement:
		return chk.checkReturnStatement(curStmt)
	case *ast.BlockStatement:
		return chk.checkBlockStatement(curStmt)
	}
	chk.typeError("unknown statement type", stmt)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkLetStatement(stmt *ast.LetStatement) types.Type {
	if stmt == nil {
		return &types.IllegalType{}
	}

	if stmt.Type == nil {
		chk.typeError("missing type in let statement", stmt)
		return &types.IllegalType{}
	}

	expectedType := chk.resolveType(stmt.Type)

	chk.curExpectedType = expectedType

	chk.TypesInfo[stmt.Name] = expectedType
	chk.ExpectedTypes[stmt] = expectedType

	_, isFuncLit := stmt.Value.(*ast.FunctionLiteral)
	if isFuncLit {
		chk.env.Set(stmt.Name.Value, expectedType, stmt.Name)
	}

	actualType := chk.checkExpression(stmt.Value)

	if !types.Equals(expectedType, actualType) {
		if expectedType != nil && actualType != nil {
			if len(chk.errors) == 0 {
				chk.typeError(fmt.Sprintf("expected type %s, got %s", expectedType.Signature(), actualType.Signature()), stmt)
			}
		}
	}

	if !isFuncLit {
		chk.env.Set(stmt.Name.Value, expectedType, stmt.Name)
	}

	chk.curExpectedType = nil

	return actualType
}

func (chk *TypeChecker) checkReturnStatement(stmt *ast.ReturnStatement) types.Type {
	if stmt == nil {
		return &types.IllegalType{}
	}

	if stmt.Value == nil {
		chk.typeError("missing val in return statement", stmt)
		return &types.IllegalType{}
	}

	if chk.curExpectedType != nil {
		chk.ExpectedTypes[stmt] = chk.curExpectedType
	}

	returnType := chk.checkExpression(stmt.Value)
	return returnType
}

func (chk *TypeChecker) checkBlockStatement(stmt *ast.BlockStatement) types.Type {
	if stmt == nil {
		return &types.IllegalType{}
	}

	if stmt.Statements == nil {
		chk.typeError("missing body in block statement", stmt)
		return &types.IllegalType{}
	}

	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.Scopes[stmt] = chk.env

	var returnType types.Type = &types.IllegalType{}

	for _, stmt := range stmt.Statements {
		switch stmt.(type) {
		case *ast.ReturnStatement:
			returnType = chk.checkStatement(stmt)
		default:
			chk.checkStatement(stmt)
		}
	}

	chk.env = chk.env.Outer

	return returnType
}
