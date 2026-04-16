package parser_test

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/parser"
	"testing"
)

func TestIdentifierExpression(t *testing.T) {
	input := "return foobar;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ReturnStatement)
	ident, ok := stmt.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Value)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "let x: int = 5;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.LetStatement)
	literal, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Value)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "let x: bool = true;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.LetStatement)
	ident, ok := stmt.Value.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Value)
	}
	if ident.Value != true {
		t.Errorf("ident.Value not %v. got=%v", true, ident.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `let x: string = "hello world";`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.LetStatement)
	literal, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Value)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "let x: [int] =[1, 2 * 2, 3 + 3];"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.LetStatement)
	array, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Value)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `let x: {string : int} = {"one": 1, "two": 2, "three": 3};`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.LetStatement)
	hash, ok := stmt.Value.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashMapLiteral. got=%T", stmt.Value)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int{
		"\"one\"":   1,
		"\"two\"":   2,
		"\"three\"": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `let x: {string : int} = {"one": 0 + 1, "two": 10 - 8, "three": 15 / 5};`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.LetStatement)
	hash, ok := stmt.Value.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashMapLiteral. got=%T", stmt.Value)
	}

	tests := map[string]func(ast.ExpressionNode){
		"\"one\"":   func(e ast.ExpressionNode) { testInfixExpression(t, e, 0, "+", 1) },
		"\"two\"":   func(e ast.ExpressionNode) { testInfixExpression(t, e, 10, "-", 8) },
		"\"three\"": func(e ast.ExpressionNode) { testInfixExpression(t, e, 15, "/", 5) },
	}

	for key, value := range hash.Pairs {
		literal, _ := key.(*ast.StringLiteral)
		testFunc := tests[literal.String()]
		testFunc(value)
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "let x: int = myArray[1 + 1];"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.LetStatement)
	indexExp, ok := stmt.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Value)
	}
	testIdentifierLiteral(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `let myFunc: fn(int, int) -> int = fn(x, y) { return x + y; };`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.LetStatement)
	varType, ok := stmt.Type.(*ast.FunctionType)
	if !ok {
		t.Fatalf("stmt.Type is not ast.FunctionType. got=%T", stmt.Type)
	}

	if varType.ParamsTypes[0].TokenLiteral() != "int" || varType.ParamsTypes[1].TokenLiteral() != "int" {
		t.Errorf("wrong param types parsed")
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T", stmt.Value)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	bodyStmt, _ := function.Body.Statements[0].(*ast.ReturnStatement)
	testInfixExpression(t, bodyStmt.Value, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"let x: fn(int, int) -> bool = fn(x, y) { return true; };", []string{"x", "y"}},
		{"let x: fn() -> bool = fn() { return true; };", []string{}},
		{"let x: fn(int, int, bool) -> bool = fn(x, y, z) { return true; };", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)
		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		stmt := prg.Statements[0].(*ast.LetStatement)
		function := stmt.Value.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestParsingCallExpressions(t *testing.T) {
	input := "let x: int = add(1, 2 * 3, 4 + 5);"
	lxr := lexer.New(input)
	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	checkParserErrors(t, prs)

	stmt, _ := prg.Statements[0].(*ast.LetStatement)
	exp, ok := stmt.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Value)
	}

	testIdentifierLiteral(t, exp.Function, "add")
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"let x: int = -15;", "-", 15},
		{"let x: bool = !false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, _ := program.Statements[0].(*ast.LetStatement)
		exp, ok := stmt.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Value)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"let x: int = 5 + 5;", 5, "+", 5},
		{"let x: int = 5 - 5;", 5, "-", 5},
		{"let x: int = 5 * 5;", 5, "*", 5},
		{"let x: int = 5 / 5;", 5, "/", 5},
		{"let x: int = 5 > 5;", 5, ">", 5},
		{"let x: int = 5 < 5;", 5, "<", 5},
		{"let x: int = false == false;", false, "==", false},
		{"let x: int = true != true;", true, "!=", true},
		{"let x: int = true && true;", true, "&&", true},
		{"let x: int = true || true;", true, "||", true},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, _ := program.Statements[0].(*ast.LetStatement)
		testInfixExpression(t, stmt.Value, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestParsingIfExpression(t *testing.T) {
	input := "let x: bool = if (x < y) { let y: int = 5; return x + y; } else { return 5; };"
	lxr := lexer.New(input)
	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	checkParserErrors(t, prs)

	stmt, _ := prg.Statements[0].(*ast.LetStatement)
	exp, ok := stmt.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Value)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	consequence1, _ := exp.Consequence.Statements[0].(*ast.LetStatement)
	consequence2, _ := exp.Consequence.Statements[1].(*ast.ReturnStatement)
	testLiteralExpression(t, consequence1.Value, 5)
	testInfixExpression(t, consequence2.Value, "x", "+", "y")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"let x: bool = a + b * c == 3 * 1 + 4 * 5 == true;",
			"let x: bool = (((a + (b * c)) == ((3 * 1) + (4 * 5))) == true);",
		},
		{
			"let x: int = (5 + 5) * 2;",
			"let x: int = ((5 + 5) * 2);",
		},
		{
			"let x: int = a *[1, 2, 3, 4][b * c] * d;",
			"let x: int = ((a * ([1, 2, 3, 4][(b * c)])) * d);",
		},
	}
	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)
		program := prs.ParseProgram()
		checkParserErrors(t, prs)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
