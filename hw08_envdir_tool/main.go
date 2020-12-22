package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Not enough arguments`)
		os.Exit(-1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("Error readind dir \"%s\": %s\n", os.Args[1], err)
		os.Exit(-1)
	}

	retCode := RunCmd(os.Args[2:], env)

	os.Exit(retCode)
}
