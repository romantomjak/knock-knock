package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
)

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
func (t *Template) Execute() error {
	tmpl := template.New("")

	tmpl, err := tmpl.Parse(t.contents)
	if err != nil {
		return fmt.Errorf("template: %s", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, nil)
	if err != nil {
		return fmt.Errorf("template: %s", err)
	}

	t.contents = out.String()

	return nil
}

// Contents returns the contents of the template
func (t *Template) Contents() string {
	return t.contents
}
