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
	defer handlePanic(context)
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

	tokensUpToCursor := getTokensUpToCursor(docText, params.Position)
	var lastTok token.Token
	if len(tokensUpToCursor) > 0 {
		lastTok = tokensUpToCursor[len(tokensUpToCursor)-1]
	}

	if isExpectedTypeContext(tokensUpToCursor) {
		items = append(items, getTypesCompletions()...)
	} else {
		switch lastTok.Type {
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
	}

	return items, nil
}

func getTypesCompletions() []protocol.CompletionItem {
	types := getTypes()
	items := []protocol.CompletionItem{}
	for _, t := range types {
		kind := protocol.CompletionItemKindClass
		item := protocol.CompletionItem{Label: t, Kind: &kind}
		sortKey := fmt.Sprintf("%s_%s", "001_", t)
		item.SortText = &sortKey
		items = append(items, item)
	}

	fnKind := protocol.CompletionItemKindInterface
	fnPriority := "009"
	items = append(items, protocol.CompletionItem{
		Label:            "fn(args) -> type",
		Kind:             &fnKind,
		InsertText:       &[]string{"fn($1) -> $2"}[0],
		SortText:         &fnPriority,
		InsertTextFormat: &[]protocol.InsertTextFormat{protocol.InsertTextFormatSnippet}[0],
	})

	arrKind := protocol.CompletionItemKindUnit
	arrPriority := "002"
	items = append(items, protocol.CompletionItem{
		Label:            "[] (array type)",
		Kind:             &arrKind,
		InsertText:       &[]string{"[$1]"}[0],
		SortText:         &arrPriority,
		InsertTextFormat: &[]protocol.InsertTextFormat{protocol.InsertTextFormatSnippet}[0],
	})

	hashMapKind := protocol.CompletionItemKindUnit
	hashMapPriority := "003"
	items = append(items, protocol.CompletionItem{
		Label:            "{key : value} (hash map type)",
		Kind:             &hashMapKind,
		InsertText:       &[]string{"{$1 : $2}"}[0],
		SortText:         &hashMapPriority,
		InsertTextFormat: &[]protocol.InsertTextFormat{protocol.InsertTextFormatSnippet}[0],
	})

	return items
}

func getValueCompletions(chk *type_checker.TypeChecker, env *types.TypeEviroment,
	hoveredNode ast.Node, hoveredType types.Type) []protocol.CompletionItem {
	fmt.Fprintf(os.Stderr, "Node hovered %v with type %T expected type: %T",
		hoveredNode, hoveredNode, chk.ExpectedTypes[hoveredNode])
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
			fmt.Fprintf(os.Stderr, "symb %T v. %T is %+v\n", symbolInfo.SymbolType, chk.ExpectedTypes[hoveredNode], matches)
		}

		// fmt.Fprintf(os.Stderr, "Matches %T symb %T with type %T is %v\n ", symbolInfo.SymbolType,
		// 	chk.ExpectedTypes[hoveredNode], hoveredNode, matches)

		switch innerType := symbolInfo.SymbolType.(type) {
		case *types.FuncType:
			kind = protocol.CompletionItemKindFunction
			if !matches {
				matches = types.Equals(innerType.ReturnType, chk.ExpectedTypes[hoveredNode])
				insertText = name + "($1)"
			}
		case *types.BuiltinFunc:
			kind = protocol.CompletionItemKindFunction
			if !matches {
				matches = types.Equals(innerType.ReturnType, chk.ExpectedTypes[hoveredNode])
				insertText = name + "($1)"
			}
		default:
			if declaredOnSameLine(symbolInfo.DeclNode, hoveredNode) {
				continue
			}
			insertText = name
		}

		item := protocol.CompletionItem{
			Label:            name,
			Kind:             &kind,
			Detail:           &[]string{symbolInfo.SymbolType.Signature()}[0],
			InsertText:       &insertText,
			InsertTextFormat: &[]protocol.InsertTextFormat{protocol.InsertTextFormatSnippet}[0],
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

func getTokensUpToCursor(text string, pos protocol.Position) []token.Token {
	lxr := lexer.New(text)
	var tokens []token.Token

	for {
		tok := lxr.NextToken()
		if tok.Type == token.EOF {
			break
		}

		// остановится, если токен начался после или на позиции курсора
		if uint32(tok.Start.Line) > pos.Line ||
			(uint32(tok.Start.Line) == pos.Line && uint32(tok.Start.Column) >= pos.Character) {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

func isExpectedTypeContext(tokens []token.Token) bool {
	if len(tokens) == 0 {
		return false
	}

	lastTok := tokens[len(tokens)-1]

	if lastTok.Type == token.RARROW {
		return true
	}

	if lastTok.Type == token.COLON {
		if len(tokens) >= 3 {
			prev1 := tokens[len(tokens)-2] // токен до двоеточия
			prev2 := tokens[len(tokens)-3]

			// контекст let
			if prev1.Type == token.IDENT && prev2.Type == token.LET {
				return true
			}

			// контекст параметров функции
			if prev1.Type == token.IDENT && (prev2.Type == token.LPAREN || prev2.Type == token.COMMA) {
				return true
			}
		}

		// контекст типа HashMap
		if len(tokens) >= 2 {
			prev1 := tokens[len(tokens)-2]

			if prev1.Type == token.INT_TYPE || prev1.Type == token.STRING_TYPE || prev1.Type == token.BOOL_TYPE {
				return true
			}
		}

		return false
	}

	return false
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
