package lsp

import (
	"funlang/ast"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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

	var clickedNode ast.Node
	var minLen int = 9999999

	for node := range chk.Definitions {
		start := node.Start()
		end := node.End()

		if isPosInside(pos, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)
			if length < minLen {
				minLen = length
				clickedNode = node
			}
		}
	}

	if clickedNode != nil {
		declNode := chk.Definitions[clickedNode]
		if declNode != nil {
			return []protocol.Location{
				{
					URI:   params.TextDocument.URI,
					Range: createLspRange(declNode.Start(), declNode.End()),
				},
			}, nil
		}
	}

	return nil, nil
}
