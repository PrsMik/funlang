package repl

import (
	"bufio"
	"fmt"
	"funlang/evaluator"
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"funlang/type_checker"
	"funlang/types"
	"io"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnviroment()

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

		typeEnv := types.NewTypeEviroment()
		chk := type_checker.New(typeEnv)
		chk.CheckProgram(prg)
		if len(chk.Errors()) != 0 {
			printCheckerErrors(out, chk.Errors())
			continue
		}

		evaluated := evaluator.Eval(prg, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		} else {
			io.WriteString(out, "eval error\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printCheckerErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Error()+"\n")
	}
}
