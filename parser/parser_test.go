package parser

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
			let x: int = 5;
			let y: int = 10;
			let foobar: bool = true;
			`
	lxr := lexer.New(input)
	prs := New(lxr)
	program := prs.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       token.TokenType
	}{
		{"x", token.INT_TYPE},
		{"y", token.INT_TYPE},
		{"foobar", token.BOOL_TYPE},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
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
