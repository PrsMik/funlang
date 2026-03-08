// main.go
package main

import (
	"flag"
	"fmt"
	"funlang/lexer"
	"funlang/parser"
	"funlang/repl"
	"funlang/types"
	"io"
	"os"
	"os/user"
)

func InterpretProgram(program string, out io.Writer) {
	lxr := lexer.New(program)
	prs := parser.New(lxr)
	prg := prs.ParseProgram()
	chk := types.New(nil)

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

	io.WriteString(out, prg.String())
}

func main() {
	relativeFilePath := flag.String("file_rel", "", "relative path to file to be interpreted")
	flag.Parse()
	fmt.Println("file:", *relativeFilePath)
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the funlang programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	if *relativeFilePath == "" {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		file, err := os.ReadFile(*relativeFilePath)
		if err != nil {
			panic(err)
		}
		InterpretProgram(string(file), os.Stdout)
	}
}
