package ast

import "funlang/token"

type MemberAccessExpression struct {
	Token    token.Token
	Left     ExpressionNode
	Property *Identifier
}

func (indExpr *MemberAccessExpression) expressionNode()       {}
func (indExpr *MemberAccessExpression) TokenLiteral() string  { return indExpr.Token.Literal }
func (indExpr *MemberAccessExpression) Start() token.Position { return indExpr.Left.Start() }
func (indExpr *MemberAccessExpression) End() token.Position   { return indExpr.Property.End() }

type TypeAccessExpression struct {
	Token     token.Token
	Left      ExpressionNode
	SemiToken token.Token
}

func (tae *TypeAccessExpression) expressionNode()       {}
func (tae *TypeAccessExpression) TokenLiteral() string  { return tae.Token.Literal }
func (tae *TypeAccessExpression) Start() token.Position { return tae.Left.Start() }
func (tae *TypeAccessExpression) End() token.Position   { return tae.SemiToken.End }

// выражение создания экземпляра (например, foo.{bar : 42} )
type ImplInstantiationExpression struct {
	Token  token.Token
	Left   ExpressionNode
	Fields ExpressionNode
}

func (impInst *ImplInstantiationExpression) expressionNode()       {}
func (impInst *ImplInstantiationExpression) TokenLiteral() string  { return impInst.Token.Literal }
func (impInst *ImplInstantiationExpression) Start() token.Position { return impInst.Token.Start }
func (impInst *ImplInstantiationExpression) End() token.Position   { return impInst.Fields.End() }
