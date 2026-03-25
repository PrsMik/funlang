package evaluator

import (
	"fmt"
	"funlang/object"
)

var builtins = map[string]*object.Builtin{
	"len": {Fn: func(args ...object.Object) object.Object {
		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: len(arg.Elements)}
		case *object.String:
			return &object.Integer{Value: len(arg.Value)}
		default:
			return newError("argument to `len` not supported, got %s",
				object.LookUpObjSignature(args[0].Type()))
		}
	},
	},
	"tail": {Fn: func(args ...object.Object) object.Object {
		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		if length > 0 {
			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[1:length])
			return &object.Array{Elements: newElements}
		}

		return NULL
	},
	},
	"push": {Fn: func(args ...object.Object) object.Object {
		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		newElements := make([]object.Object, length+1)
		copy(newElements[1:length+1], arr.Elements)
		newElements[0] = args[1]

		return &object.Array{Elements: newElements}
	},
	},
	"puts": {Fn: func(args ...object.Object) object.Object {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return NULL
	},
	},
}
