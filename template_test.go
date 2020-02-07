package main

import (
	"io/ioutil"
	"os"
	"testing"
)

type InMemoryClient struct {
	data map[string]interface{}
}

func (c *InMemoryClient) Read(path string) (interface{}, error) {
	return c.data[path], nil
}

// Test that parser returns error when template cannot be found
func TestTemplate_TemplateNotFound(t *testing.T) {
	_, err := NewTemplate("death-star-was-an-inside-job.txt")
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}

// Test that parser returns errors when unknown functions are used in templates
func TestTemplate_UnknownFunc(t *testing.T) {
	in := createTempfile([]byte(`host = {{ unknownfunc "service/myservice/db/host" }}`), t)
	defer deleteTempfile(in, t)

	tmpl, _ := NewTemplate(in.Name())
	err := tmpl.Execute(nil)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}

// Test that template rendering resolves dependencies
func TestTemplate_Render(t *testing.T) {
	in := createTempfile([]byte(`host = {{ key "service/myservice/db/host" }}`), t)
	defer deleteTempfile(in, t)

	tmpl, _ := NewTemplate(in.Name())

	cd := make(map[string]interface{})
	cd["service/myservice/db/host"] = "my-host"
	consul := &InMemoryClient{cd}

	tmpl.Execute(consul)

	out := `host = my-host`
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
