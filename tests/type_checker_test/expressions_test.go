package type_checker_test

import "testing"

func TestPrefixExpressions(t *testing.T) {
	tests := []TestCase{
		{"Minus int", "let z: int = -5;", ""},
		{"Minus int infer", "let z := -5;", ""},
		{"Minus bool (invalid)", "let z: int = -true;", "type error: unknown operator: - for type <bool>"},
		{"Minus bool (invalid) infer", "let z := -true;", "type error: unknown operator: - for type <bool>"},
		{"Bang bool", "let x: bool = !true;", ""},
		{"Bang int (invalid)", "let x: bool = !5;", "type error: unknown operator: ! for type <int>"},
	}
	runTypeCheckerTests(t, tests)
}

func TestInfixExpressions(t *testing.T) {
	tests := []TestCase{
		{"Int math", "let x: int = 5 + 5 - 2 * 3 / 1;", ""},
		{"Int math infer", "let x := 5 + 5 - 2 * 3 / 1;", ""},
		{"String concat", `let x: string = "hello" + "world!";`, ""},
		{"String minus (invalid)", `let x: string = "a" - "b";`, "type error: type mismatch between: <string> & <string>; for operator -"},
		{"Bool logic", "let x: bool = true && false || true;", ""},
		{"Int comparison", "let x: bool = 5 > 3 && 2 <= 4;", ""},
		{"Equals ints", "let x: bool = 5 == 5;", ""},
		{"Equals bools", "let x: bool = true != false;", ""},
		{"Equals type mismatch", "let x: bool = 5 == true;", "type error: type mismatch between: <int> & <bool>; for operator =="},
		{"Math with bool (invalid)", "let x: int = 5 + true;", "type error: type mismatch between: <int> & <bool>; for operator +"},
	}
	runTypeCheckerTests(t, tests)
}

func TestIdentifiers(t *testing.T) {
	tests := []TestCase{
		{"Valid shadowing and scope", "let x: int = 5; let y: int = x + 3;", ""},
		{"Valid shadowing and scope infer", "let x := 5; let y := x + 3;", ""},
		{"Use before declare (invalid)", "let y: int = x + 3; let x: int = 5;", "unknown identifier: x"},
		{"Use before declare (invalid) infer", "let y := x + 3; let x := 5;", "unknown identifier: x"},
		{"Self assignment (invalid)", "let x: int = x;", "unknown identifier: x"},
		{"Self assignment (invalid) infer", "let x := x;", "unknown identifier: x"},
		{"Valid identifier array", `let y: []int = []int{1, 2, 3}; let x: []int = y;`, ""},
		{"Valid identifier array infer", `let y := []int{1, 2, 3}; let x: []int = y;`, ""},
		{"Valid identifier nested scope", "let z: int = 5; let y: int = if (z > 1) { let x: int = 1; return x + z; } else { return 1; };", ""},
	}
	runTypeCheckerTests(t, tests)
}
