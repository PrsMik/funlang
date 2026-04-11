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
	defer handlePanic(context)

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
	posAbs := pos.Line*10000 + pos.Character
	startAbs := uint32(start.Line)*10000 + uint32(start.Column)
	endAbs := uint32(end.Line)*10000 + uint32(end.Column)

	return posAbs >= startAbs && posAbs <= endAbs
}
