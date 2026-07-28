package types

// !!!ПОРЯДОК ВАЖЕН!!!
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
		if leftType == TrueTypeType {
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

		for leftFieldName, leftFieldType := range leftType.Fields {
			rightFieldType, exists := rightType.Fields[leftFieldName]
			if !exists {
				return false
			}
			if !IsAssignable(leftFieldType, rightFieldType) {
				return false
			}
		}

		return true
	default:
		return false
	}
}
