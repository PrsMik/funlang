package lsp

import (
	"funlang/ast"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentReferences(context *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	defer handlePanic(context)

	chk, ok := documentStates[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	pos := params.Position
	var targetDeclNode ast.Node
	var minLen int = 9999999

	// проверка что узел - исопльзование
	for usageNode, declNode := range chk.Definitions {
		start := usageNode.Start()
		end := usageNode.End()

		if isPosInside(pos, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)
			if length < minLen {
				minLen = length
				targetDeclNode = declNode
			}
		}
	}

	// проверка что узел - объявление
	if targetDeclNode == nil {
		for _, declNode := range chk.Definitions {
			if declNode == nil {
				continue
			}
			start := declNode.Start()
			end := declNode.End()

			if isPosInside(pos, start, end) {
				length := (end.Line-start.Line)*1000 + (end.Column - start.Column)
				if length < minLen {
					minLen = length
					targetDeclNode = declNode
				}
			}
		}
	}

	if targetDeclNode == nil {
		return nil, nil
	}

	var locations []protocol.Location

	if params.Context.IncludeDeclaration {
		locations = append(locations, protocol.Location{
			URI:   params.TextDocument.URI,
			Range: createLspRange(targetDeclNode.Start(), targetDeclNode.End()),
		})
	}

	// поиск использований, которые ссылаются на найденное объявление
	for usageNode, declNode := range chk.Definitions {

		if declNode == targetDeclNode {
			locations = append(locations, protocol.Location{
				URI:   params.TextDocument.URI,
				Range: createLspRange(usageNode.Start(), usageNode.End()),
			})
		}
	}

	return locations, nil
}
