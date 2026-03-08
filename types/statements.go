package types

import (
	"fmt"
	"funlang/ast"
)

func (chk *TypeChecker) checkStatement(stmt ast.StatementNode) {
	switch curStmt := stmt.(type) {
	case *ast.LetStatement:
		chk.checkLetStatement(curStmt)
	case *ast.ReturnStatement:
	case *ast.BlockStatement:
	default:
		chk.errors = append(chk.errors, fmt.Errorf("unknown statement type"))
	}
}

func (chk *TypeChecker) checkLetStatement(stmt *ast.LetStatement) {
	expectedType := chk.resolveType(stmt.Type)

	actualType := chk.checkExpression(stmt.Value)

	if !Equals(expectedType, actualType) {
		if expectedType != nil && actualType != nil {
			chk.errors = append(chk.errors, fmt.Errorf("expected type %s, got %s", expectedType.Signature(), actualType.Signature()))
		}
	}

	chk.env.Set(stmt.Name.Value, expectedType)
}
