package formatter

import (
	"funlang/ast"
	"sort"
)

func (fmtr *Formatter) formatArrayLiteral(n *ast.ArrayLiteral) {
	if n.TypeExpr != nil {
		fmtr.out.WriteString("[]")
		// fmtr.out.WriteString(n.TypeExpr.String())
		fmtr.formatExpression(n.TypeExpr)
		fmtr.out.WriteString("{")
	} else {
		fmtr.out.WriteString("[")
	}

	isMulti := n.Token.Start.Line != n.End().Line

	fmtr.formatSequence(
		len(n.Elements),
		isMulti,
		func(i int) int { return n.Elements[i].Start().Line },
		func(i int) int { return n.Elements[i].End().Line },
		func(i int) { fmtr.formatExpression(n.Elements[i]) },
	)

	if n.TypeExpr != nil {
		fmtr.out.WriteString("}")
	} else {
		fmtr.out.WriteString("]")
	}
	fmtr.prevEndLine = n.End().Line
}

func (fmtr *Formatter) formatHashMapLiteral(n *ast.HashMapLiteral) {
	fmtr.out.WriteString(n.TypeExpr.String())

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
		if _, ok := k.(*ast.BadExpression); ok {
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
