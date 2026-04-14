package formatter

import (
	"bytes"
	"funlang/lexer"
	"funlang/parser"
	"strings"
	"testing"
)

func TestFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Basic let statements and comment alignment",
			input: `
let  x :int =1; // short
let veryLongVar :  string=  "hello"; // long comment


let y:bool= true;
let z:bool= true;
`,
			// Ожидаем:
			// 1. Схлопывание 2 пустых строк в 1
			// 2. Выравнивание комментариев
			// 3. Расстановку пробелов после двоеточий и вокруг '='
			expected: `let x: int = 1;                    // short
let veryLongVar: string = "hello"; // long comment

let y: bool = true;
let z: bool = true;
`,
		},
		{
			name: "Array and HashMap multi-line layout",
			input: `
let arr :[int]= [ 1,

2,3
];
let mp: {string: int} = {"a" :1,  "b":  2};
`,
			// Ожидаем:
			// 1. Сохранение пустой строки внутри массива
			// 2. Отступы для перенесенных элементов
			// 3. Отсутствие пробела перед первой и после последней скобки
			expected: `let arr: [int] = [
	1,

	2, 3
];
let mp: {string: int} = {"a": 1, "b": 2};
`,
		},
		{
			name: "Infix expressions and continuation indent",
			input: `
let math : int = 1+2 * 3-4/5;
let multiLineMath: int = 100 +
200 +
      300;
`,
			// Ожидаем:
			// 1. Пробелы вокруг инфиксных операторов
			// 2. Дополнительный таб (continuation indent) для перенесенных частей
			expected: `let math: int = 1 + 2 * 3 - 4 / 5;
let multiLineMath: int = 100 +
	200 +
	300;
`,
		},
		{
			name: "If expression and Block statements",
			input: `
let res: int = if(x ==1){
return 10;
}else {
   return 20;} ;
`,
			// Ожидаем:
			// 1. Отступы внутри блоков {}
			// 2. Корректные переносы скобок
			expected: `let res: int = if (x == 1) {
	return 10;
} else {
	return 20;
};
`,
		},
		{
			name: "Functions, prefix, index and calls",
			input: `
let f: fn(int, int)->int = fn(a: int,b: int) {
return -a +b;
};
let arr: [int] = [1, 2];
let c: int = f(
arr[0],
  !true
);
`,
			// Ожидаем:
			// 1. Отсутствие пробелов после префиксного оператора (-a, !true)
			// 2. Отсутствие пробелов в index expression (arr[0])
			// 3. Корректные аргументы вызова
			expected: `let f: fn(int, int) -> int = fn(a: int, b: int) {
	return -a + b;
};
let arr: [int] = [1, 2];
let c: int = f(
	arr[0],
	!true
);
`,
		},
		{
			name: "Standalone comments and trailing empty lines",
			input: `
// standalone 1



// standalone 2
let z: int = 0;




`,
			// Ожидаем: схлопывание лишних пустых строк в конце файла и между комментариями
			expected: `// standalone 1

// standalone 2
let z: int = 0;
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			lxr := lexer.New(tt.input)
			prs := parser.New(lxr)
			prog := prs.ParseProgram()
			fmtr := New(out)

			if len(prs.Errors()) != 0 {
				t.Fatalf("Parser encountered errors: %v", prs.Errors())
			}

			got := fmtr.FormatProgram(prog)

			expectedClean := strings.TrimLeft(tt.expected, "\n")
			gotClean := strings.TrimLeft(got, "\n")

			if gotClean != expectedClean {
				t.Errorf("FormatProgram() mismatch.\n=== GOT ===\n%s\n=== EXPECTED ===\n%s", gotClean, expectedClean)
			}
		})
	}
}
