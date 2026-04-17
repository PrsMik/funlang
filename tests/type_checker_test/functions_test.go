package type_checker_test

import "testing"

func TestFunctionLiterals(t *testing.T) {
	tests := []TestCase{
		{"Valid function", "let x: fn(int, int) -> int = fn(x, y) { return x + y; };", ""},
		{"Valid string function", "let x: fn(string, string) -> string = fn(x, y) { return x + y; };", ""},
		{"Math on strings in function (invalid)", "let x: fn(string, string) -> string = fn(x, y) { return x / y; };", "type error: type mismatch between: <string> & <string>; for operator /"},
		{"Return type mismatch", "let x: fn(int, int) -> bool = fn(x, y) { return x + y; };", "type error: function literal has return type <int>, but expected <bool>"},
		{"Params count mismatch", "let x: fn(int, int) -> int = fn(x) { return x; };", "function literal has 1 parameters, but expected 2"},
		{"Closure capturing variables", `let x: fn(bool, bool, int) -> int = fn(x, y, z) { 
			let n: int = z + if (x && y) { return 1; } else { return 0; };
			return z; 
		};`, ""},
		{"Recursive function", "let add: fn(int, int) -> int = fn(x, y) { return add(x, y); };", ""},
	}
	runTypeCheckerTests(t, tests)
}

func TestCallExpressions(t *testing.T) {
	tests := []TestCase{
		{"Valid call", "let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = add(1, 2);", ""},
		{"Valid nested call", "let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = 1 + add(1, 2);", ""},
		{"Call wrong arg count", "let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = add(1);", "wrong number of arguments"},
		{"Call wrong arg type", "let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = add(1, true);", "type error: wrong type for argument 2: <bool> in func call \" add \"; expected <int>"},
		{"Call non-function", "let x: int = 5; let y: int = x(1);", "identifier \"x\" is no a funciton"},
		{"Valid array of functions", `let y: fn() -> int = fn() { return 1; }; let x: [fn() -> int] =[y, y, y];`, ""},
	}
	runTypeCheckerTests(t, tests)
}
