package object

import "funlang/types"

func IsAssignable(rawLeftObject, rawRightObject Object) bool {
	switch leftObject := rawLeftObject.(type) {
	case *TypeObject:
		if rightObject, ok := rawRightObject.(*TypeObject); ok {
			return types.IsAssignable(leftObject.TypeValue, rightObject.TypeValue)
		}
		return false
	default:
		return false
	}
}

// !!!ПОРЯДОК ВАЖЕН!!! ЛЕВОЕ - ТО ЧЕМУ ПРИСВАИВАЕТСЯ ПРАВОЕ
func Equals(rawLeftObject, rawRightObject Object) bool {
	switch leftObject := rawLeftObject.(type) {
	case *TypeObject:
		if rightObject, ok := rawRightObject.(*TypeObject); ok {
			return types.IsAssignable(leftObject.TypeValue, rightObject.TypeValue)
		}
		return false
	case *Integer, *Boolean, *String:
		return compareSimpleObjects(rawLeftObject, rawRightObject)
	case *Array:
		if rightObject, ok := rawRightObject.(*Array); ok {
			return compareArrays(leftObject, rightObject)
		}
		return false
	case *HashMap:
		if rightObject, ok := rawRightObject.(*HashMap); ok {
			return compareHashMaps(leftObject, rightObject)
		}
		return false
	default:
		return false
	}
}

func compareSimpleObjects(rawLeftObject, rawRightObject Object) bool {
	switch leftObject := rawLeftObject.(type) {
	case *Integer:
		if rightObject, ok := rawRightObject.(*Integer); ok {
			return leftObject.Value == rightObject.Value
		}
		return false
	case *Boolean:
		if rightObject, ok := rawRightObject.(*Boolean); ok {
			return leftObject.Value == rightObject.Value
		}
		return false
	case *String:
		if rightObject, ok := rawRightObject.(*String); ok {
			return leftObject.Value == rightObject.Value
		}
		return false
	default:
		return false
	}
}

func compareArrays(leftArray, rightArray *Array) bool {
	if len(leftArray.Elements) != len(rightArray.Elements) {
		return false
	}
	for i := range leftArray.Elements {
		if !Equals(leftArray.Elements[i], rightArray.Elements[i]) {
			return false
		}
	}
	return true
}

func compareHashMaps(leftMap, rightMap *HashMap) bool {
	if len(leftMap.Pairs) != len(rightMap.Pairs) {
		return false
	}

	for k, v1 := range leftMap.Pairs {
		if v2, ok := rightMap.Pairs[k]; !ok || !Equals(v1.Value, v2.Value) {
			return false
		}
	}

	return true
}
