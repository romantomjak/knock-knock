package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// Test that parser returns error when template cannot be found
func TestTemplate_TemplateNotFound(t *testing.T) {
	_, err := NewTemplate("death-star-was-an-inside-job.txt")
	if err == nil {
		t.Fatal("expected an error")
	}
}

// Test that parser can read and return parsed template
func TestTemplate_Render(t *testing.T) {
	in := createTempfile([]byte("hello world"), t)
	defer deleteTempfile(in, t)

	tmpl, _ := NewTemplate(in.Name())

	err := tmpl.Execute()
	if err != nil {
		t.Fatal("unexpected error")
	}

	out := "hello world"
	if tmpl.Contents() != out {
		t.Fatalf("expected %q to match %q", tmpl.Contents(), out)
	}
}

func createTempfile(b []byte, t *testing.T) *os.File {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Errorf("cannot create tempfile: %s", err)
	}

	if len(b) > 0 {
		_, err = f.Write(b)
		if err != nil {
			t.Errorf("cannot write to tempfile: %s", err)
		}
	}

	return f
}

func deleteTempfile(f *os.File, t *testing.T) {
	if err := os.Remove(f.Name()); err != nil {
		t.Errorf("cannot delete tempfile: %s", err)
	}
}
