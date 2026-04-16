package ast_test

import "funlang/token"

func createToken(t token.TokenType, literal string, line, col int) token.Token {
	return token.Token{
		Type:    t,
		Literal: literal,
		Start:   token.Position{Line: line, Column: col},
		End:     token.Position{Line: line, Column: col + len(literal)},
	}
}

func pos(line, col int) token.Position {
	return token.Position{Line: line, Column: col}
}
