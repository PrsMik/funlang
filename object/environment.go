package object

import (
	"funlang/ast"
	"funlang/types"
)

type Environment struct {
	store          map[string]Object
	types          map[string]types.Type
	outer          *Environment
	ComptimeValues map[ast.Node]Object
}

var (
	MetaType = &TypeObject{TypeValue: types.MetaType}
)

func NewEnvironment() *Environment {
	env := &Environment{
		store:          make(map[string]Object),
		types:          make(map[string]types.Type),
		outer:          nil,
		ComptimeValues: make(map[ast.Node]Object)}
	env.Set("type", MetaType)
	env.Set("int", &TypeObject{TypeValue: &types.IntType{}})
	env.Set("string", &TypeObject{TypeValue: &types.StringType{}})
	env.Set("bool", &TypeObject{TypeValue: &types.BoolType{}})
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	env.ComptimeValues = outer.ComptimeValues
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	val, ok := env.store[name]
	if !ok && env.outer != nil {
		val, ok = env.outer.Get(name)
	}
	return val, ok
}

func (env *Environment) Set(name string, val Object) {
	env.store[name] = val
}
