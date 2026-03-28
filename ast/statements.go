package ast

import "funlang/token"

// инструкция let - токен инструкции (Token), идентификатор для корого let (Name),
// тип идентификатора (Type) и выражение справа от "=" (Value)
type LetStatement struct {
	Token     token.Token
	Name      *Identifier
	Type      TypeNode
	Value     ExpressionNode
	SemiToken token.Token
}

func (letStmt *LetStatement) statementNode()        {}
func (letStmt *LetStatement) TokenLiteral() string  { return letStmt.Token.Literal }
func (letStmt *LetStatement) Start() token.Position { return letStmt.Token.Start }
func (letStmt *LetStatement) End() token.Position   { return letStmt.SemiToken.End }

// инструкция return - токен инструкции (Token), выражение которое возвращает return (ReturnValue),
type ReturnStatement struct {
	Token     token.Token
	Value     ExpressionNode
	SemiToken token.Token
}

func (returnStmt *ReturnStatement) statementNode()        {}
func (returnStmt *ReturnStatement) TokenLiteral() string  { return returnStmt.Token.Literal }
func (returnStmt *ReturnStatement) Start() token.Position { return returnStmt.Token.Start }
func (returnStmt *ReturnStatement) End() token.Position   { return returnStmt.SemiToken.End }

type BlockStatement struct {
	Token      token.Token
	Statements []StatementNode
	SemiToken  token.Token
}

func (blockStmt *BlockStatement) statementNode()        {}
func (blockStmt *BlockStatement) TokenLiteral() string  { return blockStmt.Token.Literal }
func (blockStmt *BlockStatement) Start() token.Position { return blockStmt.Token.Start }
func (blockStmt *BlockStatement) End() token.Position   { return blockStmt.SemiToken.End }
