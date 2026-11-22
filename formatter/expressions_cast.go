package formatter

import "funlang/ast"

func (fmtr *Formatter) formatSafeTypeCastExpression(node *ast.SafeTypeCastExpression) {
	fmtr.formatExpression(node.Left)
	fmtr.out.WriteString(".(")
	fmtr.formatType(node.TargetType)
	fmtr.out.WriteString(")")
	fmtr.prevEndLine = node.SemiToken.End.Line
}

func (fmtr *Formatter) formatUnsafeTypeCastExpression(node *ast.UnsafeTypeCastExpression) {
	fmtr.out.WriteString(".")
	fmtr.formatType(node.TargetType)
	fmtr.out.WriteString("(")
	fmtr.formatExpression(node.Value)
	fmtr.out.WriteString(")")
	fmtr.prevEndLine = node.SemiToken.End.Line
}
