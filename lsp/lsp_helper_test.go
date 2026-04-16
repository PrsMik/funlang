package lsp

import (
	"funlang/lexer"
	"funlang/parser"
	"funlang/type_checker"
	"funlang/types"
	"testing"
)

const testURI = "file:///test.fl"

// setupTestDocument имитирует открытие документа ,
// парсинг и проход type_checker для заполнения info
func setupTestDocument(content string) {
	documents[testURI] = content

	lxr := lexer.New(content)
	prs := parser.New(lxr)
	program := prs.ParseProgram()

	if program != nil {
		env := types.NewTypeEviroment()
		info := types.NewInfo()
		info.GlobalScope = env
		chk := type_checker.New(env, info)
		chk.CheckProgram(program)

		documentStates[testURI] = info
	}
}

func clearTestState() {
	delete(documents, testURI)
	delete(documentStates, testURI)
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
