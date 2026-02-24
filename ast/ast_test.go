package ast

import (
	"funlang/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []StatementNode{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Type: &SimpleType{
					Token: token.Token{Type: token.INT_TYPE, Literal: "int"},
					Value: "int",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar: int = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
