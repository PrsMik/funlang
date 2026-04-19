package types_test

import (
	"funlang/types"
	"strings"
	"testing"
)

func getBuiltinFunc(t *testing.T, name string) *types.BuiltinFunc {
	env := types.NewTypeEviroment()
	sym, ok := env.Get(name)
	if !ok {
		t.Fatalf("builtin %s not found in env", name)
	}
	fn, ok := sym.SymbolType.(*types.BuiltinFunc)
	if !ok {
		t.Fatalf("builtin %s is not a BuiltinFunc type", name)
	}
	return fn
}

func TestBuiltinLen(t *testing.T) {
	lenFn := getBuiltinFunc(t, "len")

	tests := []struct {
		name        string
		args        []types.Type
		wantSuccess bool
	}{
		{"Valid: String", []types.Type{&types.StringType{}}, true},
		{"Valid: Array", []types.Type{&types.ArrayType{ElementsType: &types.IntType{}}}, true},
		{"Invalid: Int", []types.Type{&types.IntType{}}, false},
		{"Invalid: arg count", []types.Type{&types.StringType{}, &types.StringType{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := lenFn.CheckFunc(tt.args)
			if tt.wantSuccess {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if _, ok := res.(*types.IntType); !ok {
					t.Errorf("len should return IntType, got %T", res)
				}
			} else if err == nil {
				t.Errorf("expected error, got success")
			}
		})
	}
}

func TestBuiltinTail(t *testing.T) {
	tailFn := getBuiltinFunc(t, "tail")

	tests := []struct {
		name        string
		args        []types.Type
		wantSuccess bool
	}{
		{"Valid: Array of Int", []types.Type{&types.ArrayType{ElementsType: &types.IntType{}}}, true},
		{"Invalid: String", []types.Type{&types.StringType{}}, false},
		{"Invalid: arg count", []types.Type{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tailFn.CheckFunc(tt.args)
			if tt.wantSuccess {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if arr, ok := res.(*types.ArrayType); !ok || !types.Equals(arr.ElementsType, &types.IntType{}) {
					t.Errorf("tail should return matching ArrayType, got %v", res.Signature())
				}
			} else if err == nil {
				t.Errorf("expected error, got success")
			}
		})
	}
}

func TestBuiltinPush(t *testing.T) {
	pushFn := getBuiltinFunc(t, "push")

	tests := []struct {
		name        string
		args        []types.Type
		wantSuccess bool
		errorMsg    string
	}{
		{
			"Valid: push Int to [Int]", []types.Type{&types.ArrayType{ElementsType: &types.IntType{}}, &types.IntType{}},
			true, "",
		},
		{
			"Valid: push to empty Array", []types.Type{&types.ArrayType{}, &types.IntType{}},
			true, "",
		},
		{
			"Invalid: Type mismatch", []types.Type{&types.ArrayType{ElementsType: &types.IntType{}}, &types.StringType{}},
			false, "push cannot be performed",
		},
		{
			"Invalid: First arg not array", []types.Type{&types.IntType{}, &types.IntType{}},
			false, "push does not support type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := pushFn.CheckFunc(tt.args)
			if tt.wantSuccess {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if arr, ok := res.(*types.ArrayType); !ok || !types.Equals(arr.ElementsType, tt.args[1]) {
					t.Errorf("push should return matching ArrayType, got %v", res.Signature())
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got success")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestBuiltinPuts(t *testing.T) {
	putsFn := getBuiltinFunc(t, "puts")

	tests := []struct {
		name        string
		args        []types.Type
		wantSuccess bool
	}{
		{"Valid: Printable Int", []types.Type{&types.IntType{}}, true},
		{"Valid: Printable String", []types.Type{&types.StringType{}}, true},
		{"Invalid: Non-printable Func", []types.Type{&types.FuncType{Params: []types.FuncParam{},
			ReturnType: &types.IntType{}}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := putsFn.CheckFunc(tt.args)
			if tt.wantSuccess {
				if err != nil {
					t.Errorf("unexpected error for puts: %v", err)
				}
			} else if err == nil {
				t.Errorf("expected error for non-printable type")
			}
		})
	}
}
