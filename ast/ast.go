package ast

import "funlang/token"

type Node interface {
	TokenLiteral() string
	String() string
	Start() token.Position
	End() token.Position
}

// узел является инструкцией
type StatementNode interface {
	Node
	statementNode()
}

// узел является выражением
type ExpressionNode interface {
	Node
	expressionNode()
}

// узел является типом
type TypeNode interface {
	Node
	typeNode()
}

// корневой узел
type Program struct {
	Statements []StatementNode
}

func (prg *Program) TokenLiteral() string {
	if len(prg.Statements) > 0 {
		return prg.Statements[0].TokenLiteral()
	}
	return ""
}

func (prg *Program) Start() token.Position { return prg.Statements[0].Start() }
func (prg *Program) End() token.Position   { return prg.Statements[len(prg.Statements)-1].End() }
