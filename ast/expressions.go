package ast

import "funlang/token"

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (intLit *IntegerLiteral) expressionNode()      {}
func (intLit *IntegerLiteral) TokenLiteral() string { return intLit.Token.Literal }
func (intLit *IntegerLiteral) String() string       { return intLit.Token.Literal }

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    ExpressionNode
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
