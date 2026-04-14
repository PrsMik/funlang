package formatter

import (
	"funlang/ast"
	"sort"
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
	case *ast.IfExpression:
		fmtr.formatIfExpression(node)
	case *ast.CallExpression:
		fmtr.formatCallExpression(node)
	case *ast.IndexExpression:
		fmtr.formatIndexExpression(node)
	default:
		fmtr.out.WriteString(node.String())
		fmtr.prevEndLine = node.End().Line
	}
}

func (fmtr *Formatter) formatArrayLiteral(n *ast.ArrayLiteral) {
	fmtr.out.WriteString("[")
	isMulti := n.Token.Start.Line != n.End().Line

	fmtr.formatSequence(
		len(n.Elements),
		isMulti,
		func(i int) int { return n.Elements[i].Start().Line },
		func(i int) int { return n.Elements[i].End().Line },
		func(i int) { fmtr.formatExpression(n.Elements[i]) },
	)

	fmtr.out.WriteString("]")
	fmtr.prevEndLine = n.End().Line
}

func (fmtr *Formatter) formatHashMapLiteral(n *ast.HashMapLiteral) {
	fmtr.out.WriteString("{")
	isMulti := n.Token.Start.Line != n.End().Line

	pairs := getSortedHashMapPairs(n)

	fmtr.formatSequence(
		len(pairs),
		isMulti,
		func(i int) int { return pairs[i].k.Start().Line },
		func(i int) int { return pairs[i].v.End().Line },
		func(i int) {
			fmtr.formatExpression(pairs[i].k)
			fmtr.out.WriteString(": ")
			fmtr.prevEndLine = pairs[i].v.Start().Line
			fmtr.formatExpression(pairs[i].v)
		},
	)

	fmtr.out.WriteString("}")
	fmtr.prevEndLine = n.End().Line
}

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
		if i < len(n.ParamTypes) && n.ParamTypes[i] != nil {
			fmtr.out.WriteString(": ")
			fmtr.formatType(n.ParamTypes[i])
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

func (fmtr *Formatter) formatIfExpression(n *ast.IfExpression) {
	fmtr.out.WriteString("if ")
	fmtr.out.WriteString("(")
	fmtr.formatExpression(n.Condition)
	fmtr.out.WriteString(") ")
	fmtr.formatStatement(n.Consequence)

	if n.Alternative != nil {
		fmtr.out.WriteString(" else ")
		fmtr.formatStatement(n.Alternative)
	}
}

func (fmtr *Formatter) formatPrefixExpression(n *ast.PrefixExpression) {
	fmtr.out.WriteString(n.Operator)
	fmtr.formatExpression(n.Right)
}

func (fmtr *Formatter) formatIndexExpression(n *ast.IndexExpression) {
	fmtr.formatExpression(n.Left)
	fmtr.out.WriteString("[")
	fmtr.formatExpression(n.Index)
	fmtr.out.WriteString("]")
}

type kvPair struct {
	k ast.ExpressionNode
	v ast.ExpressionNode
}

func getSortedHashMapPairs(n *ast.HashMapLiteral) []kvPair {
	var pairs []kvPair
	for k, v := range n.Pairs {
		if _, ok := k.(*ast.VirtualNode); ok {
			continue
		}
		pairs = append(pairs, kvPair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].k.Start().Line == pairs[j].k.Start().Line {
			return pairs[i].k.Start().Column < pairs[j].k.Start().Column
		}
		return pairs[i].k.Start().Line < pairs[j].k.Start().Line
	})

	return pairs
}
