package evaluator

import "funlang/object"

var builtins = map[string]*object.Builtin{
	"len_str": {Fn: func(args ...object.Object) object.Object {
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].Inspect())
		}
	},
	},
}
