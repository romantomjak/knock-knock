package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/posener/complete"
	"gopkg.in/ini.v1"
)

const (
	usage = `
Usage: knock-knock [-help] [-autocomplete-(un)install] [options] service

  Reads a template file with KV and secret paths and renders service
  credentials on screen.

Options:

  -c=<path>
	The path to a configuration file on disk. Defaults to ~/.knock-knock.conf
`
)

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stdout, os.Args[1:]))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	var filename string

	flags := flag.NewFlagSet("knock-knock", flag.ContinueOnError)
	flags.StringVar(&filename, "c", "", "configuration file")
	flags.Usage = func() {
		fmt.Fprintln(stderr, strings.TrimSpace(usage))
	}

	// shell autocompletion
	cmp := complete.New(
		"knock-knock",
		complete.Command{
			Args: complete.PredictFunc(sectionNames),
			GlobalFlags: complete.Flags{
				"-help":                   complete.PredictNothing,
				"-autocomplete-install":   complete.PredictNothing,
				"-autocomplete-uninstall": complete.PredictNothing,
				"-c":                      complete.PredictFiles("*.conf"),
			},
		},
	)
	cmp.CLI.InstallName = "autocomplete-install"
	cmp.CLI.UninstallName = "autocomplete-uninstall"
	cmp.AddFlags(flags)

	err := flags.Parse(args)
	if err != nil {
		return 1
	}

	// in case that the completion was invoked and ran as a completion script
	// or handled a flag, the Complete method will return true, in which case,
	// the program has nothing to do and should return
	if cmp.Complete() {
		return 0
	}

	arguments := flags.Args()
	if len(arguments) != 1 {
		fmt.Fprintln(stderr, strings.TrimSpace(usage))
		return 1
	}

	if filename == "" {
		filename, err = defaultFilename()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	tmpl, err := NewTemplate(filename)
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

	config, err := ini.Load([]byte(tmpl.Contents()))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	section, err := config.GetSection(arguments[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	for _, key := range section.Keys() {
		fmt.Fprintln(stdout, fmt.Sprintf("%s = %s", key.Name(), key.Value()))
	}

	return 0
}

// defaultFilename returns the full path to the default configuration file
func defaultFilename() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("configuration file was not provided and %s", err)
	}
	return fmt.Sprintf("%s/.knock-knock.conf", home), nil
}

// sectionNames reads sections names from the default configuration file.
// It is used for providing suggestions for shell autocompletion.
func sectionNames(args complete.Args) []string {
	filename, err := defaultFilename()
	if err != nil {
		return nil
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	config, err := ini.Load([]byte(contents))
	if err != nil {
		return nil
	}

	// all sections minus the DEFAULT section
	return config.SectionStrings()[1:]
}
