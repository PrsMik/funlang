package lsp

import (
	"bytes"
	"funlang/formatter"
	"funlang/lexer"
	"funlang/parser"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	defer handlePanic(context)

	text, ok := documents[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	lxr := lexer.New(text)
	prs := parser.New(lxr)
	program := prs.ParseProgram()

	if len(prs.Errors()) > 0 {
		return nil, nil
	}

	var out bytes.Buffer
	fmtr := formatter.New(out)
	formattedText := fmtr.FormatProgram(program)

	lineCount := uint32(0)
	for _, char := range text {
		if char == '\n' {
			lineCount++
		}
	}

	fullRange := protocol.Range{
		Start: protocol.Position{Line: 0, Character: 0},
		End:   protocol.Position{Line: lineCount + 1, Character: 0},
	}

	return []protocol.TextEdit{
		{
			Range:   fullRange,
			NewText: formattedText,
		},
	}, nil
}
