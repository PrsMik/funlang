package formatter_test

import "testing"

func TestFormat_Comments(t *testing.T) {
	input := `
let  x :int =1; // short
let veryLongVar :  string=  "hello"; // long comment
// standalone 1



// standalone 2
let z: int = 0;
`
	expected := `let x: int = 1;                    // short
let veryLongVar: string = "hello"; // long comment
// standalone 1

// standalone 2
let z: int = 0;
`
	assertFormat(t, "Comments Alignment and Standalone", input, expected)
}

func TestFormat_TrailingEmptyLines(t *testing.T) {
	input := `
let x: int = 1;



`
	expected := `let x: int = 1;
`
	assertFormat(t, "Trailing Empty Lines", input, expected)
}
