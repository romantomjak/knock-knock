package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
)

func clientFunc(c Client) func(string) (interface{}, error) {
	return func(path string) (interface{}, error) {
		return c.Read(path)
	}
}

// Template is the internal representation of a configuration file
type Template struct {
	filename string
	contents string
}

// NewTemplate creates and reads a new configuration template
func NewTemplate(filename string) (*Template, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("template: %s", err)
	}
	return &Template{filename, string(contents)}, nil
}

// Execute evaluates the template
func (t *Template) Execute(consul Client, vault Client) error {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"key":    clientFunc(consul),
		"secret": clientFunc(vault),
	})

	tmpl, err := tmpl.Parse(t.contents)
	if err != nil {
		return fmt.Errorf("parse: %s", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, nil)
	if err != nil {
		return fmt.Errorf("execute: %s", err)
	}

	t.contents = out.String()

	return nil
}

// Contents returns the contents of the template
func (t *Template) Contents() string {
	return t.contents
}
