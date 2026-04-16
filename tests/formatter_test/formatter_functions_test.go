package formatter_test

import "testing"

func TestFormat_FunctionDeclarations(t *testing.T) {
	input := `
let emptyFunc: fn() -> int = fn  (  ) { return 1; };
let f: fn(int, int)->int = fn(a,b) {
return -a +b;
};
`
	expected := `let emptyFunc: fn() -> int = fn() {
	return 1;
};
let f: fn(int, int) -> int = fn(a, b) {
	return -a + b;
};
`
	assertFormat(t, "Function Declarations", input, expected)
}

func TestFormat_FunctionCalls(t *testing.T) {
	input := `
let c: int = f(
arr[0],
  !true
);
let chained:int = getFunc()();
`
	expected := `let c: int = f(
	arr[0],
	!true
);
let chained: int = getFunc()();
`
	assertFormat(t, "Function Calls", input, expected)
}
