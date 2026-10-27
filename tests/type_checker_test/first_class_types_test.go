package type_checker_test

import (
	"testing"
)

func TestFirstClassTypes(t *testing.T) {
	tests := []TestCase{
		{
			name:        "Assign primitive type to type variable",
			input:       "let T: type = int; let x: T = 5;",
			expectedErr: "",
		},
		{
			name:        "Assign type to type variable",
			input:       "let T: int = int;",
			expectedErr: "type error",
		},
		{
			name:        "Chain assign of types",
			input:       "let T: type = int; let U: type = T; let V: type = U; let x: V = 5;",
			expectedErr: "",
		},
		{
			name: "Simple impl and interface assignment",
			input: `let T: type = interface { let x: type = int; };
					let U: type = impl (T) { let x: int = 42; };
					let X: T = U.{x: 43};`,
			expectedErr: "",
		},
		{
			name: "Infinite metatype",
			input: `let T: type = interface { let x: type = type; };
					let U: type = impl (T) { let x: type = type; };
					let X: T = U.{x: type};
					let Z: X.x = int;`,
			expectedErr: "",
		},
		{
			name:        "Type mismatch using type variable",
			input:       "let T: type = string; let x: T = 5;",
			expectedErr: "in let expected variable of type <string>, got <int>",
		},
		{
			name:        "Function returning a type (Generics foundation) 1",
			input:       `let MakeType: fn() -> type = comptime fn() { return int; }; let x: comptime MakeType() = 5;`,
			expectedErr: "",
		},
		{
			name: "Function returning a type (Generics foundation) 1.1",
			input: `let MakeType: fn() -> type = comptime fn() { return int; }; 
					let mT: type = comptime MakeType();
					let x: mT = 5;`,
			expectedErr: "",
		},
		{
			name: "Function returning an interface (Generics foundation) 2",
			input: `let inter: type = interface { let x: type = int; };
					let im: type = comptime impl (inter) { let x: int = 42; };
					let MakeType: fn() -> inter = comptime fn() { return im; };
					let x: comptime MakeType() = im.{};`,
			expectedErr: "",
		},
		{
			name: "Function returning an interface (Generics foundation) 2.1",
			input: `let inter: type = interface { let x: type = int; };
					let im: type = comptime impl (inter) { let x: int = 42; };
					let MakeType: fn() -> inter = comptime fn() { return im; };
					let mT: type = comptime MakeType(); 
					let x: mT = im.{};`,
			expectedErr: "",
		},
		{
			name: "Function with typed params (Generics foundation) 3",
			input: `let makeGen: fn(T: type) -> (fn(T) -> T) = comptime fn(T) {
    					return fn(val: T) -> T {
							let x: T = val;
        					return val;
    					};
					};

					let genString: fn(string) -> string = comptime makeGen(string);
					let x: string = genString("42");`,
			expectedErr: "",
		},
		{
			name: "Generics interface",
			input: `let OkType: fn(T: type) -> type = comptime fn(T) {
						return interface{let value: type = T;};
					};

					let iface: type = comptime OkType(string);

					let ok: type = impl(iface){let value: string = "val";};

					return ok.{value:"val"};`,
			expectedErr: "",
		},
		{
			name: "Generics impl",
			input: `let MakeErr: fn(E: type) -> fn(E) -> type = comptime fn(E) {
						return fn(err_val: E) -> type {
							return impl {
								let error: E = err_val;
							};
						};
					};

					let err: fn(string) -> type = comptime MakeErr(string);
					let errImpl: comptime err("") = (comptime err("")).{error: "value"};`,
			expectedErr: "",
		},
		{
			name:        "Type mismatch with function returning a type",
			input:       "let MakeType: fn() -> type = comptime  fn() { return bool; }; let x: comptime MakeType() = 5;",
			expectedErr: "in let expected variable of type <bool>, got <int>",
		},
		{
			name: "Comptime generic type",
			input: `let x: int = comptime 2;
					let y: int = comptime 1;
					let MakeInt: fn() -> type = comptime fn() { 
						return int;
					};
					let MakeBool: fn() -> type = comptime fn() { 
						return bool;
					};					
					let MakeType: fn() -> type = comptime fn() { 
						return if (x > y) { return MakeInt(); } else { return MakeBool(); }; 
					};
					let z: comptime MakeType() = 5;`,
			expectedErr: "",
		},
	}
	runTypeCheckerTests(t, tests)
}

