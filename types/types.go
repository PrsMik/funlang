package types

type Type interface {
	isType()
	Signature() string
}

type IllegalType struct{}

func (t *IllegalType) isType()           {}
func (t *IllegalType) Signature() string { return "none" }

type NullType struct{}

func (n *NullType) isType()           {}
func (n *NullType) Signature() string { return "<null>" }

type IntType struct{}

func (t *IntType) isType()           {}
func (t *IntType) Signature() string { return "<int>" }

type BoolType struct{}

func (t *BoolType) isType()           {}
func (t *BoolType) Signature() string { return "<bool>" }

type FuncType struct {
	ParamsTypes []Type
	ReturnType  Type
}

func (t *FuncType) isType() {}
func (t *FuncType) Signature() string {
	if len(t.ParamsTypes) == 0 {
		return "<fn() " + t.ReturnType.Signature() + ">"
	}

	res := "<fn("

	for idx, tp := range t.ParamsTypes {
		res += tp.Signature()
		if idx < len(t.ParamsTypes)-1 {
			res += ", "
		}
	}

	res += ") -> " + t.ReturnType.Signature() + ">"
	return res
}

func Equals(rawLeftType, rawRightType Type) bool {
	if rawLeftType == rawRightType {
		return true
	}

	switch leftType := rawLeftType.(type) {
	case *IntType:
		_, ok := rawRightType.(*IntType)
		return ok
	case *BoolType:
		_, ok := rawRightType.(*BoolType)
		return ok
	case *FuncType:
		rightType, ok := rawRightType.(*FuncType)
		if !ok {
			return false
		}

		if len(leftType.ParamsTypes) != len(rawRightType.(*FuncType).ParamsTypes) {
			return false
		}

		for i := range leftType.ParamsTypes {
			if !Equals(leftType.ParamsTypes[i], rightType.ParamsTypes[i]) {
				return false
			}
		}

		return Equals(leftType.ReturnType, rightType.ReturnType)
	default:
		return true
	}
}
