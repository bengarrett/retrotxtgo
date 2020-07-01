//+build ignore

package main

// credit: https://dev.to/koddr/the-easiest-way-to-embed-static-files-into-a-binary-file-in-your-golang-app-no-external-dependencies-43pc

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: compress static files as bytes and decompress with Get()

const (
	name = "blob.go"
	dir  = "../../static"
)

var (
	conv = map[string]interface{}{"conv": fmtByteSlice}
	tmpl = template.Must(template.New("").Funcs(conv).Parse(`package version

// Code generated by go generate; ignore.

func init() {
	{{- range $name, $file := . }}
		pack.Add("{{ $name }}", []byte{ {{ conv $file }} })
	{{- end }}
}`))
)

func fmtByteSlice(b []byte) string {
	var builder = strings.Builder{}
	for _, v := range b {
		s := fmt.Sprintf("%d,", int(v))
		if _, err := builder.WriteString(s); err != nil {
			log.Fatal(err)
		}
	}
	return builder.String()
}

func main() {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("Directory does not exist: " + dir)
	}
	configs := make(map[string][]byte)
	// walk dir
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		var rel string
		if rel, err = filepath.Rel(dir, path); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		} else {
			log.Println("packing ", path)
			b, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf(", error reading file: %s", err)
				return err
			}
			configs[rel] = b
		}
		return nil
	})
	if err != nil {
		log.Fatal("error walking path", dir, err)
	}
	// create blob file
	f, err := os.Create(name)
	if err != nil {
		log.Fatal("error creating", name, dir)
	}
	defer f.Close()
	// create buffer
	builder := &bytes.Buffer{}
	// execute template
	if err = tmpl.Execute(builder, configs); err != nil {
		log.Fatal("error executing template", err)
	}
	// format the generated code
	code, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("error formatting code", err)
	}
	if err = ioutil.WriteFile(name, code, os.ModePerm); err != nil {
		log.Fatal("error writing file", name, err)
	}
}
