package evaluator_test

import (
	"funlang/object"
	"testing"
)

func TestEvalBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`return len("");`, 0},
		{`return len("four");`, 4},
		{`return len("hello world");`, 11},
		{`return len([]);`, 0},
		{`return len([1, 2, 3]);`, 3},
		{`return tail([1, 2, 3]);`, []int{2, 3}},
		{`return tail([]);`, nil},
		{`return push([2, 3], 1);`, []int{1, 2, 3}},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case []int:
			testArrayIntegerObject(t, evaluated, expected)
		case nil:
			if _, ok := evaluated.(*object.Null); !ok {
				t.Errorf("test[%d] - object is not Null. got=%T (%+v)", i, evaluated, evaluated)
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("test[%d] - object is not Error. got=%T (%+v)", i, evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("test[%d] - wrong error message. expected=%q, got=%q", i, expected, errObj.Message)
			}
		}
	}
}
