package ast

import "funlang/token"

type BadExpression struct {
	From token.Position
	To   token.Position
}

func (un *BadExpression) expressionNode()       {}
func (un *BadExpression) TokenLiteral() string  { return "" }
func (un *BadExpression) String() string        { return "" }
func (un *BadExpression) Start() token.Position { return un.From }
func (un *BadExpression) End() token.Position   { return un.To }

// литерал типа int (например "5")
type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (intLit *IntegerLiteral) expressionNode()       {}
func (intLit *IntegerLiteral) TokenLiteral() string  { return intLit.Token.Literal }
func (intLit *IntegerLiteral) Start() token.Position { return intLit.Token.Start }
func (intLit *IntegerLiteral) End() token.Position   { return intLit.Token.End }

// литерал типа bool (например "true")
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (boolLit *BooleanLiteral) expressionNode()       {}
func (boolLit *BooleanLiteral) TokenLiteral() string  { return boolLit.Token.Literal }
func (boolLit *BooleanLiteral) Start() token.Position { return boolLit.Token.Start }
func (boolLit *BooleanLiteral) End() token.Position   { return boolLit.Token.End }

// строковый литерал (например `"Hello world!"`)
type StringLiteral struct {
	Token token.Token
	Value string
}

func (strLit *StringLiteral) expressionNode()       {}
func (strLit *StringLiteral) TokenLiteral() string  { return strLit.Token.Literal }
func (strLit *StringLiteral) Start() token.Position { return strLit.Token.Start }
func (strLit *StringLiteral) End() token.Position   { return strLit.Token.End }

// литерал массива (например "[1, 2, 3]")
type ArrayLiteral struct {
	Token     token.Token
	Elements  []ExpressionNode
	SemiToken token.Token
}

func (arrLit *ArrayLiteral) expressionNode()       {}
func (arrLit *ArrayLiteral) TokenLiteral() string  { return arrLit.Token.Literal }
func (arrLit *ArrayLiteral) Start() token.Position { return arrLit.Token.Start }
func (arrLit *ArrayLiteral) End() token.Position   { return arrLit.SemiToken.End }

type HashMapLiteral struct {
	Token     token.Token
	Pairs     map[ExpressionNode]ExpressionNode
	SemiToken token.Token
}

func (hashMapLit *HashMapLiteral) expressionNode()       {}
func (hashMapLit *HashMapLiteral) TokenLiteral() string  { return hashMapLit.Token.Literal }
func (hashMapLit *HashMapLiteral) Start() token.Position { return hashMapLit.Token.Start }
func (hashMapLit *HashMapLiteral) End() token.Position   { return hashMapLit.SemiToken.End }

// литерал функции (например "fn(x) { return x; }")
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	ParamTypes []TypeNode
	ReturnType TypeNode
	Body       *BlockStatement
}

func (funcLit *FunctionLiteral) expressionNode()       {}
func (funcLit *FunctionLiteral) TokenLiteral() string  { return funcLit.Token.Literal }
func (funcLit *FunctionLiteral) Start() token.Position { return funcLit.Token.Start }
func (funcLit *FunctionLiteral) End() token.Position   { return funcLit.Body.End() }

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode()       {}
func (ident *Identifier) TokenLiteral() string  { return ident.Token.Literal }
func (ident *Identifier) Start() token.Position { return ident.Token.Start }
func (ident *Identifier) End() token.Position   { return ident.Token.End }

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

type IfExpression struct {
	Token       token.Token
	Condition   ExpressionNode
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()       {}
func (ifExpr *IfExpression) TokenLiteral() string  { return ifExpr.Token.Literal }
func (ifExpr *IfExpression) Start() token.Position { return ifExpr.Token.Start }
func (ifExpr *IfExpression) End() token.Position   { return ifExpr.Alternative.End() }

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
