package parser_test

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/parser"
	"testing"
)

func TestTypeParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, ast.TypeNode)
	}{
		{
			name:  "Simple Type Int",
			input: "let x: int = 1;",
			check: func(t *testing.T, node ast.TypeNode) {
				st, ok := node.(*ast.SimpleType)
				if !ok {
					t.Fatalf("expected *ast.SimpleType, got %T", node)
				}
				if st.Value != "int" {
					t.Errorf("expected 'int', got %s", st.Value)
				}
			},
		},
		{
			name:  "Array Type",
			input: "let x: [string] =[];",
			check: func(t *testing.T, node ast.TypeNode) {
				at, ok := node.(*ast.ArrayType)
				if !ok {
					t.Fatalf("expected *ast.ArrayType, got %T", node)
				}
				st, _ := at.ElementsType.(*ast.SimpleType)
				if st.Value != "string" {
					t.Errorf("expected array elements to be 'string'")
				}
			},
		},
		{
			name:  "HashMap Type",
			input: "let x: {string: int} = {};",
			check: func(t *testing.T, node ast.TypeNode) {
				ht, ok := node.(*ast.HashMapType)
				if !ok {
					t.Fatalf("expected *ast.HashMapType, got %T", node)
				}
				kt, _ := ht.KeyType.(*ast.SimpleType)
				vt, _ := ht.ElementType.(*ast.SimpleType)
				if kt.Value != "string" || vt.Value != "int" {
					t.Errorf("expected key 'string' and value 'int'")
				}
			},
		},
		{
			name:  "Function Type",
			input: "let x: fn(int, bool) -> string = fn(a, b) { return \"\"; };",
			check: func(t *testing.T, node ast.TypeNode) {
				ft, ok := node.(*ast.FunctionType)
				if !ok {
					t.Fatalf("expected *ast.FunctionType, got %T", node)
				}
				if len(ft.ParamsTypes) != 2 {
					t.Fatalf("expected 2 params")
				}
				pt1, _ := ft.ParamsTypes[0].(*ast.SimpleType)
				pt2, _ := ft.ParamsTypes[1].(*ast.SimpleType)
				rt, _ := ft.ReturnType.(*ast.SimpleType)
				if pt1.Value != "int" || pt2.Value != "bool" || rt.Value != "string" {
					t.Errorf("expected fn(int, bool) -> string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.LetStatement)
			tt.check(t, stmt.Type)
		})
	}
}
