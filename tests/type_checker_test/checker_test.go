package type_checker_test

import (
	"strings"
	"testing"

	"funlang/lexer"
	"funlang/parser"
	"funlang/type_checker"
	"funlang/types"
)

type TestCase struct {
	name        string
	input       string
	expectedErr string
}

func runTypeCheckerTests(t *testing.T, tests []TestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			prog := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser error in test '%s': %v", tt.name, p.Errors())
			}

			env := types.NewTypeEviroment()
			chk := type_checker.New(env, nil)
			chk.CheckProgram(prog)

			errors := chk.Errors()

			if tt.expectedErr == "" {
				if len(errors) != 0 {
					t.Fatalf("Expected no errors, got %d errors. First error: %q", len(errors), errors[0].Msg)
				}
			} else {
				if len(errors) == 0 {
					t.Fatalf("Expected error containing %q, but got none", tt.expectedErr)
				}

				match := false
				for _, err := range errors {
					if strings.Contains(err.Msg, tt.expectedErr) {
						match = true
						break
					}
				}

				if !match {
					t.Errorf("Expected error containing %q, but got: %q", tt.expectedErr, errors[0].Msg)
				}
			}
		})
	}
}
