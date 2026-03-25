package object

import (
	"bytes"
	"fmt"
	"funlang/types"
	"strings"
)

func (n *Null) Inspect() string { return "null" }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (s *String) Inspect() string { return s.Value }

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (h *HashMap) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

func (b *Builtin) Inspect() string { return "builtin function" }

func (e *Error) Inspect() string { return "RUNTIME ERROR: " + e.Message }

var objectTypes map[ObjectType]string = map[ObjectType]string{
	NULL_OBJ:         (&types.NullType{}).Signature(),
	INTEGER_OBJ:      (&types.IntType{}).Signature(),
	BOOLEAN_OBJ:      (&types.BoolType{}).Signature(),
	STRING_OBJ:       (&types.StringType{}).Signature(),
	ARRAY_OBJ:        (&types.ArrayType{}).Signature(),
	HASH_OBJ:         (&types.HashMapType{}).Signature(),
	RETURN_VALUE_OBJ: "<return val>",
	FUNCTION_OBJ:     "<fn>",
	BUILTIN_OBJ:      "<builtin>",
	ERROR_OBJ:        "<error>",
}

func LookUpObjSignature(objectType ObjectType) string {
	if obj, ok := objectTypes[objectType]; ok {
		return obj
	}
	return "<unknown>"
}
