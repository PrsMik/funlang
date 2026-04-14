package formatter

import (
	"funlang/ast"
)

func (fmtr *Formatter) forceNewLine(stmt ast.StatementNode) {
	if fmtr.prevEndLine >= 0 {
		delta := stmt.Start().Line - fmtr.prevEndLine
		if delta > 1 {
			fmtr.out.WriteString("\n\n")
		} else {
			fmtr.out.WriteString("\n")
		}
		fmtr.writeIndent()
	}
	fmtr.prevEndLine = stmt.Start().Line
}

func (fmtr *Formatter) formatStatement(stmt ast.StatementNode) {
	fmtr.printTrailingComments(stmt.Start().Line)

	fmtr.pushEndLine(stmt.End().Line)
	defer fmtr.popEndLine()

	switch node := stmt.(type) {
	case *ast.LetStatement:
		fmtr.formatLetStatement(node)
	case *ast.ReturnStatement:
		fmtr.formatReturnStatement(node)
	case *ast.BlockStatement:
		fmtr.formatBlockStatement(node)
	}

	fmtr.writeInlineComment(stmt)
}

func (fmtr *Formatter) formatLetStatement(node *ast.LetStatement) {
	fmtr.forceNewLine(node)
	fmtr.out.WriteString("let ")
	fmtr.out.WriteString(node.Name.Value)
	fmtr.out.WriteString(": ")
	fmtr.formatType(node.Type)
	fmtr.out.WriteString(" = ")

	fmtr.prevEndLine = node.Value.Start().Line
	fmtr.formatExpression(node.Value)

	fmtr.out.WriteString(";")
	fmtr.prevEndLine = node.SemiToken.End.Line
}

func (fmtr *Formatter) formatReturnStatement(node *ast.ReturnStatement) {
	fmtr.forceNewLine(node)
	fmtr.out.WriteString("return ")

	fmtr.prevEndLine = node.Value.Start().Line
	fmtr.formatExpression(node.Value)

	fmtr.out.WriteString(";")
	fmtr.prevEndLine = node.SemiToken.End.Line
}

func (fmtr *Formatter) formatBlockStatement(n *ast.BlockStatement) {
	fmtr.out.WriteString("{")
	fmtr.prevEndLine = n.Token.End.Line

	if len(n.Statements) > 0 {
		fmtr.indentLevel++
		for _, s := range n.Statements {
			fmtr.formatStatement(s)
		}
		fmtr.indentLevel--

		fmtr.printTrailingComments(n.End().Line)

		delta := n.End().Line - fmtr.prevEndLine
		if delta > 1 {
			fmtr.out.WriteString("\n\n")
		} else {
			fmtr.out.WriteString("\n")
		}
		fmtr.writeIndent()
	} else {
		fmtr.printTrailingComments(n.End().Line)
	}

	fmtr.out.WriteString("}")
	fmtr.prevEndLine = n.End().Line
}

func (fmtr *Formatter) writeInlineComment(stmt ast.StatementNode) {
	if fmtr.commentIndex < len(fmtr.comments) {
		nextComment := fmtr.comments[fmtr.commentIndex]
		if nextComment.Start.Line == stmt.End().Line {
			fmtr.out.WriteString("\x00")
			fmtr.out.WriteString(FormatCommentText(nextComment))
			fmtr.commentIndex++
		}
	}
}
