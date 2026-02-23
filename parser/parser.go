package parser

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
)

type Parser struct {
	lxr *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(lxr *lexer.Lexer) *Parser {
	prs := &Parser{lxr: lxr}

	prs.nextToken()
	prs.nextToken()

	return prs
}

func (prs *Parser) nextToken() {
	prs.curToken = prs.peekToken
	prs.peekToken = prs.lxr.NextToken()
}

func (prs *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.StatementNode{}

	for prs.curToken.Type != token.EOF {
		statement := prs.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		prs.nextToken()
	}

	return program
}

func (prs *Parser) parseStatement() ast.StatementNode {
	switch prs.curToken.Type {
	case token.LET:
		return prs.parseLetStatement()
	default:
		return nil
	}
}

func (prs *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: prs.curToken}

	if !prs.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: prs.curToken, Value: prs.curToken.Literal}

	if !prs.expectPeek(token.COLON) {
		return nil
	}

	prs.nextToken()

	statement.Type = prs.parseType()
	if statement.Type == nil {
		return nil
	}

	if !prs.expectPeek(token.ASSIGN) {
		return nil
	}

	for !prs.curTokenIs(token.SEMICOLON) {
		prs.nextToken()
	}

	return statement
}

func (prs *Parser) curTokenIs(tknType token.TokenType) bool {
	return prs.curToken.Type == tknType
}

func (prs *Parser) peekTokenIs(tknType token.TokenType) bool {
	return prs.peekToken.Type == tknType
}

func (prs *Parser) expectPeek(tknType token.TokenType) bool {
	if prs.peekTokenIs(tknType) {
		prs.nextToken()
		return true
	} else {
		return false
	}
}

func (prs *Parser) parseType() ast.TypeNode {
	switch prs.curToken.Type {
	case token.INT_TYPE, token.BOOL_TYPE:
		return &ast.SimpleType{Token: prs.curToken, Value: prs.curToken.Literal}
	}
	return nil
}
