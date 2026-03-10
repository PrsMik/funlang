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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (strLit *StringLiteral) expressionNode()      {}
func (strLit *StringLiteral) TokenLiteral() string { return strLit.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	ParamTypes []TypeNode
	ReturnType TypeNode
	Body       *BlockStatement
}

func (funcLit *FunctionLiteral) expressionNode()      {}
func (funcLit *FunctionLiteral) TokenLiteral() string { return funcLit.Token.Literal }

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }

type CallExpression struct {
	Token     token.Token
	Function  ExpressionNode
	Arguments []ExpressionNode
}

func (callExpr *CallExpression) expressionNode()      {}
func (callExpr *CallExpression) TokenLiteral() string { return callExpr.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    ExpressionNode
}

func (prefixExpr *PrefixExpression) expressionNode()      {}
func (prefixExpr *PrefixExpression) TokenLiteral() string { return prefixExpr.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (infixExpr *InfixExpression) expressionNode()      {}
func (infixExpr *InfixExpression) TokenLiteral() string { return infixExpr.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   ExpressionNode
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()      {}
func (ifExpr *IfExpression) TokenLiteral() string { return ifExpr.Token.Literal }
