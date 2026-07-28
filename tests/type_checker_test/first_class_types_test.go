package type_checker_test

import "testing"

func TestFirstClassTypes(t *testing.T) {
	tests := []TestCase{
		{
			name:        "Assign primitive type to type variable",
			input:       "let T: type = int; let x: T = 5;",
			expectedErr: "",
		},
		{
			name:        "Chain assign of types",
			input:       "let T: type = int; let U: type = T; let V: type = U; let x: V = 5;",
			expectedErr: "",
		},
		{
			name:        "Type mismatch using type variable",
			input:       "let T: type = string; let x: T = 5;",
			expectedErr: "expected type <string>, got <int>",
		},
		{
			name:        "Function returning a type (Generics foundation)",
			input:       "let MakeType: fn() -> type = fn() { return int; }; let x: MakeType() = 5;",
			expectedErr: "",
		},
		{
			name:        "Type mismatch with function returning a type",
			input:       "let MakeType: fn() -> type = fn() { return bool; }; let x: MakeType() = 5;",
			expectedErr: "expected type <bool>, got <int>",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestStructuralTyping(t *testing.T) {
	tests := []TestCase{
		{
			name: "Valid structural subtyping (exact match)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Person: type = impl (HasName) { let name: string = "Alice"; };
				let obj: HasName = Person;
			`,
			expectedErr: "",
		},
		{
			name: "Valid chain with structural subtyping (exact match)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let HasNameTwo: type = HasName;
				let Person: type = impl (HasName) { let name: string = "Alice"; };
				let obj: HasNameTwo = Person;
			`,
			expectedErr: "",
		},
		{
			name: "Valid impl",
			input: `
				let Person: type = impl { let name: string = "Alice"; };
			`,
			expectedErr: "",
		},
		{
			name: "Valid impl (interface)",
			input: `
				let obj: interface { let val: type = int; } = impl (interface { let val: type = int; }) { let val: int = 42; };
			`,
			expectedErr: "",
		},
		{
			name: "Valid impl (empty interface)",
			input: `
				let obj: interface {} = impl (interface {let val: type = int;}) { let val: int = 42; };
			`,
			expectedErr: "",
		},
		{
			name: "Valid structural subtyping (extra fields allowed)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Person: type = impl (HasName) { let name: string = "Alice"; let age: int = 30; };
				let obj: HasName = Person;
			`,
			expectedErr: "",
		},
		{
			name: "Valid structural subtyping (multiple implentation)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let HasYear: type = interface { let age: type = int; };
				let Person: type = impl (HasName & HasYear) { let name: string = "Alice"; let age: int = 30; };
				let obj: HasName = Person;
			`,
			expectedErr: "",
		},
		{
			name: "Invalid impl (interface)",
			input: `
				let obj: interface { let val: type = string; } = impl (interface {}) { let val: int = 42; };
			`,
			expectedErr: "type error",
		},
		{
			name: "Invalid structural subtyping (missing field)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Box: type = impl (HasName) { let size: int = 10; };
				let obj: HasName = Box;
			`,
			expectedErr: "type error",
		},
		{
			name: "Invalid structural subtyping (not implemented)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Box: type = impl { let size: string = "10"; };
				let obj: HasName = Box;
			`,
			expectedErr: "type error",
		},
		{
			name: "Invalid structural subtyping (wrong field type)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Person: type = impl (HasName) { let name: int = 123; };
				let obj: HasName = Person;
			`,
			expectedErr: "type error",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestAlgebraicTypes(t *testing.T) {
	tests := []TestCase{
		{
			name:        "Valid Union assignment (left side)",
			input:       "let x: int | string = 5;",
			expectedErr: "",
		},
		{
			name:        "Valid Union assignment (right side)",
			input:       `let x: int | string = "hello";`,
			expectedErr: "",
		},
		{
			name:        "Invalid Union assignment",
			input:       "let x: int | string = true;",
			expectedErr: "type error",
		},
		{
			name:        "Invalid Union assignment",
			input:       "let x: int | string = type;",
			expectedErr: "type error",
		},
		{
			name: "Valid Intersection assignment",
			input: `
				let A: type = interface { let a: type = int; };
				let B: type = interface { let b: type = string; };
				let obj: A & B = impl (A & B) { let a: int = 1; let b: string = "hello"; };
			`,
			expectedErr: "",
		},
		{
			name: "Invalid Intersection assignment (missing one requirement)",
			input: `
				let A: type = interface { let a: type = int; };
				let B: type = interface { let b: type = string; };
				let obj: A & B = impl (A & B) { let a: int = 1; };
			`,
			expectedErr: "type error",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestMemberAccess(t *testing.T) {
	tests := []TestCase{
		{
			name: "Valid member access",
			input: `
				let MyIface: type = interface { let val: type = int; };
				let obj: MyIface = impl (MyIface) { let val: int = 42; };
				let x: int = obj.val;
			`,
			expectedErr: "",
		},
		{
			name: "Valid member access (anonymous interfaces)",
			input: `
				let obj: interface { let val: type = int; } = impl (interface { let val: type = int; }) { let val: int = 42; };
				let x: int = obj.val;
			`,
			expectedErr: "",
		},
		{
			name: "Invalid member access",
			input: `
				let obj: type = impl { let val: int = 42; };
				let x: int = obj.val;
			`,
			expectedErr: "type error",
		},
		{
			name: "Invalid member access (unknown field)",
			input: `
				let obj: interface { let val: type = int; } = impl (interface { let val: type = int; }) { let val: int = 42; };
				let x: int = obj.unknown;
			`,
			expectedErr: "type error",
		},
		{
			name: "Member access type mismatch",
			input: `
				let obj: interface { let val: type = string; } = impl (interface { let val: type = string; }) { let val: string = "hi"; };
				let x: int = obj.val;
			`,
			expectedErr: "expected type <int>, got <string>",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestTypeAccessExpression(t *testing.T) {
	tests := []TestCase{
		{
			name: "Valid type access assignment",
			input: `
				let x: int = 5;
				let T: type = x.(type);
			`,
			expectedErr: "",
		},
		{
			name: "Type access assignment type mismatch",
			input: `
				let x: int = 5;
				let y: int = x.(type); // ожидается int, а получает type
			`,
			expectedErr: "expected type <int>, got <type(<int>)>",
		},
		{
			name: "Type access on complex expression",
			input: `
				let getNumber: fn() -> int = fn() { return 42; };
				let ResultType: type = getNumber().(type);
			`,
			expectedErr: "",
		},
		{
			name: "Type access on unresolved identifier",
			input: `
				let T: type = unknown_var.(type);
			`,
			expectedErr: "unknown identifier: unknown_var",
		},
		{
			name: "Type access on Union type",
			input: `
				let val: int | string = "hello";
				let ValType: type = val.(type);
			`,
			expectedErr: "",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestTypeNarrowing(t *testing.T) {
	tests := []TestCase{
		{
			name: "Valid narrowing in switch",
			input: `
				let res: int | string = 5;
				let final: int = switch res.(type) {
					case int { return res + 1; }
					case string { return len(res); }
				};
			`,
			expectedErr: "",
		},
		{
			name: "Invalid operation in narrowed context (catches type error correctly)",
			input: `
				let res: int | string = 5;
				let final: int = switch res.(type) {
					case int { return len(res); } // len expects string or array, not int
					case string { return len(res); }
				};
			`,
			expectedErr: "len does not support type",
		},
		{
			name: "Structural narrowing",
			input: `
				let OkType: type = interface { let value: type = int; };
				let ErrType: type = interface { let error: type = string; };
				
				let res: OkType | ErrType = impl (OkType) { let value: int = 5; };
				
				let final: int = switch res.(type) {
					case OkType { return res.value; } // Доступ разрешен, тип сужен
					case ErrType { return 0; }
				};
			`,
			expectedErr: "",
		},
		{
			name: "Invalid field access in narrowed context",
			input: `
				let OkType: type = interface { let value: type = int; };
				let ErrType: type = interface { let error: type = string; };
				
				let res: OkType | ErrType = impl (OkType) { let value: int = 5; };
				
				switch res.(type) {
					case OkType { return res.error; } // Ошибка: у OkType нет поля error
					case ErrType { return 0; }
				};
			`,
			expectedErr: "type error",
		},
	}
	runTypeCheckerTests(t, tests)
}
