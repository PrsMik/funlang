package object

import (
	"funlang/types"
)

func GetTypeFromObject(rawObj Object) types.Type {
	switch object := rawObj.(type) {
	case *Nil:
		return &types.NilType{}
	case *TypeObject:
		return &types.TypeType{Underlying: object.TypeValue}
	case *Integer:
		return &types.IntType{}
	case *Boolean:
		return &types.BoolType{}
	case *String:
		return &types.StringType{}
	case *Array:
		return object.TypeObj.TypeValue
	case *HashMap:
		return object.TypeObj.TypeValue
	case *Function:
		return object.TypeObj.TypeValue
	case *ImplInstance:
		return object.TypeObj.TypeValue
	}
	return &types.IllegalType{}
}
