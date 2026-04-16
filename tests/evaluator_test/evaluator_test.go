package evaluator_test

import (
	"funlang/evaluator"
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"funlang/type_checker"
	"funlang/types"
	"testing"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	type_env := types.NewTypeEviroment()
	chk := type_checker.New(type_env, nil)
	chk.CheckProgram(program)
	checkCheckerErrors(t, chk)

	env := object.NewEnvironment()
	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testArrayIntegerObject(t *testing.T, obj object.Object, expected []int) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}
	if len(result.Elements) != len(expected) {
		t.Errorf("object has wrong size. got=%d, want=%d", len(result.Elements), expected)
		return false
	}
	for i, val := range expected {
		testIntegerObject(t, result.Elements[i], val)
	}
	return true
}

func testBooleanObject(t *testing.T, n int, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Test %d. object is not Boolean. got=%T (%+v)", n, obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Test %d. object has wrong value. got=%t, want=%t", n, result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, prs *parser.Parser) {
	errors := prs.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err.Msg)
	}
	t.FailNow()
}

func checkCheckerErrors(t *testing.T, chk *type_checker.TypeChecker) {
	errors := chk.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Checker has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("checker error: %q", msg)
	}
	t.FailNow()
}
