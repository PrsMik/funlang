package repl

import (
	"bufio"
	"fmt"
	"funlang/lexer"
	"funlang/parser"
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
		prs := parser.New(lxr)
		prg := prs.ParseProgram()

		if len(prs.Errors()) != 0 {
			printParserErrors(out, prs.Errors())
			continue
		}

		io.WriteString(out, prg.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
