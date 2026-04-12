package lsp

import (
	"fmt"
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
	"funlang/type_checker"
	"funlang/types"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func getRenameTarget(chk *type_checker.TypeChecker, pos protocol.Position) (ast.Node, error) {
	var targetDecl ast.Node
	var targetUsage ast.Node
	var minLen int = 9999999

	for usageNode, declNode := range chk.Definitions {
		if isPosInside(pos, usageNode.Start(), usageNode.End()) {
			length := (usageNode.End().Line-usageNode.Start().Line)*1000 + (usageNode.End().Column - usageNode.Start().Column)
			if length < minLen {
				minLen = length
				targetUsage = usageNode
				targetDecl = declNode
			}
		}
	}

	if targetDecl == nil {
		for _, declNode := range chk.Definitions {
			if declNode != nil && isPosInside(pos, declNode.Start(), declNode.End()) {
				targetDecl = declNode
				targetUsage = declNode
				break
			}
		}
	}

	if targetDecl == nil {
		return nil, fmt.Errorf("symbol not found")
	}

	if tp, ok := chk.TypesInfo[targetUsage]; ok {
		if _, isBuiltin := tp.(*types.BuiltinFunc); isBuiltin {
			return nil, fmt.Errorf("cannot rename builtin '%s'", targetUsage.TokenLiteral())
		}
	}

	start := targetDecl.Start()
	if start.Line == -1 && start.Column == -1 {
		return nil, fmt.Errorf("cannot rename system identifier")
	}

	return targetDecl, nil
}

func textDocumentRename(context *glsp.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	defer handlePanic(context)

	chk, ok := documentStates[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	targetDeclNode, err := getRenameTarget(chk, params.Position)
	if err != nil {
		return nil, err
	}

	if !isValidIdentifier(params.NewName) {
		return nil, fmt.Errorf("'%s' is not a valid identifier name", params.NewName)
	}

	uniqueEdits := make(map[string]protocol.TextEdit)

	addEdit := func(node ast.Node) {
		if node == nil {
			return
		}
		start := node.Start()
		key := fmt.Sprintf("%d:%d", start.Line, start.Column)
		uniqueEdits[key] = protocol.TextEdit{
			Range:   createLspRange(node.Start(), node.End()),
			NewText: params.NewName,
		}
	}

	addEdit(targetDeclNode)

	for usageNode, declNode := range chk.Definitions {
		if declNode == targetDeclNode {
			addEdit(usageNode)
		}
	}

	var edits []protocol.TextEdit
	for _, edit := range uniqueEdits {
		edits = append(edits, edit)
	}

	changes := make(map[protocol.DocumentUri][]protocol.TextEdit)
	changes[params.TextDocument.URI] = edits

	return &protocol.WorkspaceEdit{
		Changes: changes,
	}, nil
}

func isValidIdentifier(name string) bool {
	lxr := lexer.New(name)

	tok := lxr.NextToken()

	if tok.Type != token.IDENT {
		return false
	}

	nextTok := lxr.NextToken()

	return nextTok.Type == token.EOF
}
