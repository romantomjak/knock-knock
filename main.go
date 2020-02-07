package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	usage = `
Usage: knock-knock [options] service

  Reads a template file with KV and secret paths and renders service
  credentials on screen.

Options:

  -c=<path>
	The path to a configuration file on disk. Defaults to ~/.knock-knock.toml
`
)

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stdout, os.Args[1:]))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	var config string

	flags := flag.NewFlagSet("knock-knock", flag.ContinueOnError)
	flags.StringVar(&config, "c", "", "configuration file")
	flags.Usage = func() {
		fmt.Fprintln(stderr, strings.TrimSpace(usage))
	}

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		config = fmt.Sprintf("%s/.knock-knock.toml", home)
	}

	tmpl, err := NewTemplate(config)
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
