package types

import (
	"fmt"
	"funlang/ast"
)

func (chk *TypeChecker) checkStatement(stmt ast.StatementNode) Type {
	switch curStmt := stmt.(type) {
	case *ast.LetStatement:
		return chk.checkLetStatement(curStmt)
	case *ast.ReturnStatement:
		return chk.checkReturnStatement(curStmt)
	case *ast.BlockStatement:
		return chk.checkBlockStatement(curStmt)
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown statement type"))
	return &IllegalType{}
}

func (chk *TypeChecker) checkLetStatement(stmt *ast.LetStatement) Type {
	expectedType := chk.resolveType(stmt.Type)

	actualType := chk.checkExpression(stmt.Value)

	if !Equals(expectedType, actualType) {
		if expectedType != nil && actualType != nil {
			chk.errors = append(chk.errors, fmt.Errorf("expected type %s, got %s", expectedType.Signature(), actualType.Signature()))
		}
	}

	chk.env.Set(stmt.Name.Value, expectedType)

	return actualType
}

func (chk *TypeChecker) checkReturnStatement(stmt *ast.ReturnStatement) Type {
	returnType := chk.checkExpression(stmt.Value)
	return returnType
}

func (chk *TypeChecker) checkBlockStatement(stmt *ast.BlockStatement) Type {
	chk.env = NewEnclosedTypeEviroment(chk.env)

	var returnType Type = &IllegalType{}

	for _, stmt := range stmt.Statements {
		switch stmt.(type) {
		case *ast.ReturnStatement:
			returnType = chk.checkStatement(stmt)
		}
	}

	chk.env = chk.env.outer
	return returnType
}
