package parser_test

import (
	"funlang/lexer"
	"funlang/parser"
	"testing"
)

func TestCommentCollection(t *testing.T) {
	input := `
let x: int = 5; // comment 1
// comment 2
return x; // comment 3
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectedComments := []string{
		" comment 1",
		" comment 2",
		" comment 3",
	}

	if len(program.Comments) != len(expectedComments) {
		t.Fatalf("expected %d comments, got %d", len(expectedComments), len(program.Comments))
	}

	for i, expected := range expectedComments {
		if program.Comments[i].Literal != expected {
			t.Errorf("comment %d: expected %q, got %q", i, expected, program.Comments[i].Literal)
		}
	}
}
