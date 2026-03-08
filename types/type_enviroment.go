package types

type TypeEviroment struct {
	types map[string]Type
	outer *TypeEviroment
}

func NewTypeEviroment() *TypeEviroment {
	return &TypeEviroment{types: make(map[string]Type)}
}

func NewEnclosedTypeEviroment(outer *TypeEviroment) *TypeEviroment {
	return &TypeEviroment{types: make(map[string]Type), outer: outer}
}

func (env *TypeEviroment) Get(name string) (Type, bool) {
	tp, ok := env.types[name]

	if !ok && env.outer != nil {
		return env.outer.Get(name)
	}

	return tp, ok
}

func (env *TypeEviroment) Set(name string, tp Type) {
	env.types[name] = tp
}
