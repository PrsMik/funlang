package ast

import "funlang/token"

type SimpleType struct {
	Token token.Token
	Value string
}

func (t *SimpleType) typeNode()            {}
func (t *SimpleType) TokenLiteral() string { return t.Token.Literal }

type FunctionType struct {
	Token       token.Token
	ParamsTypes []TypeNode
	ReturnType  TypeNode
}

func (t *FunctionType) typeNode()            {}
func (t *FunctionType) TokenLiteral() string { return t.Token.Literal }
