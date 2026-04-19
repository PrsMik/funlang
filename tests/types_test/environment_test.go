package types_test

import (
	"funlang/types"
	"testing"
)

func TestTypeEnvironment(t *testing.T) {
	globalEnv := types.NewTypeEviroment()

	_, ok := globalEnv.Get("len")
	if !ok {
		t.Fatalf("expected builtin 'len' to be present in global environment")
	}

	globalEnv.Set("x", &types.IntType{}, nil)

	sym, ok := globalEnv.Get("x")
	if !ok {
		t.Fatalf("expected to find 'x' in environment")
	}
	if _, isInt := sym.SymbolType.(*types.IntType); !isInt {
		t.Errorf("expected 'x' to be IntType, got %T", sym.SymbolType)
	}

	innerEnv := types.NewEnclosedTypeEviroment(globalEnv)

	if _, ok := innerEnv.Get("x"); !ok {
		t.Errorf("expected inner environment to resolve 'x' from outer scope")
	}
	if _, ok := innerEnv.Get("len"); !ok {
		t.Errorf("expected inner environment to resolve 'len' from outer scope")
	}

	innerEnv.Set("x", &types.StringType{}, nil)
	innerSym, _ := innerEnv.Get("x")
	if _, isStr := innerSym.SymbolType.(*types.StringType); !isStr {
		t.Errorf("expected local 'x' to shadow outer and be StringType, got %T", innerSym.SymbolType)
	}

	globalSym, _ := globalEnv.Get("x")
	if _, isInt := globalSym.SymbolType.(*types.IntType); !isInt {
		t.Errorf("expected global 'x' to remain IntType, got %T", globalSym.SymbolType)
	}
}

func TestGetAllNames(t *testing.T) {
	globalEnv := types.NewTypeEviroment()
	innerEnv := types.NewEnclosedTypeEviroment(globalEnv)

	globalEnv.Set("globalVar", &types.IntType{}, nil)
	innerEnv.Set("localVar", &types.BoolType{}, nil)
	innerEnv.Set("len", &types.StringType{}, nil)

	names := innerEnv.GetAllNames()

	expectedNames := map[string]bool{
		"globalVar": true,
		"localVar":  true,
		"len":       true,
		"tail":      true,
		"push":      true,
		"puts":      true,
	}

	if len(names) != len(expectedNames) {
		t.Fatalf("expected %d names, got %d", len(expectedNames), len(names))
	}

	for _, name := range names {
		if !expectedNames[name] {
			t.Errorf("found unexpected name %q in GetAllNames", name)
		}
	}
}
