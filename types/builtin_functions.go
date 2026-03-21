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
	builtins["tail"] = &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 1 {
				return &IllegalType{}, fmt.Errorf("tail expects 1 argument but got %d", len(args))
			}

			switch argType := args[0].(type) {
			case *ArrayType:
				return &IntType{}, nil
			default:
				return &IllegalType{}, fmt.Errorf("tail does not support type %s", argType.Signature())
			}
		},
	}
	builtins["push"] = &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 2 {
				return &IllegalType{}, fmt.Errorf("push expects 2 argument but got %d", len(args))
			}

			firstArg, ok := args[0].(*ArrayType)
			if !ok {
				return &IllegalType{}, fmt.Errorf("push does not support type %s", args[0].Signature())
			}

			if !Equals(firstArg.ElementsType, args[1]) {
				return &IllegalType{}, fmt.Errorf("push cannot be performed on %s with 2nd arg as %s",
					args[0].Signature(), args[1].Signature())
			}

			return &ArrayType{ElementsType: args[1]}, nil
		},
	}
	return builtins
}
