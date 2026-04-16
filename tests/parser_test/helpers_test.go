package parser_test

import (
	"fmt"
	"funlang/ast"
	"funlang/parser"
	"funlang/token"
	"testing"
)

func checkParserErrors(t *testing.T, prs *parser.Parser) {
	errors := prs.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err.Msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statementNode ast.StatementNode, name string, tknType token.TokenType) bool {
	if statementNode.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statementNode.TokenLiteral())
		return false
	}

	letStmt, ok := statementNode.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statementNode)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	simpleType, ok := letStmt.Type.(*ast.SimpleType)
	if !ok {
		t.Errorf("letStmt.Type not *ast.SimpleType. got=%T", letStmt.Type)
		return false
	}

	if simpleType.Token.Type != tknType {
		want, _ := token.LookupString(tknType)
		got, _ := token.LookupString(simpleType.Token.Type)
		t.Errorf("type not correct. expected=%s, got=%s", want, got)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.ExpressionNode, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.ExpressionNode, expected interface{}) bool {
	switch val := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, val)
	case bool:
		return testBooleanLiteral(t, exp, val)
	case string:
		return testIdentifierLiteral(t, exp, val)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.ExpressionNode, value int) bool {
	integ, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integerLiteral not *ast.IntegerLiteral. got=%T", integerLiteral)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.ExpressionNode, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testIdentifierLiteral(t *testing.T, expression ast.ExpressionNode, value string) bool {
	ident, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not ast.Identifier; got: %T", expression)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value is not %s; got :%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is not %s; got :%s", value, ident.TokenLiteral())
		return false
	}
	return true
}
