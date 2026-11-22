package formatter

import "funlang/ast"

func (fmtr *Formatter) formatInterfaceLiteral(n *ast.InterfaceLiteral) {
	fmtr.out.WriteString("interface ")
	fmtr.formatStatement(n.Body)
}

func (fmtr *Formatter) formatImplLiteral(n *ast.ImplLiteral) {
	fmtr.out.WriteString("impl ")
	if n.TargetType != nil {
		fmtr.formatExpression(n.TargetType)
	}
	fmtr.formatStatement(n.Body)
}
