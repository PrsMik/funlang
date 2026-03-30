package types

import "funlang/ast"

type SymbolInfo struct {
	SymbolType Type
	DeclNode   ast.Node
}

type TypeEviroment struct {
	types map[string]SymbolInfo
	Outer *TypeEviroment
}

func NewTypeEviroment() *TypeEviroment {
	// types: make(map[string]Type)
	env := &TypeEviroment{}
	env.types = getMapWithBuiltins()
	return env
}

func NewEnclosedTypeEviroment(outer *TypeEviroment) *TypeEviroment {
	return &TypeEviroment{types: make(map[string]SymbolInfo), Outer: outer}
}

func (env *TypeEviroment) Get(name string) (SymbolInfo, bool) {
	tp, ok := env.types[name]

	if !ok && env.Outer != nil {
		return env.Outer.Get(name)
	}

	return tp, ok
}

func (env *TypeEviroment) Set(name string, tp Type, declNode ast.Node) {
	env.types[name] = SymbolInfo{SymbolType: tp, DeclNode: declNode}
}

func (env *TypeEviroment) GetAllNames() []string {
	namesMap := make(map[string]bool)
	currEnv := env

	for currEnv != nil {
		for name := range currEnv.types {
			namesMap[name] = true
		}
		currEnv = currEnv.Outer
	}

	res := make([]string, 0, len(namesMap))
	for name := range namesMap {
		res = append(res, name)
	}
	return res
}
