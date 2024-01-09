package main

import (
	"fmt"
	"monkey-int/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("== MONKEY INTERPRETER ==\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
