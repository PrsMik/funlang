package object_test

import (
	"funlang/ast"
	"funlang/object"
	"testing"
)

func TestObjectTypesAndInspection(t *testing.T) {
	tests := []struct {
		name         string
		obj          object.Object
		expectedType object.ObjectType
		expectedIns  string
	}{
		{
			name:         "Null",
			obj:          &object.Null{},
			expectedType: object.NULL_OBJ,
			expectedIns:  "null",
		},
		{
			name:         "Integer",
			obj:          &object.Integer{Value: 42},
			expectedType: object.INTEGER_OBJ,
			expectedIns:  "42",
		},
		{
			name:         "Boolean True",
			obj:          &object.Boolean{Value: true},
			expectedType: object.BOOLEAN_OBJ,
			expectedIns:  "true",
		},
		{
			name:         "Boolean False",
			obj:          &object.Boolean{Value: false},
			expectedType: object.BOOLEAN_OBJ,
			expectedIns:  "false",
		},
		{
			name:         "String",
			obj:          &object.String{Value: "hello world"},
			expectedType: object.STRING_OBJ,
			expectedIns:  "hello world",
		},
		{
			name: "Array",
			obj: &object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
			expectedType: object.ARRAY_OBJ,
			expectedIns:  "[1, 2, 3]",
		},
		{
			name: "HashMap",
			obj: &object.HashMap{Pairs: map[object.HashKey]object.HashPair{
				(&object.String{Value: "name"}).HashKey(): {
					Key:   &object.String{Value: "name"},
					Value: &object.String{Value: "funlang"},
				},
			}},
			expectedType: object.HASH_OBJ,
			expectedIns:  "{name: funlang}",
		},
		{
			name:         "ReturnValue",
			obj:          &object.ReturnValue{Value: &object.Integer{Value: 100}},
			expectedType: object.RETURN_VALUE_OBJ,
			expectedIns:  "100",
		},
		{
			name:         "TailCall",
			obj:          &object.TailCall{},
			expectedType: object.TAIL_CALL_OBJ,
			expectedIns:  "tail call",
		},
		{
			name:         "Builtin",
			obj:          &object.Builtin{},
			expectedType: object.BUILTIN_OBJ,
			expectedIns:  "builtin function",
		},
		{
			name:         "Error",
			obj:          &object.Error{Message: "type mismatch"},
			expectedType: object.ERROR_OBJ,
			expectedIns:  "RUNTIME ERROR: type mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.obj.Type() != tt.expectedType {
				t.Errorf("Type() = %v, want %v", tt.obj.Type(), tt.expectedType)
			}
			if tt.obj.Inspect() != tt.expectedIns {
				t.Errorf("Inspect() = %v, want %v", tt.obj.Inspect(), tt.expectedIns)
			}
		})
	}
}

func TestFunctionObject(t *testing.T) {
	fn := &object.Function{
		Parameters: []*ast.Identifier{
			{Value: "x"},
			{Value: "y"},
		},
		Body: &ast.BlockStatement{},
		Env:  object.NewEnvironment(),
	}

	if fn.Type() != object.FUNCTION_OBJ {
		t.Errorf("Type() = %v, want %v", fn.Type(), object.FUNCTION_OBJ)
	}

	expectedIns := "fn(x, y) {\n\n}"
	if fn.Inspect() != expectedIns {
		t.Errorf("Inspect() = %q, want %q", fn.Inspect(), expectedIns)
	}
}

func TestLookUpObjSignature(t *testing.T) {
	validTypes := []object.ObjectType{
		object.INTEGER_OBJ,
		object.BOOLEAN_OBJ,
		object.STRING_OBJ,
		object.NULL_OBJ,
		object.ARRAY_OBJ,
		object.HASH_OBJ,
	}

	for _, ot := range validTypes {
		sig := object.LookUpObjSignature(ot)
		if sig == "<unknown>" || sig == "" {
			t.Errorf("LookUpObjSignature(%d) returned %q, expected a valid type signature", ot, sig)
		}
	}

	unknownType := object.ObjectType(999)
	if sig := object.LookUpObjSignature(unknownType); sig != "<unknown>" {
		t.Errorf("LookUpObjSignature for unknown type = %q, want %q", sig, "<unknown>")
	}
}

func TestHashKeyUniquenessAcrossTypes(t *testing.T) {
	intOne := &object.Integer{Value: 1}
	boolTrue := &object.Boolean{Value: true}

	if intOne.HashKey() == boolTrue.HashKey() {
		t.Errorf("Integer(1) and Boolean(true) must have different HashKeys, but they were identical")
	}

	strOne := &object.String{Value: "\x01"}
	if intOne.HashKey() == strOne.HashKey() {
		t.Errorf("Integer(1) and String(byte 1) must have different HashKeys")
	}
}
