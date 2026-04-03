package lsp

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Fprintf(os.Stderr, "Recovered from panic during LSP validation: %v\n", r)
	// 	}

	// 	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
	// 		URI:         params.TextDocument.URI,
	// 		Diagnostics: []protocol.Diagnostic{},
	// 	})
	// }()

	chk, ok := documentStates[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	pos := params.Position

	var hoveredNode ast.Node
	var hoveredType types.Type
	var minLen int = 9999999

	for node, tp := range chk.TypesInfo {
		start := node.Start()
		end := node.End()

		if isPosInside(pos, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)

			if length < minLen {
				minLen = length
				hoveredNode = node
				hoveredType = tp
			}
		}
	}

	if hoveredNode != nil && hoveredType != nil {
		signature := hoveredType.Signature()

		markdown := fmt.Sprintf("```funlang\n%s: %s\n```", hoveredNode.TokenLiteral(), signature)

		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: markdown,
			},
		}, nil
	}

	return nil, nil
}

func isPosInside(pos protocol.Position, start token.Position, end token.Position) bool {
	sl, sc := uint32(start.Line), uint32(start.Column)
	el, ec := uint32(end.Line), uint32(end.Column)
	pl, pc := pos.Line, pos.Character

	if pl < sl || pl > el {
		return false
	}
	if pl == sl && pc < sc {
		return false
	}
	if pl == el && pc > ec {
		return false
	}
	return true
}
