package formatter_test

import "testing"

func TestFormat_Arrays(t *testing.T) {
	input := `
let empty: []int = []int{  };
let arr :[]int= []int{ 1,

2,3
};
`
	expected := `let empty: []int = []int{};
let arr: []int = []int{
	1,

	2, 3
};
`
	assertFormat(t, "Array Formatting", input, expected)
}

func TestFormat_HashMaps(t *testing.T) {
	input := `
let emptyHash: {string}int = {string}int{  };
let mp: {string}int = {string}int{"a" :1,  "b":  2};
`
	expected := `let emptyHash: {string}int = {string}int{};
let mp: {string}int = {string}int{"a": 1, "b": 2};
`
	assertFormat(t, "HashMap Formatting", input, expected)
}

func TestFormat_IndexExpressions(t *testing.T) {
	input := `let val: int = arr [ 0 ] + map [ "a" ] ;`
	expected := `let val: int = arr[0] + map["a"];
`
	assertFormat(t, "Index Expressions", input, expected)
}
