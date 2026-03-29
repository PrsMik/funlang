package type_checker

import "funlang/ast"

func (chk *TypeChecker) typeError(msg string, node ast.Node) {
	chk.errors = append(chk.errors, TypeError{Msg: "type error: " + msg, Node: node})
}
