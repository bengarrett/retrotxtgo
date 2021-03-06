package create

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"strings"

	"github.com/bengarrett/retrotxtgo/static"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

// Font enum.
type Font uint

const (
	// Automatic uses AutoFont to suggest a font.
	Automatic Font = iota
	// Mona is a Japanese language font for ShiftJIS encoding.
	Mona
	// VGA is an all-purpose 8 pixel IBM/MS-DOS era VGA font.
	VGA
)

func (f Font) String() string {
	return [...]string{"automatic", "mona", "vga"}[f]
}

// File is the packed filename of the font.
func (f Font) File() string {
	files := [...]string{"ibm-vga8", "mona", "ibm-vga8"}
	return fmt.Sprintf("%s.woff2", files[f])
}

// AutoFont applies the automatic font-family setting to suggest a font based on the given encoding.
func AutoFont(e encoding.Encoding) Font {
	if e == japanese.ShiftJIS {
		return Mona
	}
	return VGA
}

// Family returns the named font.
func Family(name string) Font {
	switch strings.ToLower(name) {
	case Automatic.String(), "a":
		return Automatic
	case Mona.String(), "m":
		return Mona
	case VGA.String(), "v":
		return VGA
	default:
		return Automatic
	}
}

// Fonts are values for the CSS font-family attribute.
func Fonts() []string {
	return []string{Automatic.String(), Mona.String(), VGA.String()}
}

// FontCSS creates the CSS required for customized fonts.
func FontCSS(name string, e encoding.Encoding, embed bool) (b []byte, err error) {
	f := Family(name)
	if Family(name) == Automatic {
		f = AutoFont(e)
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
		Name: f.String(),
	}
	if embed {
		url := ""
		url, err = fontBase64(f)
		if err != nil {
			return nil, fmt.Errorf("binary font to base64 failed: %w", err)
		}
		data.URL = template.HTML(url) // nolint:gosec
	} else {
		data.URL = template.HTML(f.File()) // nolint:gosec
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

// FontBase64 encodes a font using the Base64 binary-to-text encoding scheme,
// for embedding into HTML or CSS textfiles.
func fontBase64(f Font) (string, error) {
	name := fmt.Sprintf("font/%s", f.File())
	b, err := static.Font.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("font base64 %q: %w", f.File(), ErrFont)
	}
	var s bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &s)
	defer encoder.Close()
	if _, err := encoder.Write(b); err != nil {
		return "", fmt.Errorf("font base64 failed to write b: %w", err)
	}
	return fmt.Sprintf("data:application/font-woff2;charset=utf-8;base64,%s", &s), nil
}
