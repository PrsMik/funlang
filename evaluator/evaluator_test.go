package evaluator

import (
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"funlang/type_checker"
	"funlang/types"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10;", 10},
		{"return 10; let y: int = 9;", 10},
		{"return 2 * 5; let y: int = 9;", 10},
		{"let x: int = 9; return 2 * 5; let y: int = 9;", 10},
		{"let x: int = if (2 > 1) { let x: int = 1; return 1; return 10; } else { return 2; return 20; };", 1},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let a: int = 5; return a;", 5},
		{"let a: int = 5 * 5; return a;", 25},
		{"let a: int = 5; let b: int = a; return b;", 5},
		{"let a: int = 5; let b: int = a; let c: int = a + b + 5; return c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let x: int = 5;", 5},
		{"let x: int = 10;", 10},
		{"let x: int =-10;", -10},
		{"let x: int =5 + 5 + 5 + 5 - 10;", 10},
		{"let x: int =2 * 2 * 2 * 2 * 2;", 32},
		{"let x: int =-50 + 100 + -50;", 0},
		{"let x: int =5 * 2 + 10;", 20},
		{"let x: int =5 + 2 * 10;", 25},
		{"let x: int =20 + 2 * -10;", 0},
		{"let x: int =50 / 2 * 2 + 10;", 60},
		{"let x: int =2 * (5 + 10);", 30},
		{"let x: int =3 * 3 * 3 + 10;", 37},
		{"let x: int =3 * (3 * 3) + 10;", 37},
		{"let x: int =(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		n        int
		input    string
		expected bool
	}{
		{1, "let x: bool = true;", true},
		{2, "let x: bool = false;", false},
		{3, "let x: bool = 1 < 2;", true},
		{4, "let x: bool = 1 > 2;", false},
		{5, "let x: bool = 1 < 1;", false},
		{6, "let x: bool = 1 > 1;", false},
		{7, "let x: bool = 1 == 1;", true},
		{8, "let x: bool = 1 != 1;", false},
		{9, "let x: bool = 1 == 2;", false},
		{10, "let x: bool = 1 != 2;", true},
		{11, "let x: bool = true == true;", true},
		{12, "let x: bool = false == false;", true},
		{13, "let x: bool = true == false;", false},
		{14, "let x: bool = true != false;", true},
		{15, "let x: bool = false != true;", true},
		{16, "let x: bool = (1 < 2) == true;", true},
		{17, "let x: bool = (1 < 2) == false;", false},
		{18, "let x: bool = (1 > 2) == true;", false},
		{19, "let x: bool = (1 > 2) == false;", true},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, tt.n, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let x: int = if (true) { return 10; } else { return 0; };", 10},
		{"let x: int = if (false) { return 10; } else { return 0; };", 0},
		{"let x: int = if (1 < 2) { return 10; } else { return 20; };", 10},
		{"let x: int = if (1 > 2) { return 10; } else { return 20; };", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"let x: bool = !true;", false},
		{"let x: bool = !false;", true},
		{"let x: bool = !!true;", true},
		{"let x: bool = !!false;", false},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, 0, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "let y: fn(int) -> int = fn(x) { return x + 2; };"
	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "return (x + 2);"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let identity: fn(int) -> int = fn(x) { return x; }; return identity(5);", 5},
		{"let double: fn(int) -> int = fn(x) { return x * 2; }; return double(5);", 10},
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; return add(5, 5);", 10},
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; return add(5 + 5, add(5, 5));", 20},
		// {"let y: int = fn(x: int) -> int { return x; }(5); return y;", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	type_env := types.NewTypeEviroment()
	chk := type_checker.New(type_env)
	chk.CheckProgram(program)
	checkCheckerErrors(t, chk)

	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, n int, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Test %d. object is not Boolean. got=%T (%+v)", n, obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Test %d. object has wrong value. got=%t, want=%t",
			n, result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, prs *parser.Parser) {
	errors := prs.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func checkCheckerErrors(t *testing.T, chk *type_checker.TypeChecker) {
	errors := chk.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Checker has %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("checker error: %q", msg)
	}
	t.FailNow()
}
