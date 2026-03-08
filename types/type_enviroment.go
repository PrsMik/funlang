package types

type TypeEviroment struct {
	types map[string]Type
	Outer *TypeEviroment
}

func NewTypeEviroment() *TypeEviroment {
	return &TypeEviroment{types: make(map[string]Type)}
}

func NewEnclosedTypeEviroment(outer *TypeEviroment) *TypeEviroment {
	return &TypeEviroment{types: make(map[string]Type), Outer: outer}
}

func (env *TypeEviroment) Get(name string) (Type, bool) {
	tp, ok := env.types[name]

	if !ok && env.Outer != nil {
		return env.Outer.Get(name)
	}

	return tp, ok
}

func (env *TypeEviroment) Set(name string, tp Type) {
	env.types[name] = tp
}
