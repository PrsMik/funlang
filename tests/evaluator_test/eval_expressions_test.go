package evaluator_test

import (
	"funlang/object"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return -10;", -10},
		{"return 5 + 5 + 5 + 5 - 10;", 10},
		{"return 2 * 2 * 2 * 2 * 2;", 32},
		{"return -50 + 100 + -50;", 0},
		{"return 5 * 2 + 10;", 20},
		{"return 5 + 2 * 10;", 25},
		{"return 20 + 2 * -10;", 0},
		{"return 50 / 2 * 2 + 10;", 60},
		{"return 2 * (5 + 10);", 30},
		{"return 3 * 3 * 3 + 10;", 37},
		{"return 3 * (3 * 3) + 10;", 37},
		{"return (5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		n        int
		input    string
		expected bool
	}{
		{1, "return true;", true},
		{2, "return false;", false},
		{3, "return 1 < 2;", true},
		{4, "return 1 > 2;", false},
		{5, "return 1 < 1;", false},
		{6, "return 1 > 1;", false},
		{7, "return 1 == 1;", true},
		{8, "return 1 != 1;", false},
		{9, "return 1 == 2;", false},
		{10, "return 1 != 2;", true},
		{11, "return true == true;", true},
		{12, "return false == false;", true},
		{13, "return true == false;", false},
		{14, "return true != false;", true},
		{15, "return false != true;", true},
		{16, "return (1 < 2) == true;", true},
		{17, "return (1 < 2) == false;", false},
		{18, "return true && true;", true},
		{19, "return true && false;", false},
		{20, "return false || true;", true},
		{21, "return false || false;", false},
		{22, "return (1 < 2) && (3 > 2);", true},
	}
	for _, tt := range tests {
		testBooleanObject(t, tt.n, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		n        int
		input    string
		expected bool
	}{
		{1, "return !true;", false},
		{2, "return !false;", true},
		{3, "return !!true;", true},
		{4, "return !!false;", false},
	}
	for _, tt := range tests {
		testBooleanObject(t, tt.n, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return "Hello World!";`, "Hello World!"},
		{`return "Hello" + " " + "World!";`, "Hello World!"},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}
		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=%q, want=%q", str.Value, tt.expected)
		}
	}
}

func TestEvalIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return if (1 > 2) { return 10; } else { return 20; };", 20},
		{"return if (1 < 2) { return 10; } else { return 20; };", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalArrayExpressions(t *testing.T) {
	input := "return[1, 2 * 2, 3 + 3];"
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestEvalArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return [1, 2, 3][0];", 1},
		{"return[1, 2, 3][1];", 2},
		{"return[1, 2, 3][2];", 3},
		{"let i: int = 0; return [1][i];", 1},
		{"return [1, 2, 3][1 + 1];", 3},
		{"let myArray: [int] = [1, 2, 3]; return myArray[2];", 3},
		{"let myArray: [int] =[1, 2, 3]; return myArray[0] + myArray[1] + myArray[2];", 6},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalHashExpressions(t *testing.T) {
	input := `let two: string = "two";
			  return {
				"one": 10 - 9,
				two: 1 + 1,
				"thr" + "ee": 6 / 2 
			  };`
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestEvalHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`return {"foo": 5}["foo"];`, 5},
		{`return {"foo": 5}["bar"];`, nil},
		{`let key: string = "foo"; return {"foo": 5}[key];`, 5},
		{`return {5: 5}[5];`, 5},
		{`return {true: 5}[true];`, 5},
		{`return {false: 5}[false];`, 5},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}
