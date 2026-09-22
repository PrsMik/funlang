package types

import (
	"fmt"
	"funlang/ast"
)

type Type interface {
	isType()
	Signature() string
}

type HashableType interface {
	isHashable()
}

type PrintableType interface {
	isPrintable()
}

type ComparableType interface {
	isComparable()
}

// ===========================================================
// ОСОБЫЕ ТИПЫ
// ===========================================================

type TypeType struct {
	Underlying Type
}

func (t *TypeType) isType() {}
func (t *TypeType) Signature() string {
	if t == MetaType {
		return "<metatype>"
	}
	if t.Underlying != nil {
		return fmt.Sprintf("<type(%s)>", t.Underlying.Signature())
	} else {
		return "<type>"
	}
}
func (t *TypeType) isComparable() {}

type IllegalType struct{}

func (t *IllegalType) isType()           {}
func (t *IllegalType) Signature() string { return "<illegal/none type>" }
func (t *IllegalType) isPrintable()      {}

type NullType struct{}

func (n *NullType) isType()           {}
func (n *NullType) Signature() string { return "<null>" }
func (t *NullType) isPrintable()      {}

// ===========================================================
// ПРОСТЫЕ ТИПЫ
// ===========================================================

type IntType struct{}

func (t *IntType) isType()           {}
func (t *IntType) Signature() string { return "<int>" }
func (t *IntType) isHashable()       {}
func (t *IntType) isPrintable()      {}
func (t *IntType) isComparable()     {}

type BoolType struct{}

func (t *BoolType) isType()           {}
func (t *BoolType) Signature() string { return "<bool>" }
func (t *BoolType) isHashable()       {}
func (t *BoolType) isPrintable()      {}
func (t *BoolType) isComparable()     {}

type StringType struct{}

func (t *StringType) isType()           {}
func (t *StringType) Signature() string { return "<string>" }
func (t *StringType) isHashable()       {}
func (t *StringType) isPrintable()      {}
func (t *StringType) isComparable()     {}

type ArrayType struct {
	ElementsType Type
}

func (t *ArrayType) isType() {}
func (t *ArrayType) Signature() string {
	if t.ElementsType != nil {
		return "<[" + t.ElementsType.Signature() + "]>"
	}
	return "<[]>"
}
func (t *ArrayType) isPrintable()  {}
func (t *ArrayType) isComparable() {}

type HashMapType struct {
	KeyType     Type
	ElementType Type
}

func (bt *HashMapType) isType() {}
func (bt *HashMapType) Signature() string {
	if bt.KeyType != nil && bt.ElementType != nil {
		return "<{" + bt.KeyType.Signature() + ":" + bt.ElementType.Signature() + "}>"
	}

	return "<{}>"
}
func (t *HashMapType) isPrintable()  {}
func (t *HashMapType) isComparable() {}

// ===========================================================
// ФУНКЦИОНАЛЬНЫЕ ТИПЫ
// ===========================================================

type BuiltinFunc struct {
	CheckFunc  func(args []Type) (Type, error)
	ReturnType Type
}

func (bt *BuiltinFunc) isType()           {}
func (bt *BuiltinFunc) Signature() string { return "<builtin_func>" }

type FuncParam struct {
	Name string
	Type Type
}

type FuncType struct {
	Params     []FuncParam
	ReturnType Type
}

func (t *FuncType) isType() {}
func (t *FuncType) Signature() string {
	if len(t.Params) == 0 {
		return "<fn() -> " + t.ReturnType.Signature() + ">"
	}

	res := "<fn("

	for idx, param := range t.Params {
		res += param.Name + ": "
		res += param.Type.Signature()
		if idx < len(t.Params)-1 {
			res += ", "
		}
	}

	res += ") -> " + t.ReturnType.Signature() + ">"
	return res
}

// ===========================================================
// КОМПЛЕКСНЫЕ ПОЛЬЗОВАТЕЛЬСКИЕ ТИПЫ
// ===========================================================

type InterfaceType struct {
	Name   string
	Fields map[string]Type
}

func (t *InterfaceType) isType()           {}
func (t *InterfaceType) Signature() string { return fmt.Sprintf("<interface %s>", t.Name) }
func (t *InterfaceType) isComparable()     {}

type ImplType struct {
	Name                string
	ImplementedType     Type
	Fields              map[string]Type
	FieldsDefaultValues map[string]ast.ExpressionNode
}

func (t *ImplType) isType() {}
func (t *ImplType) Signature() string {
	res := fmt.Sprintf("<impl %s>", t.Name)
	if t.ImplementedType != nil {
		res = fmt.Sprintf("<impl %s (%s)>", t.Name, t.ImplementedType.Signature())
	}
	return res
}
func (t *ImplType) isComparable() {}

type UnionType struct {
	Left  Type
	Right Type
}

func (t *UnionType) isType() {}
func (t *UnionType) Signature() string {
	return fmt.Sprintf("<union %s | %s>", t.Left.Signature(), t.Right.Signature())
}
func (t *UnionType) isComparable() {}

type IntersectionType struct {
	Left  Type
	Right Type
}

func (t *IntersectionType) isType() {}
func (t *IntersectionType) Signature() string {
	return fmt.Sprintf("<intersection %s & %s>", t.Left.Signature(), t.Right.Signature())
}
func (t *IntersectionType) isComparable() {}
