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

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// инструкция return - токен инструкции (Token), выражение которое возвращает return (ReturnValue),
type ReturnStatement struct {
	Token token.Token
	Value ExpressionNode
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type BlockStatement struct {
	Token      token.Token
	Statements []StatementNode
}

func (blockStmt *BlockStatement) statementNode()       {}
func (blockStmt *BlockStatement) TokenLiteral() string { return blockStmt.Token.Literal }
