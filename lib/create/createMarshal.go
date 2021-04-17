package create

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"time"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/version"
	"retrotxt.com/retrotxt/static"
)

// marshal transforms bytes into UTF-8, creates the page meta and template data.
func (args *Args) marshal(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	// templates are found in the dir static/html/*.gohtml
	switch args.layout {
	case Inline:
		if p, err = args.marshalInline(b); err != nil {
			return p, err
		}
		p = args.marshalStandard(&p)
	case Standard:
		p = args.marshalStandard(&p)
	case Compact: // disables all meta tags
		p = args.marshalCompact(&p)
	case None:
		// do nothing
	default:
		return PageData{}, fmt.Errorf("pagedata %s: %w", args.layout, ErrNoLayout)
	}
	// convert bytes into utf8
	r := bytes.Runes(*b)
	p.PreText = string(r)
	if p.MetaRetroTxt {
		lb := filesystem.LineBreaks(true, r...)
		p.Comment = args.comment(lb, r...)
	}
	return p, nil
}

func (args *Args) marshalCompact(p *PageData) PageData {
	p.PageTitle = args.pageTitle()
	p.MetaGenerator = false
	return *p
}

func (args *Args) marshalInline(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	p.ExternalEmbed = true
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	// styles
	s := bytes.TrimSpace(static.Styles)
	// font
	var f []byte
	f, err = FontCSS(args.FontFamily.Value, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return p, fmt.Errorf("pagedata font error: %w", err)
	}
	f = bytes.TrimSpace(f)
	// merge
	c := [][]byte{s, []byte("/* font */"), f}
	*b = bytes.Join(c, []byte("\n\n"))
	// compress & embed
	*b, err = m.Bytes("text/css", *b)
	if err != nil {
		return p, fmt.Errorf("pagedata minify css: %w", err)
	}
	p.CSSEmbed = template.CSS(string(*b))
	jsp := static.Scripts
	jsp, err = m.Bytes("application/javascript", jsp)
	if err != nil {
		return p, fmt.Errorf("pagedata minify javascript: %w", err)
	}
	p.ScriptEmbed = template.JS(string(jsp))
	return p, nil
}

func (args *Args) marshalStandard(p *PageData) PageData {
	p.FontEmbed = args.FontEmbed
	p.FontFamily = args.fontFamily()
	p.MetaAuthor = args.metaAuthor()
	p.MetaColorScheme = args.metaColorScheme()
	p.MetaDesc = args.metaDesc()
	p.MetaGenerator = args.Metadata.Generator
	p.MetaKeywords = args.metaKeywords()
	p.MetaNoTranslate = args.Metadata.NoTranslate
	p.MetaReferrer = args.metaReferrer()
	p.MetaRobots = args.metaRobots()
	p.MetaRetroTxt = args.Metadata.RetroTxt
	p.MetaThemeColor = args.metaThemeColor()
	p.PageTitle = args.pageTitle()
	// sauce data
	p.SauceTitle = args.SauceData.Title
	p.SauceAuthor = args.SauceData.Author
	p.SauceGroup = args.SauceData.Group
	p.SauceDescription = args.SauceData.Description
	p.SauceWidth = args.SauceData.Width
	p.SauceLines = args.SauceData.Lines
	// generate data
	t := time.Now().UTC()
	p.BuildDate = t.Format(time.RFC3339)
	p.BuildVersion = version.B.Version
	return *p
}

func (args *Args) marshalTextTransform(b *[]byte) (buf bytes.Buffer, err error) {
	tmpl, err := args.newTemplate()
	if err != nil {
		return buf, fmt.Errorf("stdout new template failure: %w", err)
	}
	d, err := args.marshal(b)
	if err != nil {
		return buf, fmt.Errorf("stdout meta and pagedata failure: %w", err)
	}
	if err = tmpl.Execute(&buf, d); err != nil {
		return buf, fmt.Errorf("stdout template execute failure: %w", err)
	}
	return buf, nil
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args *Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, fmt.Errorf("using existing template cache: %w", err)
	}
	if !args.Save.Cache {
		if err := args.templateSave(); err != nil {
			return nil, fmt.Errorf("creating a new template: %w", err)
		}
	} else {
		s, err := os.Stat(args.tmpl)
		switch {
		case os.IsNotExist(err):
			if err2 := args.templateSave(); err2 != nil {
				return nil, fmt.Errorf("saving to the template: %w", err2)
			}
		case err != nil:
			return nil, err
		case s.IsDir():
			return nil, fmt.Errorf("new template %q: %w", args.tmpl, ErrTmplDir)
		}
	}
	// to avoid a potential panic, Stat again in case os.IsNotExist() is true
	s, err := os.Stat(args.tmpl)
	if err != nil {
		return nil, fmt.Errorf("could not access file: %q: %w", args.tmpl, err)
	}
	if err = args.templatePack(); err != nil {
		return nil, fmt.Errorf("template pack: %w", err)
	}
	b, err := args.templateData()
	if s.Size() != int64(len(*b)) {
		if err != nil {
			return nil, fmt.Errorf("template data: %w", err)
		}
		if _, _, err := filesystem.Save(args.tmpl, *b...); err != nil {
			return nil, fmt.Errorf("saving template: %q: %w", args.tmpl, err)
		}
	}
	t := template.Must(template.ParseFiles(args.tmpl))
	return t, nil
}

