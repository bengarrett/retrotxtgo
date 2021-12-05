package color

import (
	"bytes"
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// ColorCSS returns the element colored using CSS syntax highlights.
func ColorCSS(elm string) string {
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "css", style, true)
}

// ColorHTML returns the element colored using HTML syntax highlights.
func ColorHTML(elm string) string {
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "html", style, true)
}

// ColorElm applies color syntax to an element.
func ColorElm(elm, lexer, style string, color bool) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style, color); err != nil {
		logs.FatalMark(fmt.Sprint("html ", lexer), logs.ErrHighlight, err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
}
