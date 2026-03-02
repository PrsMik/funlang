package repl

import (
	"bufio"
	"fmt"
	"funlang/lexer"
	"funlang/parser"
	"funlang/token"
	"io"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lxr := lexer.New(line)
		for curToken := lxr.NextToken(); curToken.Type != token.EOF; curToken = lxr.NextToken() {
			tokenStr, _ := token.LookupString(curToken.Type)
			fmt.Printf("Token: %s; Literal: %s\n", tokenStr, curToken.Literal)
		}
		lxr = lexer.New(line)
		prs := parser.New(lxr)
		prg := prs.ParseProgram()
		output := prg.String()
		fmt.Println(output)
	}
}
