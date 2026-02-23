// main.go
package main

import (
	"fmt"
	"funlang/repl"
	"os"
	"os/user"
)

func main() {
	args := os.Args
	fmt.Println(args)
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the funlang programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
