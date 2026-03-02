package ast

import (
	"bytes"
	"strings"
)

func (prg *Program) String() string {
	var out bytes.Buffer
	for _, statement := range prg.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

// ----- ТИПЫ -----

func (t *SimpleType) String() string { return t.Value }

func (ft *FunctionType) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range ft.ParamsTypes {
		params = append(params, p.String())
	}
	out.WriteString(ft.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") -> ")
	out.WriteString(ft.ReturnType.String())
	return out.String()
}

// ----- ЛИТЕРАЛЫ -----

func (intLit *IntegerLiteral) String() string { return intLit.Token.Literal }

func (boolLit *BooleanLiteral) String() string { return boolLit.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

func (i *Identifier) String() string { return i.Value }

// ----- ИНСТРУКЦИИ -----

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(": " + ls.Type.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (blockStmt *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range blockStmt.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ----- ВЫРАЖЕНИЯ -----

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ifExpr *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ifExpr.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExpr.Consequence.String())
	if ifExpr.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifExpr.Alternative.String())
	}
	return out.String()
}
