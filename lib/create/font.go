package create

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"retrotxt.com/retrotxt/internal/pack"
)

// FontFamily values for the CSS font-family.
var FontFamily = []string{"automatic", "mona", "vga"}

// FontCSS creates the CSS required for customized fonts.
func FontCSS(name string, embed bool) (b []byte, err error) {
	if !font(name) {
		return nil, errors.New("create.fontcss: font name is not known: " + name)
	}
	const css = `@font-face {
  font-family: '{{.Name}}';
  src: url({{.URL}}) format('woff2');
  font-display: swap;
}

.font-{{.Name}} {
  font-family: {{.Name}};
}

body {
  font-family: {{.Name}};
}

main pre {
  font-family: {{.Name}}; /* temp */
}
`
	f := SelectFont(name)
	data := struct {
		Name string
		URL  template.HTML // use HTML type to avoid contextual encoding
	}{
		Name: f,
	}
	if embed {
		url, err := fontBase64(f)
		if err != nil {
			return nil, err
		}
		data.URL = template.HTML(url)
	} else {
		data.URL = template.HTML(f + ".woff2")
	}
	var out bytes.Buffer
	t, err := template.New("fontface").Parse(css)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&out, data)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// SelectFont blah...
func SelectFont(name string) (font string) {
	switch name {
	case "automatic":
		// TODO:
		font = "vga"
	case "mona", "vga":
		font = name
	}
	return strings.ToLower(font)
}

func font(name string) bool {
	if name == "" {
		return false
	}
	for _, f := range FontFamily {
		if f == name {
			return true
		}
	}
	return false
}

func fontBase64(name string) (string, error) {
	if name == "" {
		return "", errors.New("create.fontbase64: empty name")
	}
	n := name
	// TODO: temp work around until font names are confirmed
	switch name {
	case "vga":
		n = "ibm-vga8"
	}
	b := pack.Get(fmt.Sprintf("font/%s.woff2", n))
	if len(b) == 0 {
		return "", errors.New("create.fontbase64: unknown pack name: " + n)
	}
	var s bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &s)
	encoder.Write(b)
	encoder.Close()
	//s := base64.StdEncoding.EncodeToString(b)
	return fmt.Sprintf("data:application/font-woff2;charset=utf-8;base64,%s", &s), nil
}
