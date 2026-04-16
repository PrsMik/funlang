package evaluator_test

import (
	"funlang/evaluator"
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"testing"
)

func TestRuntimeErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let x: int = 10 / 0;",
			"runtime error: integer divide by zero",
		},
		{
			"let arr: [int] = [1, 2]; let x: int = arr[5];",
			"runtime error: array index 5 is out of bounds (max is 1)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}
