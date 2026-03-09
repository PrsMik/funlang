package object

type Environment struct {
	store map[string]Object
}

func NewEnviroment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func (env *Environment) Get(name string) (Object, bool) {
	val, ok := env.store[name]
	return val, ok
}

func (env *Environment) Set(name string, val Object) {
	env.store[name] = val
}
