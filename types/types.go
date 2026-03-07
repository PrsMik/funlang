package types

type Type interface {
	isType()
}

type IntType struct{}

func (t *IntType) isType() {}

type BoolType struct{}

func (t *BoolType) isType() {}

type FuncType struct {
	ParamsTypes []Type
	ReturnType  Type
}

func (t *FuncType) isType() {}
