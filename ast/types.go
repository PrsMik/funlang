package ast

import "funlang/token"

type TypeType struct {
	Token token.Token
	Value string
}

func (typeType *TypeType) expressionNode()       {}
func (typeType *TypeType) TokenLiteral() string  { return typeType.Token.Literal }
func (typeType *TypeType) Start() token.Position { return typeType.Token.Start }
func (typeType *TypeType) End() token.Position   { return typeType.Token.End }

type SimpleType struct {
	Token token.Token
	Value string
}

func (simpType *SimpleType) expressionNode()       {}
func (simpType *SimpleType) TokenLiteral() string  { return simpType.Token.Literal }
func (simpType *SimpleType) Start() token.Position { return simpType.Token.Start }
func (simpType *SimpleType) End() token.Position   { return simpType.Token.End }

type ArrayType struct {
	Token        token.Token
	ElementsType ExpressionNode
}

func (arrType *ArrayType) expressionNode()       {}
func (arrType *ArrayType) TokenLiteral() string  { return arrType.Token.Literal }
func (arrType *ArrayType) Start() token.Position { return arrType.Token.Start }
func (arrType *ArrayType) End() token.Position   { return arrType.ElementsType.End() }

type HashMapType struct {
	Token       token.Token
	KeyType     ExpressionNode
	ElementType ExpressionNode
}

func (hashMapType *HashMapType) expressionNode()       {}
func (hashMapType *HashMapType) TokenLiteral() string  { return hashMapType.Token.Literal }
func (hashMapType *HashMapType) Start() token.Position { return hashMapType.Token.Start }
func (hashMapType *HashMapType) End() token.Position   { return hashMapType.ElementType.End() }

type FunctionType struct {
	Token       token.Token
	ParamsTypes []ExpressionNode
	ReturnType  ExpressionNode
}

func (funcType *FunctionType) expressionNode()      {}
func (funcType *FunctionType) TokenLiteral() string { return funcType.Token.Literal }
func (funcType *FunctionType) Start() token.Position {
	return funcType.Token.Start
}
func (funcType *FunctionType) End() token.Position { return funcType.ReturnType.End() }
