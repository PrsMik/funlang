package evaluator_test

import (
	"funlang/object"
	"testing"
)

func testTypeObject(t *testing.T, obj object.Object, expectedSig string) bool {
	result, ok := obj.(*object.TypeObject)
	if !ok {
		t.Errorf("object is not TypeObject. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Inspect() != expectedSig {
		t.Errorf("TypeObject has wrong signature. got=%q, want=%q", result.Inspect(), expectedSig)
		return false
	}
	return true
}

func TestEvalFirstClassTypes(t *testing.T) {
	tests := []struct {
		input       string
		expectedSig string
	}{
		{
			input:       `let T: type = int; return T;`,
			expectedSig: "int",
		},
		{
			input:       `let NumericOrString: type = int | string; return NumericOrString;`,
			expectedSig: "<union <int> | <string>>",
		},
		{
			input:       `let NumericAndString: type = int & string; return NumericAndString;`,
			expectedSig: "<intersection <int> & <string>>",
		},
		{
			input: `
			let Greeter: type = interface {
				let greet: type = fn() -> string;
			};
			return Greeter;`,
			expectedSig: "Greeter",
		},
		{
			input: `
			let Box: fn(T: type) -> type = fn(T: type) -> type {
				return interface {
					let value: type = T;
				};
			};
			let IntBox: type = Box(int);
			return IntBox;`,
			expectedSig: "<inteface >",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testTypeObject(t, evaluated, tt.expectedSig)
	}
}

func TestEvalImplInstantiationAndFieldAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			input: `
			let Point: type = impl {
				let x: int = 0;
				let y: int = 0;
			};
			let p: Point = Point.{x : 10, y : 20};
			return p.x + p.y;`,
			expected: 30,
		},
		{
			input: `
			let Point: type = impl {
				let x: int = 5;
				let y: int = 15;
			};
			let p: Point = Point.{x : 100};
			return p.x + p.y;`,
			expected: 115,
		},
		{
			input: `
			let Config: type = impl {
				let timeout: int = 30;
			};
			let c: Config = Config.{};
			return c.timeout;`,
			expected: 30,
		},
		{
			input: `
			let Box: type = impl {
				let size: int = 0;
			};
			let b: Box = Box.{size : 10 * 2 + 5};
			return b.size;`,
			expected: 25,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalImplMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			input: `
			let Calculator: type = impl {
				let add: fn(int, int) -> int = fn(a: int, b: int) -> int {
					return a + b;
				};
			};
			let calc: Calculator = Calculator.{};
			return calc.add(15, 25);`,
			expected: 40,
		},
		{
			input: `
			let Vector: type = impl {
				let x: int = 0;
				let y: int = 0;
			};

			let newVector: fn(int, int) -> Vector = fn(xT: int, yT: int) -> Vector {
				return Vector.{x : xT, y : yT};
			};

			let v: Vector = newVector(7, 3);
			return v.x * v.y;`,
			expected: 21,
		},
		{
			input: `
			let Inner: type = impl {
				let val: int = 42;
			};
			let Outer: type = impl {
				let in: Inner = Inner.{};
			};
			let o: Outer = Outer.{};
			return o.in.val;`,
			expected: 42,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalResultErrorHandlingAndTypeSwitch(t *testing.T) {
	inputOk := `
	let OkType: fn(T: type) -> type = fn(T: type) -> type {
		return interface {
			let value: type = T;
		};
	};

	let ErrType: fn(E: type) -> type = fn(E: type) -> type {
		return interface {
			let error: type = E;
		};
	};

	let Result: fn(T: type, E: type) -> type = fn(T: type, E: type) -> type {
		return OkType(T) | ErrType(E);
	};

	let OkImpl: type = impl (OkType(int)) {
		let value: int = 0;
	};

	let ErrImpl: type = impl (ErrType(string)) {
		let error: string = "";
	};

	let divide: fn(int, int) -> Result(int, string) = fn(a: int, b: int) -> Result(int, string) {
		if (b == 0) {
			return ErrImpl.{error : "Division by zero"};
		} else {
			return OkImpl.{value : a / b};
		};
	};

	let res: Result(int, string) = divide(10, 2);

	return switch res.(type) {
		case OkType(int) {
			return res.value;
		}
		case ErrType(string) {
			return -1;
		}
		default {
			return -999;
		}
	};`

	evaluatedOk := testEval(t, inputOk)
	testIntegerObject(t, evaluatedOk, 5)

	inputErr := `
	let OkType: fn(T: type) -> type = fn(T: type) -> type {
		return interface {
			let value: type = T;
		};
	};

	let ErrType: fn(E: type) -> type = fn(E: type) -> type {
		return interface {
			let error: type = E;
		};
	};

	let Result: fn(T: type, E: type) -> type = fn(T: type, E: type) -> type {
		return OkType(T) | ErrType(E);
	};

	let OkImpl: type = impl (OkType(int)) {
		let value: int = 0;
	};

	let ErrImpl: type = impl (ErrType(string)) {
		let error: string = "";
	};

	let divide: fn(int, int) -> Result(int, string) = fn(a: int, b: int) -> Result(int, string) {
		if (b == 0) {
			return ErrImpl.{error : "Division by zero"};
		} else {
			return OkImpl.{value : a / b};
		};
	};

	let res: Result(int, string) = divide(10, 0);

	return switch res.(type) {
		case ErrType(string) {
			return res.error;
		}
		case OkType(int) {
			return "OK";
		}
		default {
			return "UNKNOWN";
		}
	};`

	evaluatedErr := testEval(t, inputErr)
	str, ok := evaluatedErr.(*object.String)
	if !ok {
		t.Fatalf("expected string result from error branch, got=%T (%+v)", evaluatedErr, evaluatedErr)
	}
	if str.Value != "Division by zero" {
		t.Errorf("wrong error string. got=%q, want=%q", str.Value, "Division by zero")
	}
}