// filename creates a filepath for the template filenames.
func (args *Args) templateCache() (err error) {
	l := args.layout.Pack()
	if l == "" {
		return fmt.Errorf("template cache %q: %w", args.layout, ErrNoLayout)
	}
	args.tmpl, err = gap.NewScope(gap.User, "retrotxt").DataPath(l + ".gohtml")
	if err != nil {
		return fmt.Errorf("template cache path: %q: %w", args.tmpl, err)
	}
	return nil
}

func (args *Args) templatePack() error {
	l := args.layout.Pack()
	if l == "" {
		return fmt.Errorf("template pack %q: %w", args.layout, ErrNoLayout)
	}
	args.pack = fmt.Sprintf("html/%s.gohtml", l)
	return nil
}

func (args *Args) templateData() (*[]byte, error) {
	b, err := static.Tmpl.ReadFile(args.pack)
	if err != nil {
		return nil, fmt.Errorf("template data %q: %w", args.pack, err)
	}
	return &b, nil
}

func (args *Args) templateSave() error {
	err := args.templatePack()
	if err != nil {
		return fmt.Errorf("template save pack error: %w", err)
	}
	b, err := args.templateData()
	if err != nil {
		return fmt.Errorf("template save data error: %w", err)
	}
	if _, _, err := filesystem.Save(args.tmpl, *b...); err != nil {
		return fmt.Errorf("template save error: %w", err)
	}
	return nil
}

func (args *Args) comment(lb filesystem.LB, r ...rune) string {
	l, w, f := 0, 0, "n/a"
	b, lbs, e := []byte(string(r)),
		filesystem.LineBreak(lb, false),
		args.Source.Encoding
	l, err := filesystem.Lines(bytes.NewReader(b), lb)
	if err != nil {
		l = -1
	}
	w, err = filesystem.Columns(bytes.NewReader(b), lb)
	if err != nil {
		w = -1
	}
	if args.Source.Name != "" {
		f = args.Source.Name
	}
	return fmt.Sprintf("encoding: %s; line break: %s; length: %d; width: %d; name: %s", e, lbs, l, w, f)
}

func (args *Args) fontFamily() string {
	if args.FontFamily.Flag {
		return args.FontFamily.Value
	}
	return viper.GetString("html.font.family")
}

func (args *Args) metaAuthor() string {
	if args.Metadata.Author.Flag {
		return args.Metadata.Author.Value
	}
	return viper.GetString("html.meta.author")
}

func (args *Args) metaColorScheme() string {
	if args.Metadata.ColorScheme.Flag {
		return args.Metadata.ColorScheme.Value
	}
	return viper.GetString("html.meta.color-scheme")
}

func (args *Args) metaDesc() string {
	if args.Metadata.Description.Flag {
		return args.Metadata.Description.Value
	}
	return viper.GetString("html.meta.description")
}

func (args *Args) metaKeywords() string {
	if args.Metadata.Keywords.Flag {
		return args.Metadata.Keywords.Value
	}
	return viper.GetString("html.meta.keywords")
}

func (args *Args) metaReferrer() string {
	if args.Metadata.Referrer.Flag {
		return args.Metadata.Referrer.Value
	}
	return viper.GetString("html.meta.referrer")
}

func (args *Args) metaRobots() string {
	if args.Metadata.Robots.Flag {
		return args.Metadata.Robots.Value
	}
	return viper.GetString("html.meta.robots")
}

func (args *Args) metaThemeColor() string {
	if args.Metadata.ThemeColor.Flag {
		return args.Metadata.ThemeColor.Value
	}
	return viper.GetString("html.meta.theme-color")
}

func (args *Args) pageTitle() string {
	if args.Title.Flag {
		return args.Title.Value
	}
	return viper.GetString("html.title")
}
