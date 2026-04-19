// main.go v0.1.0
package main

import (
	"bytes"
	"flag"
	"fmt"
	"funlang/evaluator"
	"funlang/formatter"
	"funlang/lexer"
	"funlang/lsp"
	"funlang/object"
	"funlang/parser"
	"funlang/repl"
	"funlang/type_checker"
	"funlang/types"
	"io"
	"os"
)

func FormatFile(program string, out io.Writer) (string, error) {
	var newFileText bytes.Buffer
	typeEnv := types.NewTypeEviroment()

	lxr := lexer.New(program)

	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	if len(prs.Errors()) != 0 {
		repl.PrintParserErrors(out, prs.Errors())
		return "", fmt.Errorf("Parser erorrs!")
	}

	chk := type_checker.New(typeEnv, nil)
	chk.CheckProgram(prg)
	if len(chk.Errors()) != 0 {
		repl.PrintCheckerErrors(out, chk.Errors())
		return "", fmt.Errorf("Parser erorrs!")
	}

	fmtr := formatter.New(&newFileText)
	return fmtr.FormatProgram(prg), nil
}

func InterpretProgram(program string, out io.Writer) {
	typeEnv := types.NewTypeEviroment()
	env := object.NewEnvironment()

	lxr := lexer.New(program)

	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	if len(prs.Errors()) != 0 {
		repl.PrintParserErrors(out, prs.Errors())
		return
	}

	chk := type_checker.New(typeEnv, nil)
	chk.CheckProgram(prg)
	if len(chk.Errors()) != 0 {
		repl.PrintCheckerErrors(out, chk.Errors())
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fun <command> [args]\n\n")
		fmt.Fprintf(os.Stderr, "Available commands:\n")
		fmt.Fprintf(os.Stderr, "  run    interpretes file (fun run <file name>)\n")
		fmt.Fprintf(os.Stderr, "  lsp    start a LSP server\n")
		fmt.Fprintf(os.Stderr, "  fmt    formats file (fun fmt <file name>)\n\n")
		fmt.Fprintf(os.Stderr, "Use \"fun <command> --help\" to get more information about each command.\n")
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	lspCmd := flag.NewFlagSet("lsp", flag.ExitOnError)
	fmtCmd := flag.NewFlagSet("fmt", flag.ExitOnError)

	runCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fun run <file name>\n")
		runCmd.PrintDefaults()
	}

	lspMode := lspCmd.String("mode", "tcp", "Mode to run lsp. Can be \"tcp\" or \"stdio\" (tcp by default)")
	lspPort := lspCmd.String("port", "127.0.0.1:5007", "Used only with -mode=tcp flag to set up server port")

	lspCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fun lsp\n This will start a LSP server on 127.0.0.1:5007\n")
		lspCmd.PrintDefaults()
	}

	fmtCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fun fmt <file name>\n")
		fmtCmd.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Printf("Hello! This is the funlang programming language!\n")
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
		case "lsp":
			lspCmd.Parse(os.Args[2:])
			mode := true
			if *lspMode == "stdio" {
				mode = false
			}
			lsp.StartServer(mode, *lspPort)
			return
		case "fmt":
			fmtCmd.Parse(os.Args[2:])
			file, err := os.ReadFile(fmtCmd.Arg(0))
			if err != nil {
				panic(err)
			}

			newText, err := FormatFile(string(file), os.Stdout)
			if err != nil {
				panic(err)
			}

			newfile, err := os.Create(fmtCmd.Arg(0))
			if err != nil {
				panic(err)
			}
			defer newfile.Close()

			_, err = newfile.WriteString(newText)
			if err != nil {
				panic(err)
			}
		}
	}
}
