package types

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

		if len(leftType.Params) != len(rawRightType.(*FuncType).Params) {
			return false
		}

		for i := range leftType.Params {
			if !Equals(leftType.Params[i].Type, rightType.Params[i].Type) {
				return false
			}
		}

		return Equals(leftType.ReturnType, rightType.ReturnType)
	default:
		return false
	}
}
