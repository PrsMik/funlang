package types

// !!!ПОРЯДОК ВАЖЕН!!!
// проверяет можно ли expected (левому) присвоить значение с типом actual (правое)
func IsAssignable(expected, actual Type) bool {
	if EqualsSoft(expected, actual) {
		return true
	}

	if actualUnion, ok := actual.(*UnionType); ok {
		return IsAssignable(expected, actualUnion.Left) && IsAssignable(expected, actualUnion.Right)
	}

	if expectedUnion, ok := expected.(*UnionType); ok {
		return IsAssignable(expectedUnion.Left, actual) || IsAssignable(expectedUnion.Right, actual)
	}

	if expectedIntersection, ok := expected.(*IntersectionType); ok {
		return IsAssignable(expectedIntersection.Left, actual) && IsAssignable(expectedIntersection.Right, actual)
	}

	if actualIntersection, ok := actual.(*IntersectionType); ok {
		return IsAssignable(expected, actualIntersection.Left) || IsAssignable(expected, actualIntersection.Right)
	}

	if actualImpl, ok := actual.(*ImplType); ok {
		return IsAssignable(expected, actualImpl.ImplementedType)
	}

	if actualTypeType, ok := actual.(*TypeType); ok {
		if actualUnderImpl, ok := actualTypeType.Underlying.(*ImplType); ok {
			return IsAssignable(expected, actualUnderImpl)
		}
	}

	return false
}

// !!!ПОРЯДОК ВАЖЕН!!! ЛЕВОЕ - ТО ЧЕМУ ПРИСВАИВАЕТСЯ ПРАВОЕ
func Equals(rawLeftType, rawRightType Type) bool {
	if rawLeftType == rawRightType {
		return true
	}

	switch leftType := rawLeftType.(type) {
	case *TypeType:
		if leftType == MetaType {
			_, ok := rawRightType.(*TypeType)
			return ok
		}
		return false
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
	case *InterfaceType:
		rightType, ok := rawRightType.(*InterfaceType)
		if !ok {
			return false
		}

		return sameInterface(leftType, rightType)
	case *ImplType:
		rightType, ok := rawRightType.(*ImplType)
		if !ok {
			return false
		}

		return sameImpl(leftType, rightType)

	case *UnionType:
		rightType, ok := rawRightType.(*UnionType)
		if !ok {
			return false
		}

		return Equals(leftType.Left, rightType.Left) && Equals(leftType.Right, rightType.Right)

	case *IntersectionType:
		rightType, ok := rawRightType.(*IntersectionType)
		if !ok {
			return false
		}

		return Equals(leftType.Left, rightType.Left) && Equals(leftType.Right, rightType.Right)

	default:
		return false
	}
}

// !!!ПОРЯДОК ВАЖЕН!!! ЛЕВОЕ - ТО ЧЕМУ ПРИСВАИВАЕТСЯ ПРАВОЕ
func EqualsSoft(rawLeftType, rawRightType Type) bool {
	if rawLeftType == rawRightType {
		return true
	}

	switch leftType := rawLeftType.(type) {
	case *TypeType:
		if leftType == MetaType {
			_, ok := rawRightType.(*TypeType)
			return ok
		}
		return false
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

		return IsAssignable(leftType.ElementsType, rightType.ElementsType)
	case *HashMapType:
		rightType, ok := rawRightType.(*HashMapType)
		if !ok {
			return false
		}

		if (rightType.KeyType == nil && rightType.ElementType == nil) ||
			(leftType.KeyType == nil && leftType.ElementType == nil) {
			return true
		}

		return Equals(leftType.KeyType, rightType.KeyType) && IsAssignable(leftType.ElementType, rightType.ElementType)
	case *FuncType:
		rightType, ok := rawRightType.(*FuncType)
		if !ok {
			return false
		}

		if len(leftType.Params) != len(rawRightType.(*FuncType).Params) {
			return false
		}

		for i := range leftType.Params {
			// подставляемая функция должна иметь такой же или более широкий тип
			if !IsAssignable(rightType.Params[i].Type, leftType.Params[i].Type) {
				return false
			}
		}

		return IsAssignable(leftType.ReturnType, rightType.ReturnType)
	case *InterfaceType:
		rightType, ok := rawRightType.(*InterfaceType)
		if !ok {
			return false
		}

		return sameInterface(leftType, rightType)
	case *ImplType:
		rightType, ok := rawRightType.(*ImplType)
		if !ok {
			return false
		}

		return sameImpl(leftType, rightType)

	case *UnionType:
		rightType, ok := rawRightType.(*UnionType)
		if !ok {
			return false
		}

		return Equals(leftType.Left, rightType.Left) && Equals(leftType.Right, rightType.Right)

	case *IntersectionType:
		rightType, ok := rawRightType.(*IntersectionType)
		if !ok {
			return false
		}

		return Equals(leftType.Left, rightType.Left) && Equals(leftType.Right, rightType.Right)

	default:
		return false
	}
}

func sameInterface(left, right *InterfaceType) bool {
	if left == right {
		return true
	}
	if left.Decl == nil || right.Decl == nil || left.Decl != right.Decl {
		return false
	}
	return sameFields(left.Fields, right.Fields)
}

func sameImpl(left, right *ImplType) bool {
	if left == right {
		return true
	}
	if left.Decl == nil || right.Decl == nil || left.Decl != right.Decl {
		return false
	}
	if (left.ImplementedType == nil) != (right.ImplementedType == nil) {
		return false
	}
	if left.ImplementedType != nil && !Equals(left.ImplementedType, right.ImplementedType) {
		return false
	}
	return sameFields(left.Fields, right.Fields)
}

func sameFields(left, right map[string]Type) bool {
	if len(left) != len(right) {
		return false
	}
	for name, lt := range left {
		rt, ok := right[name]
		if !ok || !Equals(lt, rt) {
			return false
		}
	}
	return true
}
