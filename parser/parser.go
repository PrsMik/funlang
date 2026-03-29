package parser

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
)

const (
	_ int = iota
	LOWEST
	LOG_SUM     // ||
	LOG_PROD    // &&
	EQUALS      // ==
	LESSGREATER // >, <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X, !X
	CALL        // myFunction(X)
	INDEX       // arr[i]
)

var precedences = map[token.TokenType]int{
	token.OR:               LOG_SUM,
	token.AND:              LOG_PROD,
	token.EQUAL:            EQUALS,
	token.NOT_EQUAL:        EQUALS,
	token.LESS:             LESSGREATER,
	token.GREATER:          LESSGREATER,
	token.LESS_OR_EQUAL:    LESSGREATER,
	token.GREATER_OR_EQUAL: LESSGREATER,
	token.PLUS:             SUM,
	token.MINUS:            SUM,
	token.SLASH:            PRODUCT,
	token.ASTERISK:         PRODUCT,
	token.LPAREN:           CALL,
	token.LBRACKET:         INDEX,
}

type ParseError struct {
	Msg   string
	Token token.Token
}

type Parser struct {
	lxr    *lexer.Lexer
	errors []ParseError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lxr *lexer.Lexer) *Parser {
	prs := &Parser{lxr: lxr, errors: []ParseError{}}

	prs.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	prs.registerPrefix(token.IDENT, prs.parseIdentifier)
	prs.registerPrefix(token.INT, prs.parseIntegerLiteral)
	prs.registerPrefix(token.STRING, prs.parseStringLiteral)

	prs.registerPrefix(token.BANG, prs.parsePrefixExpression)
	prs.registerPrefix(token.MINUS, prs.parsePrefixExpression)

	prs.registerPrefix(token.TRUE, prs.parseBoolean)
	prs.registerPrefix(token.FALSE, prs.parseBoolean)

	prs.registerPrefix(token.LPAREN, prs.parseGroupedExpression)
	prs.registerPrefix(token.LBRACKET, prs.parseArrayLiteral)
	prs.registerPrefix(token.LBRACE, prs.parseHashMapLiteral)

	prs.registerPrefix(token.IF, prs.parseIfExpression)
	prs.registerPrefix(token.FN, prs.parseFunctionLiteral)

	prs.infixParseFns = make(map[token.TokenType]infixParseFn)
	prs.registerInfix(token.EQUAL, prs.parseInfixExpression)
	prs.registerInfix(token.NOT_EQUAL, prs.parseInfixExpression)

	prs.registerInfix(token.LESS, prs.parseInfixExpression)
	prs.registerInfix(token.GREATER, prs.parseInfixExpression)
	prs.registerInfix(token.LESS_OR_EQUAL, prs.parseInfixExpression)
	prs.registerInfix(token.GREATER_OR_EQUAL, prs.parseInfixExpression)

	prs.registerInfix(token.PLUS, prs.parseInfixExpression)
	prs.registerInfix(token.MINUS, prs.parseInfixExpression)

	prs.registerInfix(token.SLASH, prs.parseInfixExpression)
	prs.registerInfix(token.ASTERISK, prs.parseInfixExpression)

	prs.registerInfix(token.AND, prs.parseInfixExpression)
	prs.registerInfix(token.OR, prs.parseInfixExpression)

	prs.registerInfix(token.LPAREN, prs.parseCallExpression)
	prs.registerInfix(token.LBRACKET, prs.parseIndexExpression)

	prs.nextToken()
	prs.nextToken()

	return prs
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

func (prs *Parser) peekTokenPrecedence() int {
	if precedence, ok := precedences[prs.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (prs *Parser) curTokenPrecedence() int {
	if precedence, ok := precedences[prs.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
