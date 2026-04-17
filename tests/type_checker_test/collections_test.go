package type_checker_test

import "testing"

func TestArrays(t *testing.T) {
	tests := []TestCase{
		{"Valid empty array", "let x: [int] =[];", ""},
		{"Valid int array", "let x: [int] = [1, 2, 3];", ""},
		{"Mismatch element types", "let x: [int] = [1, true, 3];", "type error: array literal has elements of different types <int> & <bool>"},
		{"Nested array", "let x: [[int]] = [[1, 2], [3]];", ""},
	}
	runTypeCheckerTests(t, tests)
}

func TestHashMaps(t *testing.T) {
	tests := []TestCase{
		{"Valid int map", `let x: {string : int} = {"a": 1, "b": 2};`, ""},
		{"Mismatch map keys", `let x: {string : int} = {"a": 1, true: 2};`, "type error: map literal has keys of different types <bool> & <string>"},
		{"Mismatch map values", `let x: {string : int} = {"a": 1, "b": true};`, "type error: map literal has elements of different types <int> & <bool>"},
		{"Invalid map key type", `let x: {[int] : int} = {};`, "type error: cannot use type <[<int>]> for hash map key"},
	}
	runTypeCheckerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []TestCase{
		{"Array index int", "let x: int =[1, 2, 3][0];", ""},
		{"Array index string (invalid)", `let x: int =[1, 2]["a"];`, "array index expression has non-integer index"},
		{"Map index", `let x: int = {"a": 1}["a"];`, ""},
		{"Map wrong key type", `let x: int = {"a": 1}[0];`, "type error: type mismatch between <int> in index & <string> in keys for index operator in hash map"},
		{"Index on non-collection", "let x: int = 5[0];", "index expression has wrong left operand"},
	}
	runTypeCheckerTests(t, tests)
}
