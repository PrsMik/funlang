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

func (typeType *TypeType) String() string { return typeType.Value }

func (simpType *SimpleType) String() string { return simpType.Value }

func (arrType *ArrayType) String() string { return "[" + arrType.ElementsType.String() + "]" }

func (hashMapType *HashMapType) String() string {
	return "{" + hashMapType.KeyType.String() + " : " + hashMapType.ElementType.String() + "}"
}

func (funcType *FunctionType) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range funcType.ParamsTypes {
		params = append(params, p.String())
	}
	out.WriteString(funcType.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") -> ")
	out.WriteString(funcType.ReturnType.String())
	return out.String()
}

// ----- ЛИТЕРАЛЫ -----

func (intLit *IntegerLiteral) String() string { return intLit.Token.Literal }

func (boolLit *BooleanLiteral) String() string { return boolLit.Token.Literal }

func (strLit *StringLiteral) String() string { return strLit.Token.Literal }

func (arrLit *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, elem := range arrLit.Elements {
		elements = append(elements, elem.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (hashMapLiteral *HashMapLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}

	for key, value := range hashMapLiteral.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (funcLit *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range funcLit.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(funcLit.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(funcLit.Body.String())
	return out.String()
}

func (iface *InterfaceLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(iface.TokenLiteral())
	out.WriteString("{")
	out.WriteString(iface.Body.String())
	out.WriteString("}")
	return out.String()
}

func (imp *ImplLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(imp.TokenLiteral())
	out.WriteString("(")
	out.WriteString(imp.TargetType.String())
	out.WriteString(")")
	out.WriteString(imp.Body.String())
	return out.String()
}

func (ident *Identifier) String() string { return ident.Value }

// ----- ИНСТРУКЦИИ -----

func (letStmt *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(letStmt.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(letStmt.Name.String())
	out.WriteString(": ")
	out.WriteString(letStmt.Type.String())
	out.WriteString(" = ")
	if letStmt.Value != nil {
		out.WriteString(letStmt.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (returnStmt *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(returnStmt.TokenLiteral())
	out.WriteString(" ")
	if returnStmt.Value != nil {
		out.WriteString(returnStmt.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (blockStmt *BlockStatement) String() string {
	var out bytes.Buffer
	// out.WriteString("{\n")
	for _, s := range blockStmt.Statements {
		out.WriteString(s.String())
	}
	// out.WriteString("}")
	return out.String()
}

// ----- ВЫРАЖЕНИЯ -----

func (callExpr *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range callExpr.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(callExpr.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

func (prefixExpr *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefixExpr.Operator)
	out.WriteString(prefixExpr.Right.String())
	out.WriteString(")")
	return out.String()
}

func (infixExpr *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infixExpr.Left.String())
	out.WriteString(" ")
	out.WriteString(infixExpr.Operator)
	out.WriteString(" ")
	out.WriteString(infixExpr.Right.String())
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
		out.WriteString(" else ")
		out.WriteString(ifExpr.Alternative.String())
	}
	return out.String()
}

func (caseBlock *CaseBlock) String() string {
	var out bytes.Buffer
	out.WriteString("case")
	out.WriteString(caseBlock.Pattern.String())
	out.WriteString("{")
	out.WriteString(caseBlock.Body.String())
	out.WriteString("}")
	return out.String()
}

func (swExp *SwitchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("switch")
	out.WriteString(swExp.Value.String())
	out.WriteString("{")
	for _, caseBlock := range swExp.Cases {
		out.WriteString(caseBlock.String())
	}
	out.WriteString("default")
	out.WriteString(swExp.Default.String())
	out.WriteString("}")
	return out.String()
}

func (indExpr *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(indExpr.Left.String())
	out.WriteString("[")
	out.WriteString(indExpr.Index.String())
	out.WriteString("])")
	return out.String()
}

func (membExpr *MemberAccessExpression) String() string {
	var out bytes.Buffer
	out.WriteString(membExpr.Left.String())
	out.WriteString(".")
	out.WriteString(membExpr.Property.String())
	return out.String()
}

func (instantiationExpr *ImplInstantiationExpression) String() string {
	var out bytes.Buffer
	out.WriteString(instantiationExpr.Left.String())
	out.WriteString(".")
	pairs := []string{}

	for key, value := range instantiationExpr.Fields {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (tae *TypeAccessExpression) String() string {
	var out bytes.Buffer
	out.WriteString(tae.Left.String())
	out.WriteString(".(type)")
	return out.String()
}
