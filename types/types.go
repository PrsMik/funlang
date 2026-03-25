package types

type Type interface {
	isType()
	Signature() string
}

type HashableType interface {
	isHashable()
}

type IllegalType struct{}

func (t *IllegalType) isType()           {}
func (t *IllegalType) Signature() string { return "<none>" }

type NullType struct{}

func (n *NullType) isType()           {}
func (n *NullType) Signature() string { return "<null>" }

type IntType struct{}

func (t *IntType) isType()           {}
func (t *IntType) Signature() string { return "<int>" }
func (t *IntType) isHashable()       {}

type BoolType struct{}

func (t *BoolType) isType()           {}
func (t *BoolType) Signature() string { return "<bool>" }
func (t *BoolType) isHashable()       {}

type StringType struct{}

func (t *StringType) isType()           {}
func (t *StringType) Signature() string { return "<string>" }
func (t *StringType) isHashable()       {}

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

type BuiltinFunc struct {
	CheckFunc func(args []Type) (Type, error)
}

func (bt *BuiltinFunc) isType()           {}
func (bt *BuiltinFunc) Signature() string { return "<builtin_func>" }

type FuncType struct {
	ParamsTypes []Type
	ReturnType  Type
}

func (t *FuncType) isType() {}
func (t *FuncType) Signature() string {
	if len(t.ParamsTypes) == 0 {
		return "<fn() -> " + t.ReturnType.Signature() + ">"
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
	case *StringType:
		_, ok := rawRightType.(*StringType)
		return ok
	case *ArrayType:
		rightType, ok := rawRightType.(*ArrayType)
		if !ok {
			return false
		}

		if rightType.ElementsType == nil || leftType.ElementsType == nil {
			return true
		}

		return Equals(leftType.ElementsType, rightType.ElementsType)
	case *HashMapType:
		rightType, ok := rawRightType.(*HashMapType)
		if !ok {
			return false
		}

		if (rightType.KeyType == nil && rightType.ElementType == nil) ||
			(leftType.KeyType == nil && leftType.ElementType == nil) {
			return true
		}

		return Equals(leftType.KeyType, rightType.KeyType) && Equals(leftType.ElementType, rightType.ElementType)
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
		return false
	}
}
