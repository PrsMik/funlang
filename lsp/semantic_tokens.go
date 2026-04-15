package lsp

import (
	"funlang/ast"
	"funlang/types"
	"sort"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type semToken struct {
	line      uint32
	character uint32
	length    uint32
	tokenType uint32
	modifiers uint32
}

const (
	tokenTypeVariable  = 0
	tokenTypeFunction  = 1
	tokenTypeParameter = 2
	tokenTypeType      = 3
)

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	defer handlePanic(context)

	info, ok := documentStates[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	var rawTokens []semToken

	for node, nodeType := range info.TypesInfo {
		ident, isIdent := node.(*ast.Identifier)
		if !isIdent {
			continue
		}

		start := ident.Start()
		length := uint32(len(ident.Value))

		tokenType := uint32(tokenTypeVariable)

		switch nodeType.(type) {
		case *types.FuncType, *types.BuiltinFunc:
			tokenType = tokenTypeFunction
		default:
			tokenType = tokenTypeVariable
		}

		rawTokens = append(rawTokens, semToken{
			line:      uint32(start.Line),
			character: uint32(start.Column),
			length:    length,
			tokenType: tokenType,
			modifiers: 0,
		})
	}

	for typeNode := range info.TypeNodes {
		start := typeNode.Start()
		rawTokens = append(rawTokens, semToken{
			line:      uint32(start.Line),
			character: uint32(start.Column),
			length:    uint32(len(typeNode.TokenLiteral())),
			tokenType: uint32(tokenTypeType),
			modifiers: 0,
		})
	}

	// токены ОБЯЗАТЕЛЬНО должны быть отсортированы сверху вниз, слева направо
	sort.Slice(rawTokens, func(i, j int) bool {
		if rawTokens[i].line == rawTokens[j].line {
			return rawTokens[i].character < rawTokens[j].character
		}
		return rawTokens[i].line < rawTokens[j].line
	})

	// кодирование дельт
	var data []protocol.UInteger
	var prevLine, prevChar uint32

	for _, t := range rawTokens {
		deltaLine := t.line - prevLine
		deltaChar := t.character

		if deltaLine == 0 {
			deltaChar = t.character - prevChar
		}

		data = append(data,
			protocol.UInteger(deltaLine),
			protocol.UInteger(deltaChar),
			protocol.UInteger(t.length),
			protocol.UInteger(t.tokenType),
			protocol.UInteger(t.modifiers),
		)

		prevLine = t.line
		prevChar = t.character
	}

	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}
