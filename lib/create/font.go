package create

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"retrotxt.com/retrotxt/internal/pack"
)

type Font int

const (
	// Automatic uses AutoFont to suggest a font.
	Automatic Font = iota
	// Mona is a Japanese language font for ShiftJIS encoding.
	Mona
	// VGA is an all-purpose 8 pixel IBM/MS-DOS era VGA font.
	VGA
)

var (
	ErrName = errors.New("font name is not known")
	ErrPack = errors.New("font pack is not found")
)

func (f Font) String() string {
	fonts := [...]string{"automatic", "mona", "vga"}
	if f < Automatic || f > VGA {
		return ""
	}
	return fonts[f]
}

// File is the packed filename of the font.
func (f Font) File() string {
	files := [...]string{"ibm-vga8", "mona", "ibm-vga8"}
	if f < Automatic || f > VGA {
		return ""
	}
	return fmt.Sprintf("%s.woff2", files[f])
}

// AutoFont applies the automatic font-family setting to suggest a font based on the given encoding.
func AutoFont(e encoding.Encoding) Font {
	switch e {
	case japanese.ShiftJIS:
		return Mona
	}
	return VGA
}

// Family returns the named font.
func Family(name string) Font {
	switch name {
	case "automatic", "a":
		return Automatic
	case "mona", "m":
		return Mona
	case "vga", "v":
		return VGA
	}
	return -1
}

// Fonts are values for the CSS font-family attribute.
func Fonts() []string {
	return []string{Automatic.String(), Mona.String(), VGA.String()}
}

// FontCSS creates the CSS required for customized fonts.
func FontCSS(name string, embed bool) (b []byte, err error) {
	f := Family(name).String()
	if f == "" {
		return nil, fmt.Errorf("font css %q: %w", name, ErrName)
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
	data := struct {
		Name string
		URL  template.HTML // use HTML type to avoid contextual encoding
	}{
		Name: f,
	}
	if embed {
		url, err := fontBase64(f)
		if err != nil {
			return nil, fmt.Errorf("binary font to base64 failed: %w", err)
		}
		data.URL = template.HTML(url)
	} else {
		data.URL = template.HTML(Family(name).File())
	}
	var out bytes.Buffer
	t, err := template.New("fontface").Parse(css)
	if err != nil {
		return nil, fmt.Errorf("fontface new template failed: %w", err)
	}
	err = t.Execute(&out, data)
	if err != nil {
		return nil, fmt.Errorf("fontface execute template failed: %w", err)
	}
	return out.Bytes(), nil
}

func fontBase64(name string) (string, error) {
	f := Family(name).File()
	if f == "" {
		return "", fmt.Errorf("font base64 %q: %w", name, ErrName)
	}
	b := pack.Get(fmt.Sprintf("font/%s", f))
	if len(b) == 0 {
		return "", fmt.Errorf("font base64 %q: %w", f, ErrPack)
	}
	var s bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &s)
	defer encoder.Close()
	if _, err := encoder.Write(b); err != nil {
		return "", fmt.Errorf("font base64 failed to write b: %w", err)
	}
	return fmt.Sprintf("data:application/font-woff2;charset=utf-8;base64,%s", &s), encoder.Close()
}
