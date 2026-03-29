package lsp

import (
	"fmt"
	"funlang/lexer"
	"funlang/parser"
	"funlang/token"
	"funlang/type_checker"
	"funlang/types"
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	parseSource   string = "funlang-parser"
	checkerSource string = "funlang-typechecker"
	severity             = protocol.DiagnosticSeverityError
)

func validateDocument(context *glsp.Context, uri string, text string) {
	diagnostics := []protocol.Diagnostic{}

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Recovered from panic during LSP validation: %v\n", r)
		}

		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		})
	}()

	lxr := lexer.New(text)
	prs := parser.New(lxr)
	program := prs.ParseProgram()

	for _, err := range prs.Errors() {

		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    createLspRange(err.Token.Start, err.Token.End),
			Severity: &severity,
			Source:   &parseSource,
			Message:  err.Msg,
		})
	}

	if program != nil {
		env := types.NewTypeEviroment()
		chk := type_checker.New(env)
		chk.CheckProgram(program)

		for _, err := range chk.Errors() {
			// length := 1
			// if err.Node != nil {
			// 	length = len(err.Node.TokenLiteral())
			// }

			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range:    createLspRange(err.Node.Start(), err.Node.End()),
				Severity: &severity,
				Source:   &checkerSource,
				Message:  err.Msg,
			})
		}
	}

}

func createLspRange(start token.Position, end token.Position) protocol.Range {
	startLine := uint32(start.Line)
	startCol := uint32(start.Column)

	endLine := uint32(end.Line)
	endCol := uint32(end.Column)

	// length :=

	// if length <= 0 {
	// 	length = 1
	// }

	return protocol.Range{
		Start: protocol.Position{Line: startLine, Character: startCol},
		End:   protocol.Position{Line: endLine, Character: endCol},
	}
}
