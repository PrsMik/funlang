package evaluator_test

import (
	"funlang/object"
	"testing"
)

func TestEvalFunctionObject(t *testing.T) {
	input := "return fn(x: int) -> int { return x + 2; };"
	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "return (x + 2);"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestEvalFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let identity: fn(int) -> int = fn(x) { return x; }; return identity(5);", 5},
		{"let double: fn(int) -> int = fn(x) { return x * 2; }; return double(5);", 10},
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; return add(5, 5);", 10},
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

func TestEvalClosures(t *testing.T) {
	input := `
		let newAdder: fn(int) -> fn(int) -> int = fn(x: int) -> fn(int) -> int {
			return fn(y: int) -> int { return x + y; };
		};
		let addTwo: fn(int) -> int = newAdder(2);
		return addTwo(2);`

	testIntegerObject(t, testEval(t, input), 4)

	inputDeepClosure := `
		return fn(x: int) -> fn(int, int) -> fn() -> int {
				return fn(y: int, z: int) -> fn() -> int { 
					return fn() -> int { return x + y + z; }; 
				}; 
		}(1)(1, 1)();`

	testIntegerObject(t, testEval(t, inputDeepClosure), 3)
}
