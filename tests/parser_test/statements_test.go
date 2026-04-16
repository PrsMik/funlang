package parser_test

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/parser"
	"funlang/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
			let x: int = 5; // пять
			let y: int = 10; // десять
			// абоба
			let foobar: bool = true;
			`
	lxr := lexer.New(input)
	prs := parser.New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
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

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return 1 + 1;
			return add(1, 1);
			`
	lxr := lexer.New(input)
	prs := parser.New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}