func TestStructuralTyping(t *testing.T) {
	tests := []TestCase{
		// {
		// 	name: "Valid structural subtyping (exact match)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let Person: type = impl (HasName) { let name: string = "Alice"; };
		// 		let obj: HasName = Person;
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Valid chain with structural subtyping (exact match)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let HasNameTwo: type = HasName;
		// 		let Person: type = impl (HasName) { let name: string = "Alice"; };
		// 		let obj: HasNameTwo = Person;
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Valid impl",
		// 	input: `
		// 		let Person: type = impl { let name: string = "Alice"; };
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Recursive interface assignability",
		// 	input: `
		// 		let List1: type = interface {
		// 			let next: type = fn() -> List1;
		// 		};

		// 		let impl1: type = impl (List1) {
		// 			let next: fn() -> impl1 = fn() { return impl1.{}; };
		// 		};

		// 		let l2: List1 = impl1;
		// `,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Recursive impl",
		// 	input: `
		// 		let impl1: type = impl {
		// 			let next: fn() -> impl1 = fn() { return impl1.{}; };
		// 		};

		// 		let l2: impl1 = impl1.{};
		// `,
		// 	expectedErr: "",
		// },
		{
			name: "Generic types",
			input: `
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

					let MakeOk: fn(T: type) -> fn(T) -> OkType(T) = comptime fn(T) {
						return fn(val: T) -> OkType(T) {
							let okTp: type = comptime OkType(T);
							let ok: type = impl (okTp) {
								let value: T = val;
							};
							return ok.{value : val};
						};
					};

					let MakeErr: fn(E: type) -> fn(E) -> ErrType(E) = comptime fn(E) {
						return fn(err_val: E) -> ErrType(E) {
							let errTp: type = comptime ErrType(E);
							let err: type = impl (errTp) {
								let error: E = err_val;
							};
							return err.{error : err_val};
						};
					};

					let divide: fn(int, int) -> comptime Result(int, string) = comptime fn(a, b) {
						return if (b == 0) {
							return (comptime MakeErr(string))("Division by zero"); 
						} else {
							return (comptime MakeOk(int))(a / b); 
						};
					};

					let res: comptime Result(int, string) = divide(10, 0);
    	`,
			expectedErr: "",
		},
		// {
		// 	name: "Valid impl (interface)",
		// 	input: `
		// 		let obj: interface { let val: type = int; } = impl (interface { let val: type = int; }) { let val: int = 42; };
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Valid impl (empty interface)",
		// 	input: `
		// 		let obj: interface {} = impl (interface {let val: type = int;}) { let val: int = 42; };
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Valid structural subtyping (extra fields allowed)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let Person: type = impl (HasName) { let name: string = "Alice"; let age: int = 30; };
		// 		let obj: HasName = Person;
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Valid structural subtyping (multiple implentation)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let HasYear: type = interface { let age: type = int; };
		// 		let Person: type = impl (HasName & HasYear) { let name: string = "Alice"; let age: int = 30; };
		// 		let obj: HasName = Person;
		// 	`,
		// 	expectedErr: "",
		// },
		// {
		// 	name: "Invalid impl (interface)",
		// 	input: `
		// 		let obj: interface { let val: type = string; } = impl (interface {}) { let val: int = 42; };
		// 	`,
		// 	expectedErr: "type error",
		// },
		// {
		// 	name: "Invalid impl (interface)",
		// 	input: `
		// 		let Iface: type = interface { let val: type = int; };
		// 		let Iface1: type = interface { let val1: type = string; };
		// 		let fun: fn(Iface1) -> string = fn(a) { return a.val1; };
		// 		let obj: Iface = impl (Iface) { let val: int = 42; };
		// 		let i: int = fun(obj);
		// 	`,
		// 	expectedErr: "type error",
		// },
		// {
		// 	name: "Invalid structural subtyping (missing field)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let Box: type = impl (HasName) { let size: int = 10; };
		// 		let obj: HasName = Box;
		// 	`,
		// 	expectedErr: "type error",
		// },
		// {
		// 	name: "Invalid structural subtyping (not implemented)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let Box: type = impl { let size: string = "10"; };
		// 		let obj: HasName = Box;
		// 	`,
		// 	expectedErr: "type error",
		// },
		// {
		// 	name: "Invalid structural subtyping (wrong field type)",
		// 	input: `
		// 		let HasName: type = interface { let name: type = string; };
		// 		let Person: type = impl (HasName) { let name: int = 123; };
		// 		let obj: HasName = Person;
		// 	`,
		// 	expectedErr: "type error",
		// },
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
			name:        "Invalid Union assignment 1",
			input:       "let x: int | string = true;",
			expectedErr: "type error",
		},
		{
			name:        "Invalid Union assignment 2",
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

func TestImplInstantiationTypeChecking(t *testing.T) {
	tests := []TestCase{
		{
			name: "Succesfull init (all fields)",
			input: `
			let Point: type = impl { let x: int = 0; let y: int = 1; };
			let p: Point = Point.{x : 10, y : 20};
			`,
			expectedErr: "",
		},
		{
			name: "Succesfull init (default values)",
			input: `
			let Point: type = impl { let x: int = 0; let y: int = 1; };
			let p: Point = Point.{x : 10};
			`,
			expectedErr: "",
		},
		{
			name: "Error (wrong field type)",
			input: `
			let Point: type = impl { let x: int = 0; let y: int = 1; };
			let p: Point = Point.{x : "wrong", y : 20};
			`,
			expectedErr: "type mismatch in instatation expression of <impl Point>",
		},
		{
			name: "Error (unknown field)",
			input: `
			let Point: type = impl { let x: int = 0; let y: int = 1; };
			let p: Point = Point.{z : 10};
			`,
			expectedErr: "type <impl Point> has no field named *z*",
		},
		{
			name: "Error (no impl type)",
			input: `
			let p: int = int.{x : 10};
			`,
			expectedErr: "type error: cannot use instantation .{} on a no-impl types",
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
			name: "Valid member access",
			input: `
				let x: int = ((impl (interface { let val: type = int; }) { let val: int = 42; }).{}).val;
			`,
			expectedErr: "",
		},
		{
			name: "Invalid member access",
			input: `
				let x: int = (impl (interface { let val: type = int; }) { let val: int = 42; }).val;
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
			expectedErr: "in let expected variable of type <int>, got <string>",
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
			expectedErr: "in let expected variable of type <int>, got <type(<int>)>",
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

func TestSwitchExpression(t *testing.T) {
	tests := []TestCase{
		{
			name: "Valid switch",
			input: `
				let res: type = int;
				let final: int = switch res {
					case int { return 1; }
					case string { return 2; }
					default { return 0; }
				};
			`,
			expectedErr: "",
		},
		{
			name: "Invalid switch (type mismatch between val and case)",
			input: `
				let res: int = 5;
				let final: int = switch res {
					case int {
						// len expects string or array, not int
						return len(res);
					}
					case string {
						return len(res);
					}
					default {
						return 0;
					}
				};
			`,
			expectedErr: "case pattern type",
		},
		{
			name: "Invalid switch (no return)",
			input: `
				let res: int = 5;
				let final: int = switch res {
					case 1 { 
					} 
					case 2 { 
					}
					default { 
					}
				};
			`,
			expectedErr: "switch must return some value in case blocks",
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
					case int { 
						return res + 1; 
					}
					case string { 
						return len(res); 
					}
					default { 
						return 0; 
					}
				};
			`,
			expectedErr: "",
		},
		{
			name: "Invalid operation in narrowed context (catches type error correctly)",
			input: `
				let res: int | string = 5;
				let final: int = switch res.(type) {
					case int { 
						// len expects string or array, not int 
						return len(res); 
					} 
					case string { 
						return len(res); 
					}
					default { 
						return 0; 
					}
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
					case OkType { 
						return res.value; 
					} 
					case ErrType { 
						return 0; 
					}
					default { 
						return 0; 
					}
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
				
				let _: int = switch res.(type) {
					case OkType { 
						// у OkType нет поля error 
						return res.error; 
					} 
					case ErrType { 
						return 0; 
					}
					default { 
						return 0; 
					}
				};
			`,
			expectedErr: "type error",
		},
	}
	runTypeCheckerTests(t, tests)
}
