package object_test

import (
	"funlang/object"
	"testing"
)

func TestEnvironment_GetSet(t *testing.T) {
	env := object.NewEnvironment()
	obj := &object.Integer{Value: 42}

	env.Set("myVar", obj)

	val, ok := env.Get("myVar")
	if !ok {
		t.Fatalf("expected to find 'myVar' in environment")
	}
	if val != obj {
		t.Errorf("expected value %v, got %v", obj, val)
	}

	_, ok = env.Get("unknown")
	if ok {
		t.Errorf("expected not to find 'unknown' in environment")
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	globalEnv := object.NewEnvironment()
	globalObj := &object.Integer{Value: 1}
	globalEnv.Set("x", globalObj)

	localEnv := object.NewEnclosedEnvironment(globalEnv)
	localObj := &object.Integer{Value: 2}
	localEnv.Set("y", localObj)

	val, ok := localEnv.Get("y")
	if !ok || val != localObj {
		t.Errorf("failed to get local variable: expected %v, got %v", localObj, val)
	}

	val, ok = localEnv.Get("x")
	if !ok || val != globalObj {
		t.Errorf("failed to get outer variable: expected %v, got %v", globalObj, val)
	}

	shadowObj := &object.Integer{Value: 99}
	localEnv.Set("x", shadowObj)

	val, ok = localEnv.Get("x")
	if !ok || val != shadowObj {
		t.Errorf("failed shadowing expected %v, got %v", shadowObj, val)
	}

	val, ok = globalEnv.Get("x")
	if !ok || val != globalObj {
		t.Errorf("outer variable was improperly modified: expected %v, got %v", globalObj, val)
	}
}
