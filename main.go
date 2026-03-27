// main.go v0.1.0
package main

import (
	"flag"
	"fmt"
	"funlang/evaluator"
	"funlang/lexer"
	"funlang/object"
	"funlang/parser"
	"funlang/repl"
	"funlang/type_checker"
	"funlang/types"
	"io"
	"os"
	"os/user"
)

func InterpretProgram(program string, out io.Writer) {
	typeEnv := types.NewTypeEviroment()
	env := object.NewEnvironment()

	lxr := lexer.New(program)
	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	chk := type_checker.New(typeEnv)

	for _, err := range prs.Errors() {
		fmt.Println(err)
	}
	if len(prs.Errors()) != 0 {
		return
	}

	chk.CheckProgram(prg)
	for _, err := range chk.Errors() {
		fmt.Println(err)
	}
	if len(chk.Errors()) != 0 {
		return
	}

	evaluated := evaluator.Eval(prg, env)
	if evaluated != nil {
		// io.WriteString(out, "Result of file eval is: ")
		// io.WriteString(out, evaluated.Inspect())
		// io.WriteString(out, "\n")
	} else {
		io.WriteString(out, "eval error\n")
	}
}

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	runCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fun run <file name>\n")
		runCmd.PrintDefaults()
	}

	flag.Parse()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 1 {
		fmt.Printf("Hello %s! This is the funlang programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		switch os.Args[1] {
		case "run":
			runCmd.Parse(os.Args[2:])
			file, err := os.ReadFile(runCmd.Arg(0))
			if err != nil {
				panic(err)
			}
			InterpretProgram(string(file), os.Stdout)
		}
	}
}
