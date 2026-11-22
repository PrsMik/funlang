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
			expectedSig: "<int>",
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
			expectedSig: "<interface Greeter>",
		},
		{
			input: `
			let Box: fn(T: type) -> type = fn(T) {
				return interface {
					let value: type = T;
				};
			};
			let IntBox: type = Box(int);
			return IntBox;`,
			expectedSig: "<interface @2:11>",
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
			let Geometric: type = interface { let x: type = int; let y: type = int; };
			let Point: type = impl (Geometric) {
				let x: int = 42;
				let y: int = 42;
			};
			let p: Geometric = Point.{};
			return p.x + p.y;`,
			expected: 84,
		},
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
			let Config: type = impl {
				let timeout: int = 40 + 2;
			};
			let c: Config = Config.{};
			return c.timeout;`,
			expected: 42,
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

func TestEvalTypeCast(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Safe cast success interface -> impl",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; };
				let a: Animal = Dog.{};
				return a.(Dog);
			`,
			expected: "<impl Dog (<interface Animal>)>",
		},
		{
			name: "Safe cast success to itself",
			input: `
				let x: int = 42;
				return x.(int);
			`,
			expected: "42",
		},
		{
			name: "Safe cast unsuccess interface -> wrong impl",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; };
				let Cat := impl Animal { let name: string = "Felix"; };
				let a: Animal = Dog.{};
				return a.(Cat);
			`,
			expected: "nil",
		},
		{
			name: "Safe cast unsuccess simple",
			input: `
				let x: int = 42;
				return x.(string);
			`,
			expected: "nil",
		},
		{
			name: "Unsafe cast unsuccess interface -> impl",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; };
				let a: Animal = Dog.{};
				return .Dog(a);
			`,
			expected: "<impl Dog (<interface Animal>)>",
		},
		{
			name: "Unsafe cast success",
			input: `
				let x: int = 100;
				return .int(x);
			`,
			expected: "100",
		},
		{
			name: "Unsafe cast unsuccess interface -> wrong impl",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; };
				let Cat := impl Animal { let name: string = "Felix"; };
				let a: Animal = Dog.{};
				return .Cat(a);
			`,
			expected: "RUNTIME ERROR: runtime error: cannot unsafe cast <impl Dog (<interface Animal>)> to <impl Cat (<interface Animal>)>",
		},
		{
			name: "Unsafe cast unsuccess",
			input: `
				let x: int = 42;
				return .string(x);
			`,
			expected: "RUNTIME ERROR: runtime error: cannot unsafe cast <int> to <string>",
		},
		{
			name: "Safe cast with switch narrowing",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; let sound: string = "woof"; };
				let Cat := impl Animal { let name: string = "Felix"; let sound: string = "meow"; };
				
				let getSound := fn(a: Animal) -> string {
					let maybeDog := a.(Dog);
					return switch maybeDog.(type) {
						case Dog { 
							return maybeDog.sound; 
						}
						default  { 
							let maybeCat := a.(Cat);
							return switch maybeCat.(type) {
								case Cat { return maybeCat.sound; }
								default  { return "unknown"; }
							};
						}
					};
				};
				
				let a: Animal = Dog.{};
				return getSound(a);
			`,
			expected: "woof",
		},
		{
			name: "Safe cast unsuccess switch default",
			input: `
				let Animal := interface { let name: type = string; };
				let Dog := impl Animal { let name: string = "Rex"; };
				let Bird := impl Animal { let name: string = "Tweety"; };
				
				let a: Animal = Bird.{};
				let maybeDog := a.(Dog);
				return switch maybeDog.(type) {
					case Dog { return "is dog"; }
					default  { return "not a dog"; }
				};
			`,
			expected: "not a dog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(t, tt.input)

			if len(tt.expected) >= 13 && tt.expected[:13] == "RUNTIME ERROR" {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				}
				if errObj.Inspect() != tt.expected {
					t.Errorf("expected error message %q, got %q", tt.expected, errObj.Inspect())
				}
				return
			}

			if evaluated.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, evaluated.Inspect())
			}
		})
	}
}

func TestEvalResultErrorHandlingAndTypeSwitch(t *testing.T) {
	inputOk := `
	let OkType: fn(T: type) -> type = comptime fn(T) {
		return interface {
			let value: type = T;
		};
	};

	let ErrType: fn(E: type) -> type = comptime fn(E) {
		return interface {
			let error: type = E;
		};
	};

	let Result: fn(T: type, E: type) -> type = comptime fn(T, E) {
		return OkType(T) | ErrType(E); 
	};

	let MakeOkFor: fn(T: comptime type) -> fn(T) -> comptime OkType(T) = comptime fn(T) {
		return fn(val: comptime T) -> comptime OkType(T) {
			let ok: type = impl (comptime OkType(T)) {
				let value: T = val;
			};
			return ok.{value : val};
		};
	};

	let MakeErrFor: fn(E: comptime type) -> fn(E) -> comptime ErrType(E) = comptime fn(E) {
		return fn(err_val: comptime E) -> comptime ErrType(E) {
			let err: type = impl (comptime ErrType(E)) {
				let error: E = err_val;
			};
			return err.{error : err_val};
		};
	};

	let resTp: type = comptime Result(int, string);
	let divide: fn(int, int) -> resTp = fn(a, b) {
		return if (b == 0) {
			let res: resTp = (comptime MakeErrFor(string))("Division by zero");
			return res; 
		} else {
			let res: resTp = (comptime MakeOkFor(int))(a / b);
			return res;
		};
	};

	let res: comptime Result(int, string) = divide(10, 2);

	return switch res.(type) {
		case comptime OkType(int) {
			return res.value;
		}
		case comptime ErrType(string) {
			return -1;
		}
		default {
			return -999;
		}
	};`

	evaluatedOk := testEval(t, inputOk)
	testIntegerObject(t, evaluatedOk, 5)

	inputErr := `
	let OkType: fn(T: type) -> type = comptime fn(T) {
		return interface {
			let value: type = T;
		};
	};

	let ErrType: fn(E: type) -> type = comptime fn(E) {
		return interface {
			let error: type = E;
		};
	};

	let Result: fn(T: type, E: type) -> type = comptime fn(T, E) {
		return OkType(T) | ErrType(E); 
	};

	let MakeOkFor: fn(T: comptime type) -> fn(T) -> comptime OkType(T) = comptime fn(T) {
		return fn(val: comptime T) -> comptime OkType(T) {
			let ok: type = impl (comptime OkType(T)) {
				let value: T = val;
			};
			return ok.{value : val};
		};
	};

	let MakeErrFor: fn(E: comptime type) -> fn(E) -> comptime ErrType(E) = comptime fn(E) {
		return fn(err_val: comptime E) -> comptime ErrType(E) {
			let err: type = impl (comptime ErrType(E)) {
				let error: E = err_val;
			};
			return err.{error : err_val};
		};
	};

	let resTp: type = comptime Result(int, string);
	let divide: fn(int, int) -> resTp = fn(a, b) {
		return if (b == 0) {
			let res: resTp = (comptime MakeErrFor(string))("Division by zero");
			return res; 
		} else {
			let res: resTp = (comptime MakeOkFor(int))(a / b);
			return res;
		};
	};

	let res: comptime Result(int, string) = divide(10, 0);

	return switch res.(type) {
		case comptime ErrType(string) {
			return res.error;
		}
		case comptime OkType(int) {
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
