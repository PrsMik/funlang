package formatter

import "funlang/ast"

func (fmtr *Formatter) formatCallExpression(n *ast.CallExpression) {
	fmtr.formatExpression(n.Function)
	fmtr.out.WriteString("(")
	isMulti := n.Token.Start.Line != n.End().Line

	fmtr.formatSequence(
		len(n.Arguments),
		isMulti,
		func(i int) int { return n.Arguments[i].Start().Line },
		func(i int) int { return n.Arguments[i].End().Line },
		func(i int) { fmtr.formatExpression(n.Arguments[i]) },
	)

	fmtr.out.WriteString(")")
	fmtr.prevEndLine = n.End().Line
}

func (fmtr *Formatter) formatFunctionLiteral(n *ast.FunctionLiteral) {
	fmtr.out.WriteString("fn(")
	for i, param := range n.Parameters {
		fmtr.out.WriteString(param.Value)
		if i < len(n.ParamsTypes) && n.ParamsTypes[i] != nil {
			fmtr.out.WriteString(": ")
			fmtr.formatType(n.ParamsTypes[i])
		}
		if i < len(n.Parameters)-1 {
			fmtr.out.WriteString(", ")
		}
	}
	fmtr.out.WriteString(") ")

	if n.ReturnType != nil {
		fmtr.out.WriteString("-> ")
		fmtr.formatType(n.ReturnType)
		fmtr.out.WriteString(" ")
	}
	fmtr.formatBlockStatement(n.Body)
	fmtr.prevEndLine = n.End().Line
}
