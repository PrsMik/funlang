package lsp

import (
	"funlang/ast"
	"funlang/lexer"
	"funlang/parser"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDocumentSymbol(context *glsp.Context, params *protocol.DocumentSymbolParams) (any, error) {
	defer handlePanic(context)

	text, ok := documents[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	lxr := lexer.New(text)
	prs := parser.New(lxr)
	program := prs.ParseProgram()

	if program == nil {
		return nil, nil
	}

	symbols := getSymbolsFromStatements(program.Statements)

	return symbols, nil
}

// getSymbolsFromStatements проходит по списку инструкций и извлекает объявленные сущности
func getSymbolsFromStatements(statements []ast.StatementNode) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.LetStatement:
			sym := createLetSymbol(s)
			symbols = append(symbols, sym)
		}
	}

	return symbols
}

func createLetSymbol(letStmt *ast.LetStatement) protocol.DocumentSymbol {
	name := letStmt.Name.Value
	kind := protocol.SymbolKindVariable

	var detail string
	if letStmt.Type != nil {
		detail = letStmt.Type.String()
	}

	var children []protocol.DocumentSymbol

	if funLiteral, isFun := letStmt.Value.(*ast.FunctionLiteral); isFun {
		kind = protocol.SymbolKindFunction
		if funLiteral.Body != nil {
			children = getSymbolsFromStatements(funLiteral.Body.Statements)
		}
	}

	fullRange := createLspRange(letStmt.Start(), letStmt.End())

	selectionRange := createLspRange(letStmt.Name.Start(), letStmt.Name.End())

	return protocol.DocumentSymbol{
		Name:           name,
		Detail:         &detail,
		Kind:           kind,
		Range:          fullRange,
		SelectionRange: selectionRange,
		Children:       children,
	}
}
