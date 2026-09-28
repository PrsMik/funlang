package evaluator

import (
	"fmt"
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalMemberAccessExpression(node *ast.MemberAccessExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)
	fieldName := node.Property.Value

	switch leftObj := left.(type) {
	case *object.TypeObject:
		return evalTypeField(leftObj.TypeValue, fieldName, env, info)

	case *object.ImplInstance:
		val, ok := leftObj.Fields[fieldName]
		if !ok {
			return &object.Error{Message: fmt.Sprintf("runtime error. instance has no field '%s'", fieldName)}
		}
		return val
	}
	return &object.Error{Message: "runtime error. left operand is not a type or instance"}
}

func evalTypeField(typeVal types.Type,
	fieldName string,
	env *object.Environment,
	info *types.Info) object.Object {

	switch targetType := typeVal.(type) {
	case *types.InterfaceType:
		if fieldType, ok := targetType.Fields[fieldName]; ok {
			return &object.TypeObject{TypeValue: fieldType}
		}
	case *types.ImplType:
		if defaultExpr, ok := targetType.FieldsDefaultValues[fieldName]; ok {
			return Eval(defaultExpr, env, info)
		}
	case *types.IntersectionType:
		leftRes := evalTypeField(targetType.Left, fieldName, env, info)
		if _, isErr := leftRes.(*object.Error); !isErr {
			return leftRes
		}

		return evalTypeField(targetType.Right, fieldName, env, info)
	}
	return &object.Error{Message: "runtime error. not a type with fields"}
}

func evalImplInstantiationExpression(node *ast.ImplInstantiationExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)
	switch leftObj := left.(type) {
	case *object.TypeObject:
		fields, ok := node.Fields.(*ast.HashMapLiteral)
		if !ok {
			return &object.Error{Message: "runtime error. fields is not hash map literal"}
		}

		return evalFieldsInstantation(leftObj, fields, env, info)
	}
	return &object.Error{Message: "runtime error. object is not impl"}
}

func evalFieldsInstantation(implTypeObject *object.TypeObject,
	fields *ast.HashMapLiteral,
	env *object.Environment,
	info *types.Info) object.Object {

	implType, ok := implTypeObject.TypeValue.(*types.ImplType)
	if !ok {
		return &object.Error{Message: "runtime error. can only instantiate 'impl' types"}
	}

	instanceObject := &object.ImplInstance{TypeObj: implTypeObject, Fields: make(map[string]object.Object)}

	for key, elemExpr := range implType.FieldsDefaultValues {
		val := Eval(elemExpr, env, info)
		if isError(val) {
			return val
		}
		instanceObject.Fields[key] = val
	}

	for keyNode, elemExpr := range fields.Pairs {
		keyIndent, ok := keyNode.(*ast.Identifier)
		if !ok {
			return &object.Error{Message: fmt.Sprintf("runtime error. field key %s is not an identifier",
				keyIndent.String())}
		}

		elemVal := Eval(elemExpr, env, info)

		if isError(elemVal) {
			return elemVal
		}

		instanceObject.Fields[keyIndent.Value] = elemVal
	}

	return instanceObject
}

func evalTypeAccessExpression(node *ast.TypeAccessExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)

	tpVal := object.GetTypeFromObject(left)
	if _, isIllegal := tpVal.(*types.IllegalType); isIllegal {
		return &object.Error{Message: "runtime error. error when trying to extract type"}
	}

	return &object.TypeObject{TypeValue: tpVal}
}
