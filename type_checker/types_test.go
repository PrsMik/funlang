package type_checker

import (
	"funlang/lexer"
	"funlang/parser"
	"funlang/types"
	"testing"
)

func TestTypesEquals(t *testing.T) {
	tests := []struct {
		name string
		a    types.Type
		b    types.Type
		want bool
	}{
		{"1", &types.IntType{}, &types.IntType{}, true},
		{"2", &types.IntType{}, &types.BoolType{}, false},
		{"3", &types.IntType{}, &types.FuncType{}, false},
		{"4", &types.FuncType{ParamsTypes: []types.Type{&types.IntType{}, &types.IntType{}},
			ReturnType: &types.BoolType{}}, &types.FuncType{}, false},
		{"5", &types.FuncType{ParamsTypes: []types.Type{&types.IntType{}, &types.IntType{}},
			ReturnType: &types.BoolType{}},
			&types.FuncType{ParamsTypes: []types.Type{&types.IntType{}, &types.IntType{}},
				ReturnType: &types.BoolType{}}, true},
	}
	for _, tt := range tests {
		if got := types.Equals(tt.a, tt.b); got != tt.want {
			t.Fatalf("Equals() = %v, want %v", got, tt.want)
		}
	}
}

func TestCheckLetStatement(t *testing.T) {
	test := []struct {
		input string
		tp    types.Type
		want  bool
	}{
		{"let x: int = 5;", &types.IntType{}, true},
		{`let x: string = "hello";`, &types.StringType{}, true},
		{`let x: string = "hello" + "world!";`, &types.StringType{}, true},
		{`let x: [string] = ["hello", "world!"];`, &types.ArrayType{}, true},
		{`let x: [string] = [];`, &types.ArrayType{}, true},
		{`let x: {string : int} = {"hello" : 1};`, &types.HashMapType{}, true},
		{`let x: {string : int} = {};`, &types.HashMapType{}, true},
		{`let x: string = ["hello"][0];`, &types.StringType{}, true},
		{"let y: bool = true;", &types.BoolType{}, true},
		{"let z: int = true;", &types.BoolType{}, false},
		{"let z: int = -5;", &types.IntType{}, true},
		{"let z: bool = -5;", &types.BoolType{}, false},
		{"let z: int = -true;", &types.IntType{}, false},
		{"let z: int = 5 + 5 * 5;", &types.IntType{}, true},
		{"let z: bool = false && true || !false;", &types.IntType{}, true},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		env := types.NewTypeEviroment()

		chk := New(env)
		chk.CheckProgram(prg)
		if len(chk.errors) != 0 && tt.want != false {
			checkCheckerErrors(t, chk)
		}
	}
}

func TestCheckReturnSatement(t *testing.T) {
	test := []struct {
		input string
		want  bool
	}{
		{"let x: int = if (2 > 1) { return 5; } else { return 2; };", true},
		{"let x: [int] = if (2 > 1) { return []; } else { return [2, 3]; };", true},
		{"let x: {int : int} = if (2 > 1) { return {}; } else { return {2 : 3}; };", true},
		{"let x: int = if (2 > 1) { return 5 + 5 * 5; } else { return 1; };", true},
		{"let x: bool = if (2 > 1) { return true && false || false; } else { return true; }; ", true},
		{"let x: bool = if (2 > 1) { let x: bool = true && false || false; return x; } else { return true; }; ", true},
		{"let x: bool = if (2 + 1) { return true && false || false; } else { return true; }; ", false},
		{"let x: bool = if (2) { return true && false || false; } else { return true; }; ", false},
		{"let x: bool = if (true) { return true && false || false; } else { return true; }; ", true},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		env := types.NewTypeEviroment()

		chk := New(env)
		chk.CheckProgram(prg)
		if len(chk.errors) != 0 && tt.want != false {
			checkCheckerErrors(t, chk)
		}
	}

}

func TestCheckIdentifier(t *testing.T) {
	test := []struct {
		input string
		want  bool
	}{
		{"let x: int = 5; let y: int = x + 3;", true},
		{"let x: bool = true; let y: bool = !x;", true},
		{`let y: [int] = [1, 2, 3]; let x: [int] = y;`, true},
		{`let y: [int] = [1, 2, 3]; let x: int = y[0];`, true},
		{`let y: {int : int} = {1 : 2, 2 : 3, 3 : 4}; let x: {int : int} = y;`, true},
		{`let y: fn() -> [int] = fn() { return []; }; let x: [int] = y();`, true},
		{`let y: fn() -> int = fn() { return 1; }; 
		let x: [fn() -> int] = [y, y, y];`, true},
		{"let x: int = x;", false},
		{"let x: bool = true; let z: int = 5; let y: int = if (z > 1) { let x: int = 1; return x + z; } else { return 1; };", true},
		{"let y: int = x + 3;", false},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		env := types.NewTypeEviroment()

		chk := New(env)
		chk.CheckProgram(prg)
		if len(chk.errors) != 0 && tt.want != false {
			checkCheckerErrors(t, chk)
		}
	}
}

func TestCheckFunctionLiteral(t *testing.T) {
	test := []struct {
		input string
		want  bool
	}{
		{"let x: fn(int, int) -> int = fn(x, y) { return x + y; };", true},
		{"let x: fn(string, string) -> string = fn(x, y) { return x + y; };", true},
		{"let x: fn(string, string) -> string = fn(x, y) { return x / y; };", false},
		{"let x: fn(string, int) -> string = fn(x, y) { return x + y; };", false},
		{"let x: fn(int, int) -> bool = fn(x, y) { return x + y; };", false},
		{"let x: fn(int, int) -> bool = fn() { return x + y; };", false},
		{"let x: fn(bool, int) -> bool = fn(x, y) { return x + y; };", false},
		{`let x: fn(bool, bool, int) -> int = fn(x, y, z) { let n: int = z + if (x && y) { return 1; } else { return 0; };
		return z; };`, true},
		{"let add: fn(int, int) -> int = fn(x, y) { return add(x, y); };", true},
		{"let add: fn(int, int) -> int = fn(x, y) { return add(true, y); };", false},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		env := types.NewTypeEviroment()

		chk := New(env)
		chk.CheckProgram(prg)
		if len(chk.errors) != 0 && tt.want != false {
			checkCheckerErrors(t, chk)
		}
	}
}

func TestCheckCallExpression(t *testing.T) {
	test := []struct {
		input string
		want  bool
	}{
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = add(1, 2);", true},
		{"let add: fn(int, int) -> int = fn(x, y) { return x + y; }; let z: int = 1 + add(1, 2);", true},
		{"let add: fn(int, int) -> bool = fn(x, y) { return true; }; let z: int = add(1, 2);", false},
		{"let add: fn(int, int) -> bool = fn(x, y) { return false; }; let z: bool = add(true, 2);", false},
		{"let add: fn(int, int) -> bool = fn(x, y) { return false; }; let z: bool = add(add(1, 1), 2);", false},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		env := types.NewTypeEviroment()

		chk := New(env)
		chk.CheckProgram(prg)
		if len(chk.errors) != 0 && tt.want != false {
			checkCheckerErrors(t, chk)
		}
	}
}

func checkParserErrors(t *testing.T, prs *parser.Parser) {
	errors := prs.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func checkCheckerErrors(t *testing.T, chk *TypeChecker) {
	errors := chk.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Checker has %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("checker error: %q", msg)
	}
	t.FailNow()
}
