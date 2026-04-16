package parser_test

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/parser"
	"funlang/token"
	"testing"
)

func TestStatementBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedStart token.Position
		expectedEnd   token.Position
	}{
		{"Simple Let", "let x: int = 5;", token.Position{Line: 0, Column: 0}, token.Position{Line: 0, Column: 15}},
		{"Simple Return", "return 100;", token.Position{Line: 0, Column: 0}, token.Position{Line: 0, Column: 11}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0]
			if stmt.Start() != tt.expectedStart {
				t.Errorf("Statement Start wrong. expected=%+v, got=%+v", tt.expectedStart, stmt.Start())
			}
			if stmt.End() != tt.expectedEnd {
				t.Errorf("Statement End wrong. expected=%+v, got=%+v", tt.expectedEnd, stmt.End())
			}
		})
	}
}

func TestASTExpressionBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedStart token.Position
		expectedEnd   token.Position
	}{
		{
			"Int literal: 500",
			"let x: int = 500;",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 16},
		},
		{
			"Bool literal: true",
			"let x: int = true;",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 17},
		},
		{
			`String literal: "true"`,
			`let x: int = "true";`,
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 19},
		},
		{
			"Array literal: [1, 2, 3]",
			`let x: int = [1, 2, 3];`,
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 22},
		},
		{
			`Hash map literal: {"a" : 1, "b" : 2, "c" : 3}`,
			`let x: int = {"a" : 1, "b" : 2, "c" : 3};`,
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 40},
		},
		{
			"Function literal: fn(x) { return 1; }",
			`let x: int = fn(x) { return 1; };`,
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 32},
		},
		{
			"Ident literal: aboba",
			`let x: int = aboba;`,
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 18},
		},
		{
			"Call Expression: add(1, 2)",
			"let x: int = add(1, 2);",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 22},
		},
		{
			"Infix Expression: -1",
			"let x: int = -1;",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 15},
		},
		{
			"Infix Expression: 5 + 10",
			"let x: int = 5 + 10;",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 19},
		},
		{
			"If Expression: if (x) { y } else { z }",
			"let x: int = if (x) { return y; } else { return z; };",
			token.Position{Line: 0, Column: 13},
			token.Position{Line: 0, Column: 52},
		},
		{
			"Index Expression: [1, 2, 3][1 + 0]",
			`let x: int = [1, 2, 3][1 + 0];`,
			token.Position{Line: 0, Column: 22},
			token.Position{Line: 0, Column: 29},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0]
			var node ast.Node
			switch s := stmt.(type) {
			case *ast.LetStatement:
				node = s.Value
			case *ast.ReturnStatement:
				node = s.Value
			default:
				node = s
			}

			if node.Start() != tt.expectedStart {
				t.Errorf("Start wrong. expected=%+v, got=%+v", tt.expectedStart, node.Start())
			}
			if node.End() != tt.expectedEnd {
				t.Errorf("End wrong. expected=%+v, got=%+v", tt.expectedEnd, node.End())
			}
		})
	}
}
