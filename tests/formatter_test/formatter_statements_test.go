package formatter_test

import "testing"

func TestFormat_LetAndReturn(t *testing.T) {
	input := `
let  x :int =1;
let y  : string= "hello";
return   x;
`
	expected := `let x: int = 1;
let y: string = "hello";
return x;
`
	assertFormat(t, "Basic Let and Return", input, expected)
}

func TestFormat_IfElse(t *testing.T) {
	input := `
let res: int = if(x ==1){
return 10;
}else {
   return 20;} ;

let empty: int = if (true) {} else {};
`
	expected := `let res: int = if (x == 1) {
	return 10;
} else {
	return 20;
};

let empty: int = if (true) {
} else {
};
`
	assertFormat(t, "If Else Blocks", input, expected)
}
