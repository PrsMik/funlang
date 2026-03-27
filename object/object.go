package object

import (
	"funlang/ast"
	"hash/fnv"
)

type ObjectType int

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

const (
	ILLEGAL ObjectType = iota
	NULL_OBJ
	INTEGER_OBJ
	BOOLEAN_OBJ
	STRING_OBJ
	ARRAY_OBJ
	HASH_OBJ
	RETURN_VALUE_OBJ
	FUNCTION_OBJ
	TAIL_CALL_OBJ
	BUILTIN_OBJ
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

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

type HashKey struct {
	Type  ObjectType
	Value uint
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint(i.Value)}
}

func (b *Boolean) HashKey() HashKey {
	var value uint

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (str *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(str.Value))

	return HashKey{Type: str.Type(), Value: uint(h.Sum64())}
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashMap struct {
	Pairs map[HashKey]HashPair
}

func (hm *HashMap) Type() ObjectType { return HASH_OBJ }

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

type TailCall struct {
	Function  Object
	Arguments []Object
}

func (f *TailCall) Type() ObjectType { return TAIL_CALL_OBJ }

type BuiltinFunction func(...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (e *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
