package evaluator

import "funlang/object"

var builtins = map[string]*object.Builtin{
	"len": {Fn: func(args ...object.Object) object.Object {
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int(len(arg.Value))}
		default:
			return newError("argument to `len_str` not supported, got %s",
				object.LookUpObjSignature(args[0].Type()))
		}
	},
	},
}
