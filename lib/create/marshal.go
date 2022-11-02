package create

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
)

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	ExternalEmbed    bool
	FontEmbed        bool
	MetaGenerator    bool
	MetaNoTranslate  bool
	MetaRetroTxt     bool
	BuildVersion     string
	CacheRefresh     string
	Comment          string
	FontFamily       string
	MetaAuthor       string
	MetaColorScheme  string
	MetaDesc         string
	MetaKeywords     string
	MetaReferrer     string
	MetaRobots       string
	MetaThemeColor   string
	PageTitle        string
	PreText          string
	SauceTitle       string
	SauceAuthor      string
	SauceGroup       string
	SauceDescription string
	SauceWidth       uint
	SauceLines       uint
	CSSEmbed         template.CSS
	HTMLEmbed        template.HTML
	ScriptEmbed      template.JS
}

// Comment content for the meta retrotxt attribute.
func (args *Args) Comment(lb filesystem.LB, r ...rune) string {
	const na = "n/a"
	var (
		l int
		w int
	)
	name, e, b := na, na, []byte(string(r))
	if args.Source.Encoding != nil {
		e = fmt.Sprint(args.Source.Encoding)
	}
	lbs := filesystem.LineBreak(lb, false)
	l, err := filesystem.Lines(bytes.NewReader(b), lb)
	if err != nil {
		l = -1
	}
	w, err = filesystem.Columns(bytes.NewReader(b), lb)
	if err != nil {
		w = -1
	}
	if args.Source.Name != "" {
		name = args.Source.Name
	}
	return fmt.Sprintf("encoding: %s; line break: %s; length: %d; width: %d; name: %s",
		e, lbs, l, w, name)
}

// Marshal transforms bytes into UTF-8, creates the page meta and template data.
func (args *Args) Marshal(b *[]byte) (PageData, error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	p := PageData{}
	var out bytes.Buffer
	bt := args.Source.BBSType
	r := bytes.Runes(*b)
	switch {
	case bt == bbs.ANSI:
		// temp placeholder
		p.PreText = string(r)
		if p.MetaRetroTxt {
			lb := filesystem.LineBreaks(true, r...)
			p.Comment = args.Comment(lb, r...)
		}
	case bt < bbs.ANSI:
		// convert bytes into utf8
		p.PreText = string(r)
		if p.MetaRetroTxt {
			lb := filesystem.LineBreaks(true, r...)
			p.Comment = args.Comment(lb, r...)
		}
	case bt > bbs.ANSI:
		if err := bt.HTML(&out, []byte(string(r))); err != nil {
			return PageData{}, err
		}
		p.HTMLEmbed = template.HTML(out.Bytes()) //nolint:gosec
	}
	return p, nil
}
