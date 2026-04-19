package types_test

import (
	"funlang/types"
	"testing"
)

func TestNewInfo(t *testing.T) {
	info := types.NewInfo()

	if info == nil {
		t.Fatal("NewInfo returned nil")
	}

	if info.TypesInfo == nil {
		t.Error("TypesInfo map is nil")
	}
	if info.TypeNodes == nil {
		t.Error("TypeNodes map is nil")
	}
	if info.Definitions == nil {
		t.Error("Definitions map is nil")
	}
	if info.Scopes == nil {
		t.Error("Scopes map is nil")
	}
	if info.ExpectedTypes == nil {
		t.Error("ExpectedTypes map is nil")
	}
}
