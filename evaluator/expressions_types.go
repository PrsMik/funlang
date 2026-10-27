package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalArrayType(node *ast.ArrayType, env *object.Environment, info *types.Info) object.Object {
	elem := Eval(node.ElementsType, env, info)
	if isError(elem) {
		return elem
	}
	if tObj, ok := elem.(*object.TypeObject); ok {
		return &object.TypeObject{TypeValue: &types.ArrayType{ElementsType: tObj.TypeValue}, Env: env}
	}
	return &object.Error{Message: "runtime error: not a type in array type definition"}
}

func evalHashMapType(node *ast.HashMapType, env *object.Environment, info *types.Info) object.Object {
	key := Eval(node.KeyType, env, info)
	if isError(key) {
		return key
	}
	elem := Eval(node.ElementType, env, info)
	if isError(elem) {
		return elem
	}

	kObj, kOk := key.(*object.TypeObject)
	eObj, eOk := elem.(*object.TypeObject)
	if kOk && eOk {
		return &object.TypeObject{
			TypeValue: &types.HashMapType{
				KeyType:     kObj.TypeValue,
				ElementType: eObj.TypeValue,
			},
			Env: env}
	}
	return &object.Error{Message: "runtime error: not a type in hash map type definition"}
}

func evalFunctionType(node *ast.FunctionType, env *object.Environment, info *types.Info) object.Object {
	params := make([]types.FuncParam, len(node.ParamsTypes))
	for i, pTypeNode := range node.ParamsTypes {
		pObj := Eval(pTypeNode, env, info)
		if isError(pObj) {
			return pObj
		}
		if tObj, ok := pObj.(*object.TypeObject); ok {
			params[i] = types.FuncParam{Type: tObj.TypeValue}
			if i < len(node.Parameters) && node.Parameters[i] != nil {
				params[i].Name = node.Parameters[i].Value
			}
		} else {
			return &object.Error{Message: "runtime error: parameter is not a type"}
		}
	}

	retObj := Eval(node.ReturnType, env, info)
	if isError(retObj) {
		return retObj
	}

	if tObj, ok := retObj.(*object.TypeObject); ok {
		return &object.TypeObject{TypeValue: &types.FuncType{Params: params, ReturnType: tObj.TypeValue}, Env: env}
	}

	return &object.Error{Message: "runtime error: return type is not a type"}
}
