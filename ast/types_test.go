package ast

import (
	"funlang/token"
	"testing"
)

func TestTypes(t *testing.T) {
	tests := []struct {
		name                 string
		node                 TypeNode
		expectedTokenLiteral string
		expectedString       string
		expectedStart        token.Position
		expectedEnd          token.Position
	}{
		{
			name: "SimpleType",
			node: &SimpleType{
				Token: createToken(token.INT_TYPE, "int", 1, 1),
				Value: "int",
			},
			expectedTokenLiteral: "int",
			expectedString:       "int",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 4),
		},
		{
			name: "ArrayType",
			node: &ArrayType{
				Token: createToken(token.LBRACKET, "[", 1, 1),
				ElementsType: &SimpleType{
					Token: createToken(token.INT_TYPE, "int", 1, 2),
					Value: "int",
				},
			},
			expectedTokenLiteral: "[",
			expectedString:       "[int]",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 5),
		},
		{
			name: "HashMapType",
			node: &HashMapType{
				Token:   createToken(token.LBRACE, "{", 1, 1),
				KeyType: &SimpleType{Value: "string"},
				ElementType: &SimpleType{
					Token: createToken(token.INT_TYPE, "int", 1, 10),
					Value: "int",
				},
			},
			expectedTokenLiteral: "{",
			expectedString:       "{string : int}",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 13),
		},
		{
			name: "FunctionType (with params and return)",
			node: &FunctionType{
				Token: createToken(token.FN, "fn", 1, 1),
				ParamsTypes: []TypeNode{
					&SimpleType{Value: "int"},
					&SimpleType{Value: "string"},
				},
				ReturnType: &SimpleType{
					Token: createToken(token.BOOL_TYPE, "bool", 1, 25),
					Value: "bool",
				},
			},
			expectedTokenLiteral: "fn",
			expectedString:       "fn(int, string) -> bool",
			expectedStart:        pos(1, 1),
			expectedEnd:          pos(1, 29),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.typeNode()

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
