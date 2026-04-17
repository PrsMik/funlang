package type_checker_test

import "testing"

func TestLetStatements(t *testing.T) {
	tests := []TestCase{
		{"Valid int", "let x: int = 5;", ""},
		{"Valid string", `let x: string = "hello";`, ""},
		{"Valid bool", "let y: bool = true;", ""},
		{"Type mismatch int/bool", "let z: int = true;", "type error: expected type <int>, got <bool>"},
		{"Type mismatch bool/int", "let z: bool = -5;", "type error: expected type <bool>, got <int>"},
		{"Complex valid int", "let z: int = 5 + 5 * 5;", ""},
		{"Complex valid bool", "let z: bool = false && true || !false;", ""},
		{"Undefined identifier in let", "let y: int = unknown_var;", "unknown identifier"},
		{"Let with array", `let x: [string] = ["hello", "world!"];`, ""},
		{"Let with map", `let x: {string : int} = {"hello" : 1};`, ""},
		{"Let with empty map", `let x: {string : int} = {};`, ""},
	}
	runTypeCheckerTests(t, tests)
}

func TestReturnStatements(t *testing.T) {
	tests := []TestCase{
		{"Return matching types", "let x: int = if (2 > 1) { return 5; } else { return 2; };", ""},
		{"Return arrays", "let x: [int] = if (2 > 1) { return []; } else { return[2, 3]; };", ""},
		{"Return maps", "let x: {int : int} = if (2 > 1) { return {}; } else { return {2 : 3}; };", ""},
		{"Return bool complex", "let x: bool = if (true) { return true && false || false; } else { return true; };", ""},
		{"Condition not bool", "let x: bool = if (2 + 1) { return true; } else { return true; };", "type error: wrong type <int> for if condition"},
		{"Return branches mismatch", "let x: int = if (true) { return 5; } else { return true; };", "type error: type mismatch between <int> & <bool> in if/else branches"},
	}
	runTypeCheckerTests(t, tests)
}
