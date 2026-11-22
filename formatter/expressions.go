package formatter

import (
	"funlang/ast"
)

func (fmtr *Formatter) formatExpression(expr ast.ExpressionNode) {
	fmtr.advanceTo(expr.Start().Line)

	switch node := expr.(type) {
	case *ast.IntegerLiteral, *ast.BooleanLiteral, *ast.StringLiteral:
		fmtr.out.WriteString(node.String())
		fmtr.prevEndLine = node.End().Line

	case *ast.ArrayLiteral:
		fmtr.formatArrayLiteral(node)
	case *ast.HashMapLiteral:
		fmtr.formatHashMapLiteral(node)
	case *ast.FunctionLiteral:
		fmtr.formatFunctionLiteral(node)
	case *ast.InfixExpression:
		fmtr.formatInfixExpression(node)
	case *ast.PrefixExpression:
		fmtr.formatPrefixExpression(node)
	case *ast.GroupedExpression:
		fmtr.formatGroupedExpression(node)
	case *ast.IfExpression:
		fmtr.formatIfExpression(node)
	case *ast.SwitchExpression:
		fmtr.formatSwitchExpression(node)
	case *ast.CallExpression:
		fmtr.formatCallExpression(node)
	case *ast.IndexExpression:
		fmtr.formatIndexExpression(node)
	case *ast.ComptimeExpression:
		fmtr.formatComptimeExpression(node)
	case *ast.InterfaceLiteral:
		fmtr.formatInterfaceLiteral(node)
	case *ast.ImplLiteral:
		fmtr.formatImplLiteral(node)
	case *ast.SafeTypeCastExpression:
		fmtr.formatSafeTypeCastExpression(node)
	case *ast.UnsafeTypeCastExpression:
		fmtr.formatUnsafeTypeCastExpression(node)
	default:
		fmtr.out.WriteString(node.String())
		fmtr.prevEndLine = node.End().Line
	}
}

func (fmtr *Formatter) formatInfixExpression(n *ast.InfixExpression) {
	fmtr.formatExpression(n.Left)

	opLine := n.Token.Start.Line
	leftEndLine := n.Left.End().Line

	if opLine > leftEndLine {
		fmtr.indentLevel++
		fmtr.advanceTo(opLine)
		fmtr.out.WriteString(n.Operator)
		fmtr.indentLevel--
	} else {
		fmtr.out.WriteString(" ")
		fmtr.out.WriteString(n.Operator)
	}

	rightStartLine := n.Right.Start().Line
	if rightStartLine > opLine {
		fmtr.indentLevel++
		fmtr.advanceTo(rightStartLine)
		fmtr.formatExpression(n.Right)
		fmtr.indentLevel--
	} else {
		fmtr.out.WriteString(" ")
		fmtr.prevEndLine = rightStartLine
		fmtr.formatExpression(n.Right)
	}
}

func (fmtr *Formatter) formatPrefixExpression(n *ast.PrefixExpression) {
	fmtr.out.WriteString(n.Operator)
	fmtr.formatExpression(n.Right)
}

func (fmtr *Formatter) formatGroupedExpression(n *ast.GroupedExpression) {
	fmtr.out.WriteString("(")
	fmtr.formatExpression(n.Value)
	fmtr.out.WriteString(")")
	fmtr.prevEndLine = n.SemiToken.End.Line
}

func (fmtr *Formatter) formatComptimeExpression(n *ast.ComptimeExpression) {
	fmtr.out.WriteString("comptime ")
	fmtr.formatExpression(n.Value)
}
