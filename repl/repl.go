package repl

import (
	"bufio"
	"fmt"
	"funlang/lexer"
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
		lexer := lexer.New(line)
		for curToken := lexer.NextToken(); curToken.Type != token.EOF; curToken = lexer.NextToken() {
			tokenStr, _ := token.LookupString(curToken.Type)
			fmt.Printf("Token: %s; Literal: %s\n", tokenStr, curToken.Literal)
		}
	}
}
