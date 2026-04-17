package token_test

import (
	"funlang/token"
	"testing"
)

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name     string
		position token.Position
		expected string
	}{
		{
			name:     "Zero values",
			position: token.Position{Line: 0, Column: 0},
			expected: "0:0",
		},
		{
			name:     "Typical position",
			position: token.Position{Line: 10, Column: 5},
			expected: "10:5",
		},
		{
			name:     "Large position numbers",
			position: token.Position{Line: 9999, Column: 1234},
			expected: "9999:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.position.String()
			if actual != tt.expected {
				t.Errorf("Position.String() error. expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
