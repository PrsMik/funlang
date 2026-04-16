package formatter_test

import (
	"bytes"
	"funlang/formatter"
	"funlang/lexer"
	"funlang/parser"
	"strings"
	"testing"
)

func assertFormat(t *testing.T, name, input, expected string) {
	t.Run(name, func(t *testing.T) {
		t.Helper()

		lxr := lexer.New(input)
		prs := parser.New(lxr)
		prog := prs.ParseProgram()

		if len(prs.Errors()) != 0 {
			t.Fatalf("Parser encountered errors: %v", prs.Errors())
		}

		var out bytes.Buffer
		fmtr := formatter.New(&out)
		got := fmtr.FormatProgram(prog)

		expectedClean := strings.TrimLeft(expected, "\n")
		gotClean := strings.TrimLeft(got, "\n")

		if gotClean != expectedClean {
			t.Errorf("FormatProgram() mismatch.\n=== GOT ===\n%s\n=== EXPECTED ===\n%s", gotClean, expectedClean)
		}
	})
}
