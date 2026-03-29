package repl

import (
	"bufio"
	"fmt"
	"funlang/evaluator"
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"funlang/token"
	"funlang/type_checker"
	"funlang/types"
	"io"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	typeEnv := types.NewTypeEviroment()
	env := object.NewEnvironment()

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
			PrintParserErrors(out, prs.Errors())
			continue
		}

		chk := type_checker.New(typeEnv)
		chk.CheckProgram(prg)
		if len(chk.Errors()) != 0 {
			PrintCheckerErrors(out, chk.Errors())
			continue
		}

		evaluated := evaluator.Eval(prg, env)
		if evaluated != nil {
			// io.WriteString(out, evaluated.Inspect())
			// io.WriteString(out, "\n")
		} else {
			io.WriteString(out, "eval error\n")
		}
	}
}

func PrintParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, err := range errors {
		tknType, _ := token.LookupString(err.Token.Type)
		io.WriteString(out, fmt.Sprintf("\t%s\n\ton token %s: from %s to %s\n",
			err.Msg, tknType, err.Token.Start.String(), err.Token.End.String()))
	}
}

func PrintCheckerErrors(out io.Writer, errors []type_checker.TypeError) {
	for _, err := range errors {
		start := err.Node.Start()
		end := err.Node.End()
		io.WriteString(out, fmt.Sprintf("\t%s\n\ton node %s: from %s to %s\n",
			err.Msg, err.Node.TokenLiteral(), start.String(), end.String()))
	}
}
