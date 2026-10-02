package ast

import "funlang/token"

type BadExpression struct {
	From token.Position
	To   token.Position
}

func (un *BadExpression) expressionNode()       {}
func (un *BadExpression) TokenLiteral() string  { return "" }
func (un *BadExpression) Start() token.Position { return un.From }
func (un *BadExpression) End() token.Position   { return un.To }

type ComptimeExpression struct {
	Token token.Token
	Value ExpressionNode
}

func (ce *ComptimeExpression) expressionNode()       {}
func (ce *ComptimeExpression) TokenLiteral() string  { return ce.Token.Literal }
func (ce *ComptimeExpression) Start() token.Position { return ce.Token.Start }
func (ce *ComptimeExpression) End() token.Position   { return ce.Value.End() }

type CallExpression struct {
	Token     token.Token
	Function  ExpressionNode
	Arguments []ExpressionNode
	SemiToken token.Token
}

func (callExpr *CallExpression) expressionNode()       {}
func (callExpr *CallExpression) TokenLiteral() string  { return callExpr.Token.Literal }
func (callExpr *CallExpression) Start() token.Position { return callExpr.Function.Start() }
func (callExpr *CallExpression) End() token.Position   { return callExpr.SemiToken.End }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    ExpressionNode
}

func (prefixExpr *PrefixExpression) expressionNode()       {}
func (prefixExpr *PrefixExpression) TokenLiteral() string  { return prefixExpr.Token.Literal }
func (prefixExpr *PrefixExpression) Start() token.Position { return prefixExpr.Token.Start }
func (prefixExpr *PrefixExpression) End() token.Position   { return prefixExpr.Right.End() }

type InfixExpression struct {
	Token    token.Token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (infixExpr *InfixExpression) expressionNode()       {}
func (infixExpr *InfixExpression) TokenLiteral() string  { return infixExpr.Token.Literal }
func (infixExpr *InfixExpression) Start() token.Position { return infixExpr.Left.Start() }
func (infixExpr *InfixExpression) End() token.Position   { return infixExpr.Right.End() }

type IndexExpression struct {
	Token     token.Token
	Left      ExpressionNode
	Index     ExpressionNode
	SemiToken token.Token
}

func (indExpr *IndexExpression) expressionNode()       {}
func (indExpr *IndexExpression) TokenLiteral() string  { return indExpr.Token.Literal }
func (indExpr *IndexExpression) Start() token.Position { return indExpr.Token.Start }
func (indExpr *IndexExpression) End() token.Position   { return indExpr.SemiToken.End }
