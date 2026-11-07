package evaluator

import (
	"fmt"
	"funlang/ast"
	"funlang/object"
	"funlang/types"
	"sort"
)

func evalMemberAccessExpression(node *ast.MemberAccessExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)
	if isError(left) {
		return left
	}

	fieldName := node.Property.Value

	switch leftObj := left.(type) {
	case *object.TypeObject:
		return evalTypeField(leftObj, fieldName, env, info)

	case *object.ImplInstance:
		val, ok := leftObj.Fields[fieldName]
		if !ok {
			return &object.Error{Message: fmt.Sprintf("runtime error. instance has no field '%s'", fieldName)}
		}
		return val
	}
	return &object.Error{Message: "runtime error. left operand is not a type or instance"}
}

func implEnv(implType *types.ImplType, objEnv *object.Environment, fallback *object.Environment) *object.Environment {
	if typeEnv, ok := implType.Env.(*object.Environment); ok && typeEnv != nil {
		return typeEnv
	}
	if objEnv != nil {
		return objEnv
	}
	return fallback
}

func evalTypeField(typeObj *object.TypeObject,
	fieldName string,
	env *object.Environment,
	info *types.Info) object.Object {

	switch targetType := typeObj.TypeValue.(type) {
	case *types.InterfaceType:
		if fieldType, ok := targetType.Fields[fieldName]; ok {
			return &object.TypeObject{TypeValue: fieldType, Env: typeObj.Env}
		}
	case *types.ImplType:
		if defaultExpr, ok := targetType.FieldsDefaultValues[fieldName]; ok {
			return Eval(defaultExpr, implEnv(targetType, typeObj.Env, env), info)
		}
	case *types.IntersectionType:
		leftObj := &object.TypeObject{TypeValue: targetType.Left, Env: typeObj.Env}
		leftRes := evalTypeField(leftObj, fieldName, env, info)
		if _, isErr := leftRes.(*object.Error); !isErr {
			return leftRes
		}

		rightObj := &object.TypeObject{TypeValue: targetType.Right, Env: typeObj.Env}

		return evalTypeField(rightObj, fieldName, env, info)
	}
	return &object.Error{Message: "runtime error. not a type with fields"}
}

func evalImplInstantiationExpression(node *ast.ImplInstantiationExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)
	if isError(left) {
		return left
	}

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

	provided := make(map[string]ast.ExpressionNode, len(fields.Pairs))
	for keyNode, elemExpr := range fields.Pairs {
		if _, isBad := keyNode.(*ast.BadExpression); isBad {
			continue
		}

		keyIdent, isIdent := keyNode.(*ast.Identifier)
		if !isIdent {
			return &object.Error{Message: fmt.Sprintf("runtime error. field key %s is not an identifier",
				keyNode.String())}
		}

		provided[keyIdent.Value] = elemExpr
	}

	defaultsEnv := object.NewEnclosedEnvironment(implEnv(implType, implTypeObject.Env, env))

	for _, fieldName := range implFieldsOrder(implType) {
		var val object.Object
		if elemExpr, isProvided := provided[fieldName]; isProvided {
			val = Eval(elemExpr, env, info)
		} else {
			val = Eval(implType.FieldsDefaultValues[fieldName], defaultsEnv, info)
		}

		if isError(val) {
			return val
		}

		instanceObject.Fields[fieldName] = val
		defaultsEnv.Set(fieldName, val)
	}

	// если дополнительные поля каким-то образом прошли проверку чекера
	extraNames := make([]string, 0)
	for name := range provided {
		if _, exists := instanceObject.Fields[name]; !exists {
			extraNames = append(extraNames, name)
		}
	}
	sort.Strings(extraNames)

	for _, name := range extraNames {
		val := Eval(provided[name], env, info)
		if isError(val) {
			return val
		}
		instanceObject.Fields[name] = val
	}

	return instanceObject
}

func implFieldsOrder(implType *types.ImplType) []string {
	if len(implType.FieldsOrder) == len(implType.FieldsDefaultValues) {
		return implType.FieldsOrder
	}

	names := make([]string, 0, len(implType.FieldsDefaultValues))
	for name := range implType.FieldsDefaultValues {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func evalTypeAccessExpression(node *ast.TypeAccessExpression, env *object.Environment,
	info *types.Info) object.Object {

	left := Eval(node.Left, env, info)
	if isError(left) {
		return left
	}

	tpVal := object.GetTypeFromObject(left)
	if _, isIllegal := tpVal.(*types.IllegalType); isIllegal {
		return &object.Error{Message: "runtime error. error when trying to extract type"}
	}

	return &object.TypeObject{TypeValue: tpVal, Env: env}
}
