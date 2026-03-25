package ast

import "funlang/token"

type SimpleType struct {
	Token token.Token
	Value string
}

func (simpType *SimpleType) typeNode()            {}
func (simpType *SimpleType) TokenLiteral() string { return simpType.Token.Literal }

type ArrayType struct {
	Token        token.Token
	ElementsType TypeNode
}

func (arrType *ArrayType) typeNode()            {}
func (arrType *ArrayType) TokenLiteral() string { return arrType.Token.Literal }

type HashMapType struct {
	Token       token.Token
	KeyType     TypeNode
	ElementType TypeNode
}

func (hashMapType *HashMapType) typeNode()            {}
func (hashMapType *HashMapType) TokenLiteral() string { return hashMapType.Token.Literal }

type FunctionType struct {
	Token       token.Token
	ParamsTypes []TypeNode
	ReturnType  TypeNode
}

func (funcType *FunctionType) typeNode()            {}
func (funcType *FunctionType) TokenLiteral() string { return funcType.Token.Literal }
