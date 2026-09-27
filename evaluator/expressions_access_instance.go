package evaluator

import (
	"fmt"
	"funlang/ast"
	"funlang/object"
	"funlang/type_checker"
	"funlang/types"
)

func evalMemberAccessExpression(node *ast.MemberAccessExpression, env *object.Environment,
	chk *type_checker.TypeChecker) object.Object {

	left := Eval(node.Left, env, chk)
	switch leftObj := left.(type) {
	case *object.TypeObject:
		tp, ok := chk.GetFieldType(leftObj.TypeValue, node.Property.Value, node)
		if !ok {
			return &object.Error{Message: "runtime error. no field in type"}
		}
		return evalTypeField(leftObj.TypeValue, node.Property.Value, tp, node.Left, env, chk)
	case *object.ImplInstance:
		_, ok := chk.GetFieldType(leftObj.TypeObj.TypeValue, node.Property.Value, node)
		if !ok {
			return &object.Error{Message: "runtime error. no field in type"}
		}
		return leftObj.Fields[node.Property.Value]
	}
	return &object.Error{Message: "runtime error. object is not a type"}
}

func evalTypeField(typeVal types.Type,
	fieldName string,
	fieldType types.Type,
	node ast.ExpressionNode,
	env *object.Environment,
	chk *type_checker.TypeChecker) object.Object {

	switch targetType := typeVal.(type) {
	case *types.InterfaceType:
		return &object.TypeObject{TypeValue: fieldType}
	case *types.ImplType:
		return Eval(targetType.FieldsDefaultValues[fieldName], env, chk)
	case *types.IntersectionType:
		if fieldType, ok := chk.GetFieldType(targetType.Left, fieldName, node); ok {
			return evalTypeField(targetType.Left, fieldName, fieldType, node, env, chk)
		}

		if fieldType, ok := chk.GetFieldType(targetType.Right, fieldName, node); ok {
			return evalTypeField(targetType.Right, fieldName, fieldType, node, env, chk)
		}
	}
	return &object.Error{Message: "runtime error. not a type with fields"}
}

func evalImplInstantiationExpression(node *ast.ImplInstantiationExpression, env *object.Environment,
	chk *type_checker.TypeChecker) object.Object {

	left := Eval(node.Left, env, chk)
	switch leftObj := left.(type) {
	case *object.TypeObject:
		fields, ok := node.Fields.(*ast.HashMapLiteral)
		if !ok {
			return &object.Error{Message: "runtime error. fields is not hash map literal"}
		}

		return evalFieldsInstantation(leftObj, fields, node, env, chk)
	}
	return &object.Error{Message: "runtime error. object is not impl"}
}

func evalFieldsInstantation(implTypeObject *object.TypeObject,
	fields *ast.HashMapLiteral,
	node ast.ExpressionNode,
	env *object.Environment,
	chk *type_checker.TypeChecker) object.Object {

	implType, ok := implTypeObject.TypeValue.(*types.ImplType)
	if !ok {
		return &object.Error{Message: "runtime error. type is not impl"}
	}

	instanceObject := &object.ImplInstance{TypeObj: implTypeObject, Fields: make(map[string]object.Object)}

	for key, elem := range implType.FieldsDefaultValues {
		instanceObject.Fields[key] = Eval(elem, env, chk)
	}

	for key, elem := range fields.Pairs {
		keyIndent, ok := key.(*ast.Identifier)
		if !ok {
			return &object.Error{Message: fmt.Sprintf("%s is not an identifier", keyIndent.String())}
		}

		elemVal := Eval(elem, env, chk)

		instanceObject.Fields[keyIndent.Value] = elemVal
	}

	return instanceObject
}
