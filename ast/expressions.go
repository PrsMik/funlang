package ast

import "funlang/token"

// литерал типа int (например "5")
type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (intLit *IntegerLiteral) expressionNode()      {}
func (intLit *IntegerLiteral) TokenLiteral() string { return intLit.Token.Literal }

// литерал типа bool (например "true")
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (boolLit *BooleanLiteral) expressionNode()      {}
func (boolLit *BooleanLiteral) TokenLiteral() string { return boolLit.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (funLit *FunctionLiteral) expressionNode()      {}
func (funLit *FunctionLiteral) TokenLiteral() string { return funLit.Token.Literal }

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

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

type IfExpression struct {
	Token       token.Token
	Condition   ExpressionNode
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()      {}
func (ifExpr *IfExpression) TokenLiteral() string { return ifExpr.Token.Literal }
