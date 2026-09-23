package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/type_checker"
	"funlang/types"
)

func evalInterfaceLiteral(node *ast.InterfaceLiteral, env *object.Environment, chk *type_checker.TypeChecker) object.Object {
	ifaceType, ok := chk.CheckInterfaceLiteral(node).(*types.TypeType)
	// TODO: перепроверить, что здесь всё безопасно для рантайма 1
	if !ok {
		return &object.Error{Message: "runtime error. interface is not a type. in theory this should never happen :/"}
	}
	return &object.TypeObject{TypeValue: ifaceType.Underlying}
}

func evalImplLiteral(node *ast.ImplLiteral, env *object.Environment, chk *type_checker.TypeChecker) object.Object {
	implType, ok := chk.CheckImplLiteral(node).(*types.TypeType)
	// TODO: перепроверить, что здесь всё безопасно для рантайма 2
	if !ok {
		return &object.Error{Message: "runtime error. impl is not a type. in theory this should never happen :/"}
	}
	return &object.TypeObject{TypeValue: implType.Underlying}
}
