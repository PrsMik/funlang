package types

import "fmt"

func getMapWithBuiltins() map[string]Type {
	builtins := make(map[string]Type)
	builtins["len"] = &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 1 {
				return &IllegalType{}, fmt.Errorf("len expects 1 argument but got %d", len(args))
			}

			switch argType := args[0].(type) {
			case *StringType:
				return &IntType{}, nil
			case *ArrayType:
				return &IntType{}, nil
			default:
				return &IllegalType{}, fmt.Errorf("len does not support type %s", argType.Signature())
			}
		},
	}
	return builtins
}
