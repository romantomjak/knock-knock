package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stdout, os.Args[1:]))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	tmpl, err := NewTemplate("/Users/romantomjak/.knock-knock.toml")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	consul, err := NewConsulClient()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	vault, err := NewVaultClient()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	err = tmpl.Execute(consul, vault)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintln(stdout, tmpl.Contents())
	return 0
}
