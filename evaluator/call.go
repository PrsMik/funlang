package evaluator

import (
	"funlang/object"
	"funlang/type_checker"
)

func applyFunction(fn object.Object, args []object.Object, chk *type_checker.TypeChecker) object.Object {
	for {
		switch f := fn.(type) {
		case *object.Function:
			extendedEnv := extendFunctionEnv(f, args)
			evaluated := Eval(f.Body, extendedEnv, chk)

			res := unwrapReturnValue(evaluated)

			if tc, ok := res.(*object.TailCall); ok {
				fn = tc.Function
				args = tc.Arguments
				continue
			}

			return res

		case *object.Builtin:
			return f.Fn(args...)

		default:
			return newError("not a function: %s", fn.Inspect())
		}
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
