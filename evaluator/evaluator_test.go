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

func TestStringLiteral(t *testing.T) {
	input := `let x: string = "Hello World!";`
	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "let x: [int] = [1, 2 * 2, 3 + 3];"
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"let x: int = [1, 2, 3][0];",
			1,
		},
		{
			"let x: int = [1, 2, 3][1];",
			2,
		},
		{
			"let x: int = [1, 2, 3][2];",
			3,
		},
		{
			"let i: int = 0; let x: int = [1][i];",
			1,
		},
		{
			"let x: int = [1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray: [int] = [1, 2, 3]; let x: int = myArray[2];",
			3,
		},
		{
			"let myArray: [int] = [1, 2, 3]; let x: int = myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray: [int] = [1, 2, 3]; let i: int = myArray[0]; let x: int = myArray[i];",
			2,
		},
		{
			"let x: int = [1, 2, 3][3];",
			nil,
		},
		{
			"let x: int = [][3];",
			nil,
		},
		{
			"let x: int = [1, 2, 3][-1];",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two: string = "two";
			  let map: {string : int} = {
				"one": 10 - 9,
				two: 1 + 1,
				"thr" + "ee": 6 / 2 };`
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`let x: int = {"foo": 5}["foo"];`,
			5,
		},
		{
			`let x: int = {"foo": 5}["bar"];`,
			nil,
		},
		{
			`let key: string = "foo"; let x: int = {"foo": 5}[key];`,
			5,
		},
		// {
		// 	`let x: int = {}["foo"];`,
		// 	nil,
		// },
		{
			`let x: int = {5: 5}[5];`,
			5,
		},
		{
			`let x: int = {true: 5}[true];`,
			5,
		},
		{
			`let x: int = {false: 5}[false];`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `let x: string = "Hello" + " " + "World!";`
	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
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
		{"let add: fn(int, int) -> int = fn(x: int, y: int) -> int { return x + y; }; return add(5, 5);", 10},
		{"let add: fn() -> int = fn() { return 1; }; return add();", 1},
		{"return fn() -> int { return 1; }();", 1},
		{"return fn() -> fn() -> int { return fn() -> int { return 1; }; }()();", 1},
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; return add(5 + 5, add(5, 5));", 20},
		{"let y: int = fn(x: int) -> int { return x; }(5); return y;", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	// input := `
	// 	let newAdderX: fn(int) -> fn(int) -> fn() -> int = fn(x) {
	// 			return fn(y) { return fn() { return x + y; }; };
	// 	};
	// 	let setY: fn(int) -> fn() -> int = newAdderX(2);
	// 	let addTwo: fn() -> int = setY(1);
	// 	return addTwo();`
	input := `
		return fn(x: int) -> fn(int, int) -> fn() -> int {
				return fn(y, z) -> fn() -> int { return fn() -> int { return x + y + z; }; }; }(1)(1, 1)();`
	testIntegerObject(t, testEval(t, input), 3)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`return len("");`, 0},
		{`return len("four");`, 4},
		{`return len("hello world");`, 11},
		{`return len([]);`, 0},
		{`return len([1, 2, 3]);`, 3},
		{`return tail([1, 2, 3]);`, []int{2, 3}},
		{`return tail([]);`, nil},
		{`return push([2, 3], 1);`, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case []int:
			testArrayIntegerObject(t, evaluated, expected)
		case nil:
			_, ok := evaluated.(*object.Null)

			if !ok {
				t.Errorf("object is not Null. got=%T (%+v)", evaluated, evaluated)
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
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

func testArrayIntegerObject(t *testing.T, obj object.Object, expected []int) bool {
	result, ok := obj.(*object.Array)

	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("object has wrong size. got=%d, want=%d",
			len(result.Elements), expected)
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
	for _, err := range errors {
		t.Errorf("parser error: %q", err.Msg)
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
