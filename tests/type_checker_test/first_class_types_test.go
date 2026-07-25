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
			name: "Valid structural subtyping (extra fields allowed)",
			input: `
				let HasName: type = interface { let name: type = string; };
				let Person: type = impl (HasName) { let name: string = "Alice"; let age: int = 30; };
				let obj: HasName = Person;
			`,
			expectedErr: "",
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
				let obj: interface { let val: type = int; } = impl (interface { let val: type = int; }) { let val: int = 42; };
				let x: int = obj.val;
			`,
			expectedErr: "",
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
