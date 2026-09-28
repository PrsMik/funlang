package types

import "funlang/ast"

type Info struct {
	GlobalScope   *TypeEviroment
	TypesInfo     map[ast.Node]Type
	TypeNodes     map[ast.Node]bool
	Definitions   map[ast.Node]ast.Node
	Scopes        map[ast.Node]*TypeEviroment
	ExpectedTypes map[ast.Node]Type
}

func NewInfo(env *TypeEviroment) *Info {
	return &Info{
		GlobalScope:   env,
		TypesInfo:     make(map[ast.Node]Type),
		TypeNodes:     make(map[ast.Node]bool),
		Definitions:   make(map[ast.Node]ast.Node),
		Scopes:        make(map[ast.Node]*TypeEviroment),
		ExpectedTypes: make(map[ast.Node]Type),
	}
}
