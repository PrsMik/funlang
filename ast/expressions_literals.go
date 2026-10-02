package ast

import "funlang/token"

// идентификатор - токен и литерал
type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode()       {}
func (ident *Identifier) TokenLiteral() string  { return ident.Token.Literal }
func (ident *Identifier) Start() token.Position { return ident.Token.Start }
func (ident *Identifier) End() token.Position   { return ident.Token.End }

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
	ParamTypes []ExpressionNode
	ReturnType ExpressionNode
	Body       *BlockStatement
}

func (funcLit *FunctionLiteral) expressionNode()       {}
func (funcLit *FunctionLiteral) TokenLiteral() string  { return funcLit.Token.Literal }
func (funcLit *FunctionLiteral) Start() token.Position { return funcLit.Token.Start }
func (funcLit *FunctionLiteral) End() token.Position   { return funcLit.Body.End() }

// литерал объявления интерфейса ( например, interface { let x: type = int; } )
type InterfaceLiteral struct {
	Token token.Token
	Body  *BlockStatement
}

func (iface *InterfaceLiteral) expressionNode()       {}
func (iface *InterfaceLiteral) TokenLiteral() string  { return iface.Token.Literal }
func (iface *InterfaceLiteral) Start() token.Position { return iface.Token.Start }
func (iface *InterfaceLiteral) End() token.Position   { return iface.Body.End() }

// литерал объявления реализации (например, impl { let x: int = 1; } )
type ImplLiteral struct {
	Token      token.Token
	TargetType ExpressionNode
	Body       *BlockStatement
}

func (imp *ImplLiteral) expressionNode()       {}
func (imp *ImplLiteral) TokenLiteral() string  { return imp.Token.Literal }
func (imp *ImplLiteral) Start() token.Position { return imp.Token.Start }
func (imp *ImplLiteral) End() token.Position   { return imp.Body.End() }
