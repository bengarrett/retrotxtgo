package ansi

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/sgr"
)

const (
	CSI = "\u001b[" // CSI control sequence introducer.

	reset = -1
)

// Column contains attributes for the idiomatic HTML elements.
type Column struct {
	Class   string // A space-separated list of the CSS classes.
	Style   string // CSS styling declarations overriding any CSS classes.
	Content string // Text content.
}

// Row contains collections of columns that will be used with content division HTML elements.
type Row [][]Column

// HTML writes to w the HTML equivalent of the ANSI encoded plain text data b.
func HTML(w io.Writer, b []byte) error {
	reader := bytes.NewReader(b)
	return HTMLReader(w, reader)
}

// HTMLString returns the HTML equivalent of the ANSI encoded plain text data s.
func HTMLString(s string) (string, error) {
	reader := strings.NewReader(s)
	writer := new(strings.Builder)
	if err := HTMLReader(writer, reader); err != nil {
		return "", err
	}
	return writer.String(), nil
}

// HTMLReader writes to w the HTML equivalent of the ANSI encoded plain text data reader.
func HTMLReader(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	rows := make(Row, 0)
	line, row := reset, reset
	t := sgr.NewSGR()
	for scanner.Scan() {
		line++
		row = reset
		reader := bytes.NewReader(scanner.Bytes())
		scanner := bufio.NewScanner(reader)
		scanner.Split(ScanANSI)
		rows = append(rows, []Column{})
		for scanner.Scan() {
			row++
			if len(scanner.Bytes()) == 0 {
				// TODO: if string is empty but sgr changes, override instead of append?
				// t.DataStream(nil)
				t.Bytes = nil
				continue
			}
			t.DataStream(scanner.Bytes())
			rows[line] = append(rows[line], Column{
				Class:   t.String(),
				Content: string(t.Bytes),
			})
		}
	}
	temp := `
{{- range $index, $element := .}}
{{- if not $element}}<p id="row-{{printf "%d" $index}}"></p>
{{else}}<div id="row-{{printf "%d" $index}}">
{{range .}}
	{{- if not .Content}}{{else}}  <i class="{{.Class}}">{{.Content}}</i>
{{end}}{{end}}</div>
{{end}}{{end}}`
	tmpl, err := template.New("d").Parse(temp)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, rows)
}

// ScanANSI is a split function for a Scanner that returns each ANSI control sequence introducer as a token.
func ScanANSI(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, CSI); i >= 0 {
		return i + 1, stripCR(data[0:i]), nil
	}
	if atEOF {
		return len(data), stripCR(data), nil
	}
	return 0, nil, nil
}

// stripCR strips the carriage return value from the end of the data.
func stripCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
