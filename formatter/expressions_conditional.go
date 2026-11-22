package formatter

import "funlang/ast"

func (fmtr *Formatter) formatIfExpression(n *ast.IfExpression) {
	fmtr.out.WriteString("if ")
	fmtr.out.WriteString("")
	fmtr.formatExpression(n.Condition)
	fmtr.out.WriteString(" ")
	fmtr.formatStatement(n.Consequence)

	if n.Alternative != nil {
		fmtr.out.WriteString(" else ")
		fmtr.formatStatement(n.Alternative)
	}
}

func (fmtr *Formatter) formatSwitchExpression(node *ast.SwitchExpression) {
	fmtr.out.WriteString("switch ")
	fmtr.formatExpression(node.Value)
	fmtr.out.WriteString(" {")
	fmtr.indentLevel++
	for _, caseBlock := range node.Cases {
		fmtr.out.WriteString("\n")
		fmtr.writeIndent()
		fmtr.out.WriteString("case ")
		fmtr.prevEndLine = caseBlock.End().Line
		fmtr.formatExpression(caseBlock.Pattern)
		fmtr.out.WriteString(" ")
		fmtr.formatBlockStatement(caseBlock.Body)
	}

	fmtr.out.WriteString("\n")
	fmtr.writeIndent()
	fmtr.out.WriteString("default ")
	fmtr.formatBlockStatement(node.Default)
	fmtr.indentLevel--
	fmtr.out.WriteString("\n")
	fmtr.out.WriteString("}")
}
