package evaluator_test

import "testing"

func TestEvalLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let a: int = 5; return a;", 5},
		{"let a: int = 5 * 5; return a;", 25},
		{"let a: int = 5; let b: int = a; return b;", 5},
		{"let a: int = 5; let b: int = a; let c: int = a + b + 5; return c;", 15},
		{"let a: int = 5; let a: int = 10; return a;", 10},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10;", 10},
		{"return 10; let y: int = 9;", 10},
		{"return 2 * 5; let y: int = 9;", 10},
		{"let x: int = 9; return 2 * 5; let y: int = 9;", 10},
		{`let x: int = if (2 > 1) { 
			return if (3 > 1) { return 1; } else { return 10; };
		} else { return 2; }; 
		return x;`, 1},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
