package lsp

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
	"os"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentSignatureHelp(context *glsp.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	defer handlePanic(context)

	chk, ok := documentStates[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}

	pos := params.Position

	var targetCall *ast.CallExpression
	var minLen int = 9999999

	for node := range chk.TypesInfo {
		callExpr, isCall := node.(*ast.CallExpression)
		if !isCall {
			continue
		}

		start := callExpr.Function.End()
		end := callExpr.SemiToken.End
		if callExpr.SemiToken.Type == 0 {
			end = callExpr.End()
		}

		if isPosInside(pos, start, end) {
			length := (end.Line-start.Line)*1000 + (end.Column - start.Column)
			if length < minLen {
				minLen = length
				targetCall = callExpr
			}
		}
	}

	if targetCall == nil {
		return nil, nil
	}

	funcType, ok := chk.TypesInfo[targetCall.Function]
	fmt.Fprintf(os.Stderr, "Final map: ")

	for key, value := range chk.TypesInfo {
		fmt.Printf("Key: %+v, Value: %T\n", key, value)
	}

	if !ok {
		return nil, nil
	}

	var signatureInfo protocol.SignatureInformation
	var paramsInfo []protocol.ParameterInformation

	switch ft := funcType.(type) {
	case *types.FuncType:
		fmt.Printf("Server Fn type: %v\n", ft)
		var paramLabels []string
		for _, param := range ft.Params {
			label := param.Name + ": " + param.Type.Signature()
			paramLabels = append(paramLabels, label)
			paramsInfo = append(paramsInfo, protocol.ParameterInformation{
				Label: label,
			})
		}

		funcName := targetCall.Function.TokenLiteral()
		label := funcName + "(" + strings.Join(paramLabels, ", ") + ") -> " + ft.ReturnType.Signature()

		signatureInfo = protocol.SignatureInformation{
			Label:      label,
			Parameters: paramsInfo,
		}

	case *types.BuiltinFunc:
		signatureInfo = protocol.SignatureInformation{
			Label: targetCall.Function.TokenLiteral() + "(...) -> " + ft.ReturnType.Signature(),
		}
	default:
		return nil, nil
	}

	activeParam := calculateActiveParameter(targetCall, pos)

	signature := protocol.UInteger(0)
	parameter := protocol.UInteger(uint32(activeParam))

	return &protocol.SignatureHelp{
		Signatures:      []protocol.SignatureInformation{signatureInfo},
		ActiveSignature: &signature,
		ActiveParameter: &parameter,
	}, nil
}

func calculateActiveParameter(callExpr *ast.CallExpression, pos protocol.Position) int {
	if len(callExpr.Arguments) == 0 {
		return 0
	}

	for i, arg := range callExpr.Arguments {
		end := arg.End()
		if uint32(end.Line) > pos.Line || (uint32(end.Line) == pos.Line && uint32(end.Column) >= pos.Character) {
			return i
		}
	}

	return len(callExpr.Arguments)
}
