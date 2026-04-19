package types

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
)

func getMapWithBuiltins() map[string]SymbolInfo {
	builtins := make(map[string]SymbolInfo)

	builtins["len"] = SymbolInfo{SymbolType: &BuiltinFunc{
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
		ReturnType: &IntType{},
	}, DeclNode: getFakeNode("len", "len")}

	var argType Type
	builtins["tail"] = SymbolInfo{SymbolType: &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 1 {
				return &IllegalType{}, fmt.Errorf("tail expects 1 argument but got %d", len(args))
			}

			switch tempArgType := args[0].(type) {
			case *ArrayType:
				argType = tempArgType.ElementsType
				return &ArrayType{ElementsType: tempArgType.ElementsType}, nil
			default:
				return &IllegalType{}, fmt.Errorf("tail does not support type %s", tempArgType.Signature())
			}
		},
	}, DeclNode: getFakeNode("tail", "tail")}
	builtins["tail"].SymbolType.(*BuiltinFunc).ReturnType = &ArrayType{ElementsType: argType}

	builtins["push"] = SymbolInfo{SymbolType: &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 2 {
				return &IllegalType{}, fmt.Errorf("push expects 2 argument but got %d", len(args))
			}

			firstArg, ok := args[0].(*ArrayType)
			if !ok {
				return &IllegalType{}, fmt.Errorf("push does not support type %s", args[0].Signature())
			}

			if !Equals(firstArg.ElementsType, args[1]) && firstArg.ElementsType != nil {
				return &IllegalType{}, fmt.Errorf("push cannot be performed on %s with 2nd arg as %s",
					args[0].Signature(), args[1].Signature())
			}

			argType = args[1]

			return &ArrayType{ElementsType: args[1]}, nil
		},
	}, DeclNode: getFakeNode("push", "push")}
	builtins["push"].SymbolType.(*BuiltinFunc).ReturnType = &ArrayType{ElementsType: argType}

	builtins["puts"] = SymbolInfo{SymbolType: &BuiltinFunc{
		CheckFunc: func(args []Type) (Type, error) {
			if len(args) != 1 {
				return &IllegalType{}, fmt.Errorf("puts expects 1 argument but got %d", len(args))
			}

			_, ok := args[0].(PrintableType)
			if !ok {
				return &IllegalType{}, fmt.Errorf("puts does not support type %s", args[0].Signature())
			}

			// if !Equals(firstArg.ElementsType, args[1]) {
			// 	return &IllegalType{}, fmt.Errorf("push cannot be performed on %s with 2nd arg as %s",
			// 		args[0].Signature(), args[1].Signature())
			// }

			return &IntType{}, nil
		},
		ReturnType: &IntType{},
	}, DeclNode: getFakeNode("puts", "puts")}

	return builtins
}

func getFakeNode(literal string, value string) ast.Node {
	fakeToken := token.Token{Literal: literal,
		Start: token.Position{Column: -1, Line: -1},
		End:   token.Position{Column: -1, Line: -1}}
	return &ast.Identifier{Token: fakeToken, Value: value}
}
