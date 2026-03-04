package ast

import "funlang/token"

type SimpleType struct {
	Token token.Token
	Value string
}

func (simpType *SimpleType) typeNode()            {}
func (simpType *SimpleType) TokenLiteral() string { return simpType.Token.Literal }

type FunctionType struct {
	Token       token.Token
	ParamsTypes []TypeNode
	ReturnType  TypeNode
}

func (funcType *FunctionType) typeNode()            {}
func (funcType *FunctionType) TokenLiteral() string { return funcType.Token.Literal }
