package types

func getMapWithBuiltins() map[string]Type {
	builtins := make(map[string]Type)
	builtins["len_str"] = &FuncType{ParamsTypes: []Type{&StringType{}}, ReturnType: &IntType{}}
	return builtins
}
