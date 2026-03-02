package parser

import (
	"fmt"
	"funlang/ast"
	"funlang/lexer"
	"funlang/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
			let x: int = 5;
			let y: int = 10;
			let foobar: bool = true;
			`
	lxr := lexer.New(input)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       token.TokenType
	}{
		{"x", token.INT_TYPE},
		{"y", token.INT_TYPE},
		{"foobar", token.BOOL_TYPE},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statementNode ast.StatementNode, name string, tknType token.TokenType) bool {
	if statementNode.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statementNode.TokenLiteral())
		return false
	}

	letStmt, ok := statementNode.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statementNode)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	simpleType, ok := letStmt.Type.(*ast.SimpleType)
	if !ok {
		t.Errorf("letStmt.Type not *ast.SimpleType. got=%T", letStmt.Type)
		return false
	}

	if simpleType.Token.Type != tknType {
		want, _ := token.LookupString(tknType)
		got, _ := token.LookupString(simpleType.Token.Type)
		t.Errorf("type not correct. expected=%s, got=%s", want, got)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return 1 + 1;
			return add(1, 1);
			`
	lxr := lexer.New(input)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "return foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Value)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "let x: int = 5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Value)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "let x: bool = true;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Value.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Value)
	}
	if ident.Value != true {
		t.Errorf("ident.Value not %v. got=%v", true, ident.Value)
	}
	if ident.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "true",
			ident.TokenLiteral())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `let myFunc: fn(int, int) -> int = fn(x, y) { return x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Value, "x", "+", "y")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"let x: int = -15;", "-", 15},
		{"let x: bool = !false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Value)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"let x: int = 5 + 5;", 5, "+", 5},
		{"let x: int = 5 - 5;", 5, "-", 5},
		{"let x: int = 5 * 5;", 5, "*", 5},
		{"let x: int = 5 / 5;", 5, "/", 5},
		{"let x: int = 5 > 5;", 5, ">", 5},
		{"let x: int = 5 < 5;", 5, "<", 5},
		{"let x: int = false == false;", false, "==", false},
		{"let x: int = true != true;", true, "!=", true},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Value)
		}
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestParsingIfExpression(t *testing.T) {
	input := "let x: bool = if (x < y) { let y: int = 5; return x + y; } else { return 5; };"

	lxr := lexer.New(input)
	prs := New(lxr)
	prg := prs.ParseProgram()
	checkParserErrors(t, prs)

	if len(prg.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(prg.Statements))
	}

	stmt, ok := prg.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", prg.Statements[0])
	}

	exp, ok := stmt.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Value)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 2 {
		t.Errorf("consequence is not 2 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	consequence1, ok1 := exp.Consequence.Statements[0].(*ast.LetStatement)
	if !ok1 {
		t.Fatalf("Statements[0] is not ast.ReturnStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	consequence2, ok2 := exp.Consequence.Statements[1].(*ast.ReturnStatement)
	if !ok2 {
		t.Fatalf("Statements[0] is not ast.ReturnStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence1.Value, 5) {
		return
	}

	if !testInfixExpression(t, consequence2.Value, "x", "+", "y") {
		return
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{
		"let x: bool = a + b * c == 3 * 1 + 4 * 5 == true;",
		"let x: bool = (((a + (b * c)) == ((3 * 1) + (4 * 5))) == true);",
	},
		{
			"let x: int = (5 + 5) * 2;",
			"let x: int = ((5 + 5) * 2);",
		},
	}
	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prs := New(lxr)
		program := prs.ParseProgram()
		checkParserErrors(t, prs)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.ExpressionNode, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.ExpressionNode, expected interface{}) bool {
	switch val := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, val)
	case bool:
		return testBooleanLiteral(t, exp, val)
	case string:
		return testIdentifierLiteral(t, exp, val)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.ExpressionNode, value int) bool {
	integ, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integerLiteral not *ast.IntegerLiteral. got=%T", integerLiteral)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.ExpressionNode, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}

func testIdentifierLiteral(t *testing.T, expression ast.ExpressionNode, value string) bool {
	ident, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not ast.Identifier; got: %T", expression)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value is not %s; got :%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is not %s; got :%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, prs *Parser) {
	errors := prs.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
