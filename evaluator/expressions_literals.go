package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/type_checker"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	} else if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

func evalTypeIdentifier(node *ast.SimpleType, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	} else if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("type identifier not found: %s", node.Value)
}

func evalHashMapLiteral(node *ast.HashMapLiteral, env *object.Environment, chk *type_checker.TypeChecker) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env, chk)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Inspect())
		}

		value := Eval(valNode, env, chk)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashMap{Pairs: pairs}
}
