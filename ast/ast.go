package ast

type Node interface {
	TokenLiteral() string
	String() string
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
