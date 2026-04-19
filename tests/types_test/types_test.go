package types_test

import (
	"funlang/types"
	"testing"
)

var (
	_ types.Type          = (*types.IntType)(nil)
	_ types.HashableType  = (*types.IntType)(nil)
	_ types.PrintableType = (*types.IntType)(nil)

	_ types.Type          = (*types.StringType)(nil)
	_ types.HashableType  = (*types.StringType)(nil)
	_ types.PrintableType = (*types.StringType)(nil)

	_ types.Type          = (*types.ArrayType)(nil)
	_ types.PrintableType = (*types.ArrayType)(nil)
)

func TestTypeSignatures(t *testing.T) {
	tests := []struct {
		name     string
		tp       types.Type
		expected string
	}{
		{"Illegal", &types.IllegalType{}, "<none>"},
		{"Null", &types.NullType{}, "<null>"},
		{"Int", &types.IntType{}, "<int>"},
		{"Bool", &types.BoolType{}, "<bool>"},
		{"String", &types.StringType{}, "<string>"},
		{"Empty Array", &types.ArrayType{}, "<[]>"},
		{"Int Array", &types.ArrayType{ElementsType: &types.IntType{}}, "<[<int>]>"},
		{"Empty HashMap", &types.HashMapType{}, "<{}>"},
		{"String:Int HashMap", &types.HashMapType{KeyType: &types.StringType{}, ElementType: &types.IntType{}}, "<{<string>:<int>}>"},
		{"Builtin Func", &types.BuiltinFunc{}, "<builtin_func>"},
		{"Empty Func", &types.FuncType{ReturnType: &types.IntType{}}, "<fn() -> <int>>"},
		{"Func with params", &types.FuncType{
			Params: []types.FuncParam{
				{Name: "x", Type: &types.IntType{}},
				{Name: "y", Type: &types.StringType{}},
			},
			ReturnType: &types.BoolType{},
		}, "<fn(x: <int>, y: <string>) -> <bool>>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.Signature(); got != tt.expected {
				t.Errorf("Signature() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name string
		a    types.Type
		b    types.Type
		want bool
	}{
		// Базовые типы
		{"Int == Int", &types.IntType{}, &types.IntType{}, true},
		{"Int != Bool", &types.IntType{}, &types.BoolType{}, false},
		{"String == String", &types.StringType{}, &types.StringType{}, true},

		// Массивы
		{"Array[Int] == Array[Int]", &types.ArrayType{ElementsType: &types.IntType{}}, &types.ArrayType{ElementsType: &types.IntType{}}, true},
		{"Array[Int] != Array[Bool]", &types.ArrayType{ElementsType: &types.IntType{}}, &types.ArrayType{ElementsType: &types.BoolType{}}, false},
		{"Array[empty] matches Array[Int]", &types.ArrayType{}, &types.ArrayType{ElementsType: &types.IntType{}}, true},
		{"Array[Int] matches Array[empty]", &types.ArrayType{ElementsType: &types.IntType{}}, &types.ArrayType{}, true},
		{"Array != Int", &types.ArrayType{}, &types.IntType{}, false},

		// Хэшмапы
		{"HashMap == HashMap", &types.HashMapType{KeyType: &types.StringType{}, ElementType: &types.IntType{}}, &types.HashMapType{KeyType: &types.StringType{}, ElementType: &types.IntType{}}, true},
		{"HashMap[empty] matches typed HashMap", &types.HashMapType{}, &types.HashMapType{KeyType: &types.StringType{}, ElementType: &types.IntType{}}, true},
		{"HashMap != Array", &types.HashMapType{}, &types.ArrayType{}, false},

		// Функции
		{"Func == Func (same params)",
			&types.FuncType{Params: []types.FuncParam{{Name: "x", Type: &types.IntType{}}}, ReturnType: &types.BoolType{}},
			&types.FuncType{Params: []types.FuncParam{{Name: "y", Type: &types.IntType{}}}, ReturnType: &types.BoolType{}},
			true,
		},
		{"Func != Func (diff param types)",
			&types.FuncType{Params: []types.FuncParam{{Name: "x", Type: &types.IntType{}}}, ReturnType: &types.BoolType{}},
			&types.FuncType{Params: []types.FuncParam{{Name: "y", Type: &types.StringType{}}}, ReturnType: &types.BoolType{}},
			false,
		},
		{"Func != Func (diff return types)",
			&types.FuncType{Params: []types.FuncParam{{Name: "x", Type: &types.IntType{}}}, ReturnType: &types.BoolType{}},
			&types.FuncType{Params: []types.FuncParam{{Name: "y", Type: &types.IntType{}}}, ReturnType: &types.IntType{}},
			false,
		},
		{"Func != Func (diff params count)",
			&types.FuncType{Params: []types.FuncParam{}, ReturnType: &types.BoolType{}},
			&types.FuncType{Params: []types.FuncParam{{Name: "y", Type: &types.IntType{}}}, ReturnType: &types.BoolType{}},
			false,
		},
		{"Func != Int", &types.FuncType{}, &types.IntType{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.Equals(tt.a, tt.b); got != tt.want {
				t.Errorf("Equals(%s, %s) = %v, want %v", tt.a.Signature(), tt.b.Signature(), got, tt.want)
			}
		})
	}
}
