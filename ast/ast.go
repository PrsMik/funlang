package ast

import (
	"funlang/token"
)

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

type SimpleType struct {
	Token token.Token
	Value string
}

func (t *SimpleType) typeNode()            {}
func (t *SimpleType) TokenLiteral() string { return t.Token.Literal }
func (t *SimpleType) String() string       { return t.Value }

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

// инструкция return - токен инструкции (Token), выражение которое возвращает return (ReturnValue),
type ReturnStatement struct {
	Token token.Token
	Value ExpressionNode
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
