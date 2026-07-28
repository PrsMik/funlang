package types

// проверяет можно ли expected (левому) присвоить значение с типом actual (правое)
func IsAssignable(expected, actual Type) bool {
	if Equals(expected, actual) {
		return true
	}

	if expectedUnion, ok := expected.(*UnionType); ok {
		return IsAssignable(expectedUnion.Left, actual) || IsAssignable(expectedUnion.Right, actual)
	}

	if actualIntersection, ok := actual.(*IntersectionType); ok {
		if exprectedIntersection, ok := expected.(*IntersectionType); ok {
			return IsAssignable(exprectedIntersection.Left, actualIntersection.Left) &&
				IsAssignable(exprectedIntersection.Right, actualIntersection.Right)
		}
		return IsAssignable(expected, actualIntersection.Left) || IsAssignable(expected, actualIntersection.Right)
	}

	if actualImpl, ok := actual.(*ImplType); ok {
		return IsAssignable(expected, actualImpl.ImplementedType)
	}

	if actuaTypeType, ok := actual.(*TypeType); ok {
		return IsAssignable(expected, actuaTypeType.Underlying)
	}

	return false
}

func Equals(rawLeftType, rawRightType Type) bool {
	if rawLeftType == rawRightType {
		return true
	}

	switch leftType := rawLeftType.(type) {
	case *TypeType:
		// return true
		_, ok := rawRightType.(*TypeType)
		return ok
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
