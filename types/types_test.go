package types

import (
	"funlang/lexer"
	"funlang/parser"
	"testing"
)

func TestTypesEquals(t *testing.T) {
	tests := []struct {
		name string
		a    Type
		b    Type
		want bool
	}{
		{"1", &IntType{}, &IntType{}, true},
		{"2", &IntType{}, &BoolType{}, false},
		{"3", &IntType{}, &FuncType{}, false},
		{"4", &FuncType{ParamsTypes: []Type{&IntType{}, &IntType{}}, ReturnType: &BoolType{}}, &FuncType{}, false},
		{"5", &FuncType{ParamsTypes: []Type{&IntType{}, &IntType{}}, ReturnType: &BoolType{}},
			&FuncType{ParamsTypes: []Type{&IntType{}, &IntType{}}, ReturnType: &BoolType{}}, true},
	}
	for _, tt := range tests {
		if got := Equals(tt.a, tt.b); got != tt.want {
			t.Fatalf("Equals() = %v, want %v", got, tt.want)
		}
	}
}

func TestCheckLetStatement(t *testing.T) {
	test := []struct {
		input string
		tp    Type
		want  bool
	}{
		{"let x: int = 5;", &IntType{}, true},
		{"let y: bool = true;", &BoolType{}, true},
		{"let z: int = true;", &BoolType{}, false},
		{"let z: int = -5;", &IntType{}, true},
		{"let z: bool = -5;", &BoolType{}, false},
		{"let z: int = -true;", &IntType{}, false},
		{"let z: int = 5 + 5 * 5;", &IntType{}, true},
		{"let z: bool = false && true || !false;", &IntType{}, true},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		chk := New(nil)
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
		{"let x: int = if (2 > 1) { return 5; };", true},
		{"let x: int = if (2 > 1) { return 5 + 5 * 5; } else { return 1; };", true},
		{"let x: bool = if (2 > 1) { return true && false || false; }; ", true},
		{"let x: bool = if (2 > 1) { let x: bool = true && false || false; return x; }; ", false},
		{"let x: bool = if (2 + 1) { return true && false || false; }; ", false},
		{"let x: bool = if (2) { return true && false || false; }; ", false},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		chk := New(nil)
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
		{"let x: bool = true; let z: int = 5; let y: int = if (z > 1) { let x: int = 1; return x + z; } else { return 1; };", true},
		{"let y: int = x + 3;", false},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		chk := New(nil)
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
		{"let x: fn(int, int) -> bool = fn(x, y) { return x + y; };", false},
		{"let x: fn(int, int) -> bool = fn() { return x + y; };", false},
		{"let x: fn(bool, int) -> bool = fn(x, y) { return x + y; };", false},
		{`let x: fn(bool, bool, int) -> int = fn(x, y, z) { let n: int = z + if (x && y) { return 1; } else { return 0; }; 
		return z; };`, true},
	}
	for _, tt := range test {
		lxr := lexer.New(tt.input)
		prs := parser.New(lxr)

		prg := prs.ParseProgram()
		checkParserErrors(t, prs)

		chk := New(nil)
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
