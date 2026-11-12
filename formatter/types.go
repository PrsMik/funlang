package formatter

import "funlang/ast"

func (fmtr *Formatter) formatType(node ast.ExpressionNode) {
	if node == nil {
		return
	}
	switch t := node.(type) {
	case *ast.SimpleType:
		fmtr.out.WriteString(t.Value)
	case *ast.ArrayType:
		fmtr.out.WriteString("[]")
		fmtr.formatType(t.ElementsType)
	case *ast.HashMapType:
		fmtr.out.WriteString("{")
		fmtr.formatType(t.KeyType)
		fmtr.out.WriteString("}")
		fmtr.formatType(t.ElementType)
	case *ast.FunctionType:
		fmtr.out.WriteString("fn(")
		for i, pt := range t.ParamsTypes {
			if t.Parameters[i] != nil {
				fmtr.formatExpression(t.Parameters[i])
				fmtr.out.WriteString(": ")
			}
			fmtr.formatType(pt)
			if i < len(t.ParamsTypes)-1 {
				fmtr.out.WriteString(", ")
			}
		}
		fmtr.out.WriteString(")")
		if t.ReturnType != nil {
			fmtr.out.WriteString(" -> ")
			fmtr.formatType(t.ReturnType)
		}
	default:
		fmtr.out.WriteString(t.String())
	}
}
