package parser_test

import (
	"funlang/lexer"
	"funlang/parser"
	"testing"
)

func TestParserRecovery(t *testing.T) {
	input := `
let a: int = ;    // Ошибка: пропущено значение
let b: int = 10;  // Это должно распарситься корректно
return a + ;      // Ошибка: неполное выражение
return b;         // Это должно распарситься корректно
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(p.Errors()))
	}

	if len(program.Statements) < 2 {
		t.Errorf("parser failed to recover and collect valid statements. got=%d", len(program.Statements))
	}
}
