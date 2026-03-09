package object

import "funlang/ast"

type ObjectType int

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	ILLEGAL ObjectType = iota
	NULL_OBJ
	INTEGER_OBJ
	BOOLEAN_OBJ
	RETURN_VALUE_OBJ
	FUNCTION_OBJ
	ERROR_OBJ
)

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
