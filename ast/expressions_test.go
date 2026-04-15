// ast/expressions_test.go
package ast

import (
	"funlang/token"
	"testing"
)

func TestExpressions(t *testing.T) {
	tests := []struct {
		name                 string
		node                 ExpressionNode
		expectedTokenLiteral string
		expectedString       string
		expectedStart        token.Position
		expectedEnd          token.Position
	}{
		{
			name:                 "BadExpression",
			node:                 &BadExpression{From: pos(1, 1), To: pos(1, 5)},
			expectedTokenLiteral: "",
			expectedString:       "",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 5),
		},
		{
			name:                 "IntegerLiteral",
			node:                 &IntegerLiteral{Token: createToken(token.INT, "5", 1, 1), Value: 5},
			expectedTokenLiteral: "5",
			expectedString:       "5",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 2),
		},
		{
			name:                 "BooleanLiteral",
			node:                 &BooleanLiteral{Token: createToken(token.TRUE, "true", 1, 1), Value: true},
			expectedTokenLiteral: "true",
			expectedString:       "true",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 5),
		},
		{
			name:                 "StringLiteral",
			node:                 &StringLiteral{Token: createToken(token.STRING, "hello", 1, 1), Value: "hello"},
			expectedTokenLiteral: "hello",
			expectedString:       "hello",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 6),
		},
		{
			name:                 "Identifier",
			node:                 &Identifier{Token: createToken(token.IDENT, "myVar", 1, 1), Value: "myVar"},
			expectedTokenLiteral: "myVar",
			expectedString:       "myVar",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 6),
		},
		{
			name: "ArrayLiteral",
			node: &ArrayLiteral{
				Token:     createToken(token.LBRACKET, "[", 1, 1),
				Elements:  []ExpressionNode{&IntegerLiteral{Token: createToken(token.INT, "1", 1, 2)}},
				SemiToken: createToken(token.RBRACKET, "]", 1, 3),
			},
			expectedTokenLiteral: "[",
			expectedString:       "[1]",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 4),
		},
		{
			name: "HashMapLiteral",
			node: &HashMapLiteral{
				Token: createToken(token.LBRACE, "{", 1, 1),
				Pairs: map[ExpressionNode]ExpressionNode{
					&Identifier{Value: "a"}: &IntegerLiteral{Token: createToken(token.INT, "1", 1, 5)},
				},
				SemiToken: createToken(token.RBRACE, "}", 1, 6),
			},
			expectedTokenLiteral: "{",
			expectedString:       "{a:1}",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 7),
		},
		{
			name: "FunctionLiteral",
			node: &FunctionLiteral{
				Token:      createToken(token.FN, "fn", 1, 1),
				Parameters: []*Identifier{{Value: "x"}},
				Body:       &BlockStatement{SemiToken: createToken(token.RBRACE, "}", 1, 15)},
			},
			expectedTokenLiteral: "fn",
			expectedString:       "fn(x) ",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 16),
		},
		{
			name: "CallExpression",
			node: &CallExpression{
				Token:     createToken(token.LPAREN, "(", 1, 5),
				Function:  &Identifier{Token: createToken(token.IDENT, "add", 1, 1), Value: "add"},
				Arguments: []ExpressionNode{&Identifier{Value: "a"}},
				SemiToken: createToken(token.RPAREN, ")", 1, 8),
			},
			expectedTokenLiteral: "(",
			expectedString:       "add(a)",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 9),
		},
		{
			name: "PrefixExpression",
			node: &PrefixExpression{
				Token:    createToken(token.BANG, "!", 1, 1),
				Operator: "!",
				Right:    &Identifier{Token: createToken(token.IDENT, "a", 1, 2), Value: "a"},
			},
			expectedTokenLiteral: "!",
			expectedString:       "(!a)",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 3),
		},
		{
			name: "InfixExpression",
			node: &InfixExpression{
				Token:    createToken(token.PLUS, "+", 1, 3),
				Operator: "+",
				Left:     &Identifier{Token: createToken(token.IDENT, "a", 1, 1), Value: "a"},
				Right:    &Identifier{Token: createToken(token.IDENT, "b", 1, 5), Value: "b"},
			},
			expectedTokenLiteral: "+",
			expectedString:       "(a + b)",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 6),
		},
		{
			name: "IndexExpression",
			node: &IndexExpression{
				Token:     createToken(token.LBRACKET, "[", 1, 2),
				Left:      &Identifier{Token: createToken(token.IDENT, "a", 1, 1), Value: "a"},
				Index:     &IntegerLiteral{Token: createToken(token.INT, "1", 1, 3), Value: 1},
				SemiToken: createToken(token.RBRACKET, "]", 1, 4),
			},
			expectedTokenLiteral: "[",
			expectedString:       "(a[1])",
			expectedStart:        pos(1, 2),
			expectedEnd:          pos(1, 5),
		},
		{
			name: "IfExpression (with Alternative)",
			node: &IfExpression{
				Token:       createToken(token.IF, "if", 1, 1),
				Condition:   &Identifier{Value: "x"},
				Consequence: &BlockStatement{Token: createToken(token.LBRACE, "{", 1, 7)},
				Alternative: &BlockStatement{SemiToken: createToken(token.RBRACE, "}", 1, 20)},
			},
			expectedTokenLiteral: "if",
			expectedString:       "ifx  else ",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 21),
		},
		{
			name: "IfExpression (without Alternative)",
			node: &IfExpression{
				Token:       createToken(token.IF, "if", 1, 1),
				Condition:   &Identifier{Value: "x"},
				Consequence: &BlockStatement{SemiToken: createToken(token.RBRACE, "}", 1, 10)},
				Alternative: nil,
			},
			expectedTokenLiteral: "if",
			expectedString:       "ifx ",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 11),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.expressionNode()

			if got := tt.node.TokenLiteral(); got != tt.expectedTokenLiteral {
				t.Errorf("TokenLiteral() wrong. expected=%q, got=%q", tt.expectedTokenLiteral, got)
			}
			if got := tt.node.String(); got != tt.expectedString {
				t.Errorf("String() wrong. expected=%q, got=%q", tt.expectedString, got)
			}
			if got := tt.node.Start(); got != tt.expectedStart {
				t.Errorf("Start() wrong. expected=%v, got=%v", tt.expectedStart, got)
			}
			if got := tt.node.End(); got != tt.expectedEnd {
				t.Errorf("End() wrong. expected=%v, got=%v", tt.expectedEnd, got)
			}
		})
	}
}
