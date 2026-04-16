package object_test

import (
	"funlang/object"
	"testing"
)

func TestIntegerHashKey(t *testing.T) {
	int1 := &object.Integer{Value: 5}
	int2 := &object.Integer{Value: 5}
	int3 := &object.Integer{Value: 10}

	if int1.HashKey() != int2.HashKey() {
		t.Errorf("integers with same content must have the same hash key")
	}
	if int1.HashKey() == int3.HashKey() {
		t.Errorf("integers with different content must have different hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	bool1 := &object.Boolean{Value: true}
	bool2 := &object.Boolean{Value: true}
	bool3 := &object.Boolean{Value: false}
	bool4 := &object.Boolean{Value: false}

	if bool1.HashKey() != bool2.HashKey() {
		t.Errorf("booleans with same content must have the same hash key")
	}
	if bool3.HashKey() != bool4.HashKey() {
		t.Errorf("booleans with same content must have the same hash key")
	}
	if bool1.HashKey() == bool3.HashKey() {
		t.Errorf("booleans with different content must have different hash keys")
	}
}

func TestStringHashKey(t *testing.T) {
	str1 := &object.String{Value: "Hello World"}
	str2 := &object.String{Value: "Hello World"}
	str3 := &object.String{Value: "My name is johnny"}
	str4 := &object.String{Value: "My name is johnny"}

	if str1.HashKey() != str2.HashKey() {
		t.Errorf("strings with same content must have the same hash key")
	}
	if str3.HashKey() != str4.HashKey() {
		t.Errorf("strings with same content must have the same hash key")
	}
	if str1.HashKey() == str3.HashKey() {
		t.Errorf("strings with different content must have different hash keys")
	}
}
