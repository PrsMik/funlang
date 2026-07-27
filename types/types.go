package types

import "fmt"

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

type TypeType struct {
	Underlying Type
}

func (t *TypeType) isType() {}
func (t *TypeType) Signature() string {
	if t.Underlying != nil {
		return fmt.Sprintf("<type(%s)>", t.Underlying.Signature())
	} else {
		return "<type>"
	}
}

type IllegalType struct{}

func (t *IllegalType) isType()           {}
func (t *IllegalType) Signature() string { return "<none>" }
func (t *IllegalType) isPrintable()      {}

type NullType struct{}

func (n *NullType) isType()           {}
func (n *NullType) Signature() string { return "<null>" }
func (t *NullType) isPrintable()      {}

type IntType struct{}

func (t *IntType) isType()           {}
func (t *IntType) Signature() string { return "<int>" }
func (t *IntType) isHashable()       {}
func (t *IntType) isPrintable()      {}

type BoolType struct{}

func (t *BoolType) isType()           {}
func (t *BoolType) Signature() string { return "<bool>" }
func (t *BoolType) isHashable()       {}
func (t *BoolType) isPrintable()      {}

type StringType struct{}

func (t *StringType) isType()           {}
func (t *StringType) Signature() string { return "<string>" }
func (t *StringType) isHashable()       {}
func (t *StringType) isPrintable()      {}

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
func (t *ArrayType) isPrintable() {}

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
func (t *HashMapType) isPrintable() {}

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

type InterfaceType struct {
	Name   string
	Fields map[string]Type
}

func (t *InterfaceType) isType()           {}
func (t *InterfaceType) Signature() string { return fmt.Sprintf("<inteface %s>", t.Name) }

type ImplType struct {
	Name            string
	ImplementedType Type
	Fields          map[string]Type
}

func (t *ImplType) isType() {}
func (t *ImplType) Signature() string {
	res := fmt.Sprintf("<impl %s (%s)>", t.Name, t.ImplementedType.Signature())
	if t.ImplementedType != nil {
		res = fmt.Sprintf("<impl %s (%s)>", t.Name, t.ImplementedType.Signature())
	}
	return res
}

type UnionType struct {
	Left  Type
	Right Type
}

func (t *UnionType) isType() {}
func (t *UnionType) Signature() string {
	return fmt.Sprintf("<union %s | %s>", t.Left.Signature(), t.Right.Signature())
}

type IntersectionType struct {
	Left  Type
	Right Type
}

func (t *IntersectionType) isType() {}
func (t *IntersectionType) Signature() string {
	return fmt.Sprintf("<intersection %s & %s>", t.Left.Signature(), t.Right.Signature())
}
