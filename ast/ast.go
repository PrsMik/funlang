package ast

import "funlang/token"

type Node interface {
	TokenLiteral() string
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

func (prog *Program) TokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].TokenLiteral()
	}
	return ""
}

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type SimpleType struct {
	Token token.Token
	Value string
}

func (i *SimpleType) typeNode()            {}
func (t *SimpleType) TokenLiteral() string { return t.Token.Literal }

// инструкция let - токен инструкции (Token), идентификатор для корого let (Name),
// тип идентификатора (Type) и выражение справа от "=" (Value)
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Type  TypeNode
	Value ExpressionNode
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
