package example

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

var Config = fmt.Sprintf("  %s %s %s\n%s %s %s",
	meta.Bin, "config setup", "# Walk through all the settings",
	meta.Bin, "config set --list", "# List all the settings in use")

var Create = fmt.Sprintf("  %s%s\n%s%s\n%s%s\n%s%s\n%s%s",
	"# print a HTML file created from file.txt\n",
	fmt.Sprintf("%s create file.txt --title \"A text file\" --description \"Some text goes here\"", meta.Bin),
	"# save HTML files created from file1.txt and file2.asc\n",
	fmt.Sprintf("%s create file1.txt file2.asc --save", meta.Bin),
	"# save and compress a HTML file created from file.txt in Downloads.\n",
	fmt.Sprintf("%s create ~{{.}}Downloads{{.}}file.txt --compress", meta.Bin),
	"# host the HTML file created from file.txt\n",
	fmt.Sprintf("%s create file.txt --serve=%d", meta.Bin, meta.WebPort),
	"# pipe a HTML file created from file.txt\n",
	fmt.Sprintf("%s create file.txt | %s", meta.Bin, catCmd()))

// catCmd returns the os command name to concatenate a file to standard output.
func catCmd() string {
	s := "cat"
	if runtime.GOOS == "windows" {
		s = "type"
	}
	return s
}

var Info = fmt.Sprintf("  %s\n%s",
	fmt.Sprintf("%s config info   # List the default setting values", meta.Bin),
	fmt.Sprintf("%s config set -c # List the settings and help hints", meta.Bin))

var Set = fmt.Sprintf("  %s %s %s\n%s %s %s\n%s %s %s",
	meta.Bin, "config set --list", "# List the available settings",
	meta.Bin, "config set html.meta.description", "# Edit the meta description setting",
	meta.Bin, "config set style.info style.html", fmt.Sprintf("# Edit both the %s color styles", meta.Name),
)

// exampleCmd returns help usage examples.
func Print(tmpl string) string {
	if tmpl == "" {
		return ""
	}
	var b bytes.Buffer
	// change example operating system path separator
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(&b, string(os.PathSeparator))
	if err != nil {
		log.Fatal(err)
	}
	// color the example text except text following
	// the last hash #, which is treated as a comment
	const cmmt, sentence = "#", 2
	scanner, s := bufio.NewScanner(&b), ""
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), cmmt)
		l := len(ss)
		if l < sentence {
			s += str.ColInf(scanner.Text()) + "\n  "
			continue
		}
		// do not the last hash as a comment
		ex := strings.Join(ss[:l-1], cmmt)
		s += str.ColInf(ex)
		s += fmt.Sprintf("%s%s\n  ", color.Secondary.Sprint(cmmt), ss[l-1])
	}
	return strings.TrimSpace(s)
}
