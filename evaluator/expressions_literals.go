package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
	"sort"
)

func evalDynamicFunctionType(node *ast.FunctionLiteral, env *object.Environment,
	info *types.Info, resolvedType types.Type) types.Type {

	var dynamicFuncType types.Type = resolvedType
	if statFuncType, ok := resolvedType.(*types.FuncType); ok {
		dynFuncType := &types.FuncType{
			Params:     make([]types.FuncParam, len(statFuncType.Params)),
			ReturnType: statFuncType.ReturnType,
		}
		copy(dynFuncType.Params, statFuncType.Params)

		if node.ReturnType != nil {
			rtVal := Eval(node.ReturnType, env, info)
			if tObj, isType := rtVal.(*object.TypeObject); isType {
				dynFuncType.ReturnType = tObj.TypeValue
			}
		}

		for i, pTypeNode := range node.ParamTypes {
			if pTypeNode != nil {
				ptVal := Eval(pTypeNode, env, info)
				if tObj, isType := ptVal.(*object.TypeObject); isType {
					dynFuncType.Params[i].Type = tObj.TypeValue
				}
			}
		}

		dynamicFuncType = dynFuncType
	}
	return dynamicFuncType
}

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

func evalTypeTypeIdentifier(node *ast.TypeType, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	} else if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("type identifier not found: %s", node.Value)
}

func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment, info *types.Info) object.Object {
	elements := evalExpressions(node.Elements, env, info)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}

	return &object.Array{
		Elements: elements,
		TypeObj:  &object.TypeObject{TypeValue: info.TypesInfo[node], Env: env}}
}

func evalHashMapLiteral(node *ast.HashMapLiteral, env *object.Environment, info *types.Info) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	// var firstKey, firstValue object.Object

	for _, keyNode := range sortedPairKeys(node.Pairs) {
		if _, isBad := keyNode.(*ast.BadExpression); isBad {
			continue
		}

		valNode := node.Pairs[keyNode]

		key := Eval(keyNode, env, info)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			if key == nil {
				return newError("unusable as hash key: %s", keyNode.String())
			}
			return newError("unusable as hash key: %s", key.Inspect())
		}

		value := Eval(valNode, env, info)
		if isError(value) {
			return value
		}

		// if firstKey == nil {
		// 	firstKey, firstValue = key, value
		// }

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashMap{
		Pairs:   pairs,
		TypeObj: &object.TypeObject{TypeValue: info.TypesInfo[node], Env: env}}
}

// ключи литерала map в порядке записи
func sortedPairKeys(pairs map[ast.ExpressionNode]ast.ExpressionNode) []ast.ExpressionNode {
	keys := make([]ast.ExpressionNode, 0, len(pairs))
	for key := range pairs {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		left, right := keys[i].Start(), keys[j].Start()
		if left.Line != right.Line {
			return left.Line < right.Line
		}
		return left.Column < right.Column
	})

	return keys
}
