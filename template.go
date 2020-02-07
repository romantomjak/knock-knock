package main

import (
	"fmt"
	"io/ioutil"
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
