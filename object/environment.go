package object

import "funlang/types"

type Environment struct {
	store map[string]Object
	outer *Environment
}

var (
	MetaType = &TypeObject{TypeValue: types.MetaType}
)

func NewEnvironment() *Environment {
	env := &Environment{store: make(map[string]Object), outer: nil}
	env.Set("type", MetaType)
	env.Set("int", &TypeObject{TypeValue: &types.IntType{}})
	env.Set("string", &TypeObject{TypeValue: &types.StringType{}})
	env.Set("bool", &TypeObject{TypeValue: &types.BoolType{}})
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
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
