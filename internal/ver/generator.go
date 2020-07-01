//+build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const (
	name   = "value.go"
	dir    = "../../lib/version"
	deploy = "../.version"
)

type data struct {
	Version string
	Date    string
}

var (
	t = time.Now()
	d = []data{
		{version(), t.Format(time.RFC3339)},
	}
	tmpl = template.Must(template.New("versionTest").Parse(`package version

// Code generated by go generate; ignore.

// B holds the build and version information.
var B = Build{
	Commit:  "n/a",
	Date:    "{{.Date}}",
	Domain:  "retrotxt.com",
	Version: "{{.Version}}",
}`))
)

func version() string {
	b, err := ioutil.ReadFile(deploy)
	if err != nil {
		log.Fatal("error reading:", deploy, err)
	}
	b = bytes.ReplaceAll(b, []byte("\n"), nil)
	return string(b)
}

func main() {
	fmt.Println("version syntax:", string(version()))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("directory does not exist:", dir)
	}
	// create value.go file
	p, err := filepath.Abs(filepath.Join(dir, name))
	f, err := os.Create(p)
	if err != nil {
		log.Fatal("error creating:", p)
	}
	defer f.Close()
	buf := &bytes.Buffer{}
	// execute template
	if err = tmpl.Execute(buf, d[0]); err != nil {
		log.Fatal("error executing template:", err)
	}
	// format the generated code
	code, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(string(buf.Bytes()))
		log.Fatal("error formatting code:", err, code)
	}
	if err = ioutil.WriteFile(p, code, os.ModePerm); err != nil {
		log.Fatal("error writing file:", p, err)
	}
}
