package lsp

import (
	"fmt"
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
	"funlang/type_checker"
	"funlang/types"
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
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

	var closestScope *ast.BlockStatement
	var minLen int = 999999

	for node := range chk.Scopes {
		block, isBlock := node.(*ast.BlockStatement)
		if !isBlock {
			continue
		}

		start := block.Start()
		end := block.End()

		if isPosInside(params.Position, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)
			if length < minLen {
				minLen = length
				closestScope = block
			}
		}
	}

	var hoveredNode ast.Node
	var hoveredType types.Type

	for node, tp := range chk.ExpectedTypes {
		start := node.Start()
		end := node.End()

		if isPosInside(params.Position, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)

			if length < minLen {
				minLen = length
				hoveredNode = node
				hoveredType = tp
			}
		}
	}

	env := chk.GetEnv()
	if closestScope != nil {
		env = *chk.Scopes[closestScope]
	}

	docText, ok := documents[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	items := []protocol.CompletionItem{}

	lastTok := getPrecedingToken(docText, params.Position)

	switch lastTok.Type {
	case token.COLON, token.RARROW:
		items = append(items, getTypesCompletions()...)

	case token.ASSIGN, token.LPAREN, token.COMMA, token.PLUS, token.ASTERISK:
		items = append(items, getValueCompletions(chk, &env, hoveredNode, hoveredType)...)

	case token.LET, token.INT_TYPE, token.BOOL_TYPE, token.STRING_TYPE:
		items = append(items, []protocol.CompletionItem{}...)

	default:
		var res []protocol.CompletionItem
		if hoveredNode != nil {
			res = getValueCompletions(chk, &env, hoveredNode, hoveredType)
		}
		keywords := getKeywords()
		for _, kw := range keywords {
			kind := protocol.CompletionItemKindKeyword
			res = append(res, protocol.CompletionItem{Label: kw, Kind: &kind})
		}
		items = append(items, res...)
	}

	return items, nil
}

func getTypesCompletions() []protocol.CompletionItem {
	types := getTypes()
	items := []protocol.CompletionItem{}
	for _, t := range types {
		kind := protocol.CompletionItemKindClass
		items = append(items, protocol.CompletionItem{Label: t, Kind: &kind})
	}
	return items
}

func getValueCompletions(chk *type_checker.TypeChecker, env *types.TypeEviroment,
	hoveredNode ast.Node, hoveredType types.Type) []protocol.CompletionItem {
	fmt.Fprintf(os.Stderr, "Node hovered %v with type %T ", hoveredNode, hoveredNode)
	// fmt.Fprintf(os.Stderr, "Final map: ")

	// for key, value := range chk.ExpectedTypes {
	// 	fmt.Fprintf(os.Stderr, "Key: %+v, Value: %T\n", key, value)
	// }

	items := []protocol.CompletionItem{}
	for _, name := range env.GetAllNames() {
		symbolInfo, _ := env.Get(name)
		kind := protocol.CompletionItemKindVariable
		var insertText string

		// fmt.Fprintf(os.Stderr, "symb %s: %T %v %v\n", name, symbolInfo.SymbolType, symbolInfo.DeclNode, symbolInfo.Depth)

		if declaredLater(symbolInfo.DeclNode, hoveredNode) {
			continue
		}

		matches := false
		if chk.ExpectedTypes[hoveredNode] == nil {
			matches = true
		} else {
			matches = types.Equals(symbolInfo.SymbolType, chk.ExpectedTypes[hoveredNode])
		}

		// fmt.Fprintf(os.Stderr, "Matches %T symb %T with type %T is %v\n ", symbolInfo.SymbolType,
		// 	chk.ExpectedTypes[hoveredNode], hoveredNode, matches)

		switch innerType := symbolInfo.SymbolType.(type) {
		case *types.FuncType:
			kind = protocol.CompletionItemKindFunction
			matches = types.Equals(innerType.ReturnType, chk.ExpectedTypes[hoveredNode])
			insertText = name + "()"
		case *types.BuiltinFunc:
			kind = protocol.CompletionItemKindFunction
			matches = types.Equals(innerType.ReturnType, chk.ExpectedTypes[hoveredNode])
			insertText = name + "()"
		default:
			if declaredOnSameLine(symbolInfo.DeclNode, hoveredNode) {
				continue
			}
			insertText = name
		}

		item := protocol.CompletionItem{
			Label:      name,
			Kind:       &kind,
			Detail:     &[]string{symbolInfo.SymbolType.Signature()}[0],
			InsertText: &insertText,
		}

		var typePriority string
		if matches {
			if *item.Kind == protocol.CompletionItemKindFunction {
				// item.SortText = &[]string{"002_" + name}[0]
				typePriority = "002"
			} else {
				// item.SortText = &[]string{"001_" + name}[0]
				typePriority = "001"
			}
			detail := "(matches type) " + *item.Detail
			item.Detail = &detail
		} else {
			if *item.Kind == protocol.CompletionItemKindFunction {
				// item.SortText = &[]string{"004_" + name}[0]
				typePriority = "004"
			} else {
				// item.SortText = &[]string{"003_" + name}[0]
				typePriority = "005"
			}
		}

		distStr := fmt.Sprintf("%03d", symbolInfo.Depth)
		sortKey := fmt.Sprintf("%s_%s_%s", typePriority, distStr, name)

		item.SortText = &sortKey

		items = append(items, item)
	}
	return items
}

func getKeywords() []string {
	var res []string
	for key := range token.Keywords {
		res = append(res, key)
	}
	return res
}

func getTypes() []string {
	var res []string
	for key := range token.Types {
		res = append(res, key)
	}
	return res
}

func getPrecedingToken(text string, pos protocol.Position) token.Token {
	lxr := lexer.New(text)
	var lastTok token.Token

	for {
		tok := lxr.NextToken()
		if tok.Type == token.EOF {
			break
		}

		// если токен перед курсором или на его позиции
		if uint32(tok.Start.Line) > pos.Line ||
			(uint32(tok.Start.Line) == pos.Line && uint32(tok.Start.Column) >= pos.Character) {
			break
		}
		lastTok = tok
	}
	return lastTok
}

func declaredLater(firstNode ast.Node, secondNode ast.Node) bool {
	if firstNode.Start().Line > secondNode.End().Line ||
		(firstNode.Start().Line == secondNode.End().Line && firstNode.Start().Column > secondNode.End().Column) {
		return true
	}
	return false
}

func declaredOnSameLine(firstNode ast.Node, secondNode ast.Node) bool {
	if firstNode.Start().Line != secondNode.Start().Line {
		return false
	}
	return true
}
