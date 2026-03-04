package ast

import "funlang/token"

// инструкция let - токен инструкции (Token), идентификатор для корого let (Name),
// тип идентификатора (Type) и выражение справа от "=" (Value)
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Type  TypeNode
	Value ExpressionNode
}

func (letStmt *LetStatement) statementNode()       {}
func (letStmt *LetStatement) TokenLiteral() string { return letStmt.Token.Literal }

// инструкция return - токен инструкции (Token), выражение которое возвращает return (ReturnValue),
type ReturnStatement struct {
	Token token.Token
	Value ExpressionNode
}

func (returnStmt *ReturnStatement) statementNode()       {}
func (returnStmt *ReturnStatement) TokenLiteral() string { return returnStmt.Token.Literal }

type BlockStatement struct {
	Token      token.Token
	Statements []StatementNode
}

func (blockStmt *BlockStatement) statementNode()       {}
func (blockStmt *BlockStatement) TokenLiteral() string { return blockStmt.Token.Literal }
