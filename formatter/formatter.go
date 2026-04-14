package formatter

import (
	"bytes"
	"funlang/ast"
	"funlang/token"
)

type Formatter struct {
	out            bytes.Buffer
	indentLevel    int
	comments       []token.Token
	commentIndex   int
	prevEndLine    int
	activeEndLines []int
}

func New(out bytes.Buffer) *Formatter {
	return &Formatter{out: out, indentLevel: 0, comments: make([]token.Token, 0), commentIndex: 0, prevEndLine: -1}
}

func (fmtr *Formatter) FormatProgram(prog *ast.Program) string {
	fmtr.comments = prog.Comments

	for _, stmt := range prog.Statements {
		fmtr.formatStatement(stmt)
	}

	fmtr.printTrailingComments(999999)

	raw := fmtr.out.String()

	return alignComments(raw)
}
