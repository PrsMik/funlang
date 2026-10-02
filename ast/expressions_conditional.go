package ast

import "funlang/token"

type IfExpression struct {
	Token       token.Token
	Condition   ExpressionNode
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()       {}
func (ifExpr *IfExpression) TokenLiteral() string  { return ifExpr.Token.Literal }
func (ifExpr *IfExpression) Start() token.Position { return ifExpr.Token.Start }
func (ifExpr *IfExpression) End() token.Position {
	if ifExpr.Alternative != nil {
		return ifExpr.Alternative.End()
	}
	return ifExpr.Consequence.End()
}

type CaseBlock struct {
	Token   token.Token
	Pattern ExpressionNode
	Body    *BlockStatement
}

func (caseBlock *CaseBlock) TokenLiteral() string  { return caseBlock.Token.Literal }
func (caseBlock *CaseBlock) Start() token.Position { return caseBlock.Token.Start }
func (caseBlock *CaseBlock) End() token.Position {
	if caseBlock.Body != nil {
		return caseBlock.Body.End()
	}
	return caseBlock.Token.End
}

type SwitchExpression struct {
	Token     token.Token
	Value     ExpressionNode
	Cases     []*CaseBlock
	Default   *BlockStatement
	SemiToken token.Token
}

func (swExp *SwitchExpression) expressionNode()       {}
func (swExp *SwitchExpression) TokenLiteral() string  { return swExp.Token.Literal }
func (swExp *SwitchExpression) Start() token.Position { return swExp.Token.Start }
func (swExp *SwitchExpression) End() token.Position   { return swExp.SemiToken.End }
