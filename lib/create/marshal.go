package create

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/layout"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/static"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
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
	// templates are found in the dir static/html/*.gohtml
	switch args.Layouts {
	case layout.Inline:
		var err error
		p, err = args.marshalInline(&p)
		if err != nil {
			return PageData{}, err
		}
	case layout.Standard:
		p = args.marshalStandard(&p)
	case layout.Compact: // disables all meta tags
		p = args.marshalCompact(&p)
	case layout.None:
		// do nothing
	default:
		return PageData{}, fmt.Errorf("pagedata %s: %w", args.Layouts, logs.ErrTmplName)
	}
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

// TemplateSave saves an embedded template.
func (args *Args) TemplateSave() error {
	err := args.templatePack()
	if err != nil {
		return fmt.Errorf("template save pack error: %w", err)
	}
	b, err := args.templateData()
	if err != nil {
		return fmt.Errorf("template save data error: %w", err)
	}
	if _, _, err := filesystem.Write(args.Tmpl, *b...); err != nil {
		return fmt.Errorf("template save error: %w", err)
	}
	return nil
}

// marshalCompact is used by the compact layout argument.
func (args *Args) marshalCompact(p *PageData) PageData {
	p.PageTitle = args.pageTitle()
	p.MetaGenerator = false
	return *p
}

// marshalInline is used by the inline layout argument.
func (args *Args) marshalInline(p *PageData) (PageData, error) {
	*p = args.marshalStandard(p)
	p.ExternalEmbed = true
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	// styles
	s := bytes.TrimSpace(static.CSSStyles)
	// font
	f, err := FontCSS(args.FontFamily.Value, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return PageData{}, fmt.Errorf("pagedata font error: %w", err)
	}
	f = bytes.TrimSpace(f)
	// merge
	var embed []byte
	c := [][]byte{s, []byte("/* font */"), f}
	embed = bytes.Join(c, []byte("\n\n"))
	// compress & embed
	embed, err = m.Bytes("text/css", embed)
	if err != nil {
		return PageData{}, fmt.Errorf("pagedata minify css: %w", err)
	}
	p.CSSEmbed = template.CSS(string(embed))
	jsp := static.Scripts
	jsp, err = m.Bytes("application/javascript", jsp)
	if err != nil {
		return PageData{}, fmt.Errorf("pagedata minify javascript: %w", err)
	}
	p.ScriptEmbed = template.JS(string(jsp)) // nolint:gosec
	return *p, nil
}

// marshalStandard is used by the standard layout argument.
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
	p.BuildVersion = metaAppVersion()
	return *p
}

func metaAppVersion() string {
	t := time.Now().UTC()
	date := t.Format(time.RFC3339)
	if !meta.IsGoBuild() {
		return fmt.Sprintf("%s %s; %s",
			meta.Name, meta.Print(), date)
	}
	return fmt.Sprintf("%s; %s", meta.Name, date)
}

// marshalTextTransform marshals the bytes data with the HTML template.
func (args *Args) marshalTextTransform(b *[]byte) (bytes.Buffer, error) {
	buf := bytes.Buffer{}
	Tmpl, err := args.newTemplate()
	if err != nil {
		return buf, fmt.Errorf("stdout new template failure: %w", err)
	}
	d, err := args.Marshal(b)
	if err != nil {
		return buf, fmt.Errorf("stdout meta and pagedata failure: %w", err)
	}
	if err = Tmpl.Execute(&buf, d); err != nil {
		return buf, fmt.Errorf("stdout template execute failure: %w", err)
	}
	return buf, nil
}

// newTemplate creates and parses a new template file.
func (args *Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, fmt.Errorf("using existing template cache: %w", err)
	}
	if !args.Save.Cache {
		if err := args.TemplateSave(); err != nil {
			return nil, fmt.Errorf("creating a new template: %w", err)
		}
	} else {
		s, err := os.Stat(args.Tmpl)
		switch {
		case os.IsNotExist(err):
			if err2 := args.TemplateSave(); err2 != nil {
				return nil, fmt.Errorf("saving to the template: %w", err2)
			}
		case err != nil:
			return nil, err
		case s.IsDir():
			return nil, fmt.Errorf("new template %q: %w", args.Tmpl, logs.ErrTmplIsDir)
		}
	}
	// to avoid a potential panic, Stat again in case os.IsNotExist() is true
	s, err := os.Stat(args.Tmpl)
	if err != nil {
		return nil, fmt.Errorf("could not access file: %q: %w", args.Tmpl, err)
	}
	if err = args.templatePack(); err != nil {
		return nil, fmt.Errorf("template pack: %w", err)
	}
	b, err := args.templateData()
	if s.Size() != int64(len(*b)) {
		if err != nil {
			return nil, fmt.Errorf("template data: %w", err)
		}
		if _, _, err := filesystem.Write(args.Tmpl, *b...); err != nil {
			return nil, fmt.Errorf("saving template: %q: %w", args.Tmpl, err)
		}
	}
	t := template.Must(template.ParseFiles(args.Tmpl))
	return t, nil
}

// templateCache creates a filepath for the cache templates.
func (args *Args) templateCache() error {
	const ext = ".gohtml"
	name := args.Layouts.Pack()
	if name == "" {
		return fmt.Errorf("template cache %q: %w", args.Layouts, logs.ErrTmplName)
	}
	var err error
	filename := name + ext
	args.Tmpl, err = gap.NewScope(gap.User, meta.Dir).DataPath(filename)
	if err != nil {
		return fmt.Errorf("template cache path: %q: %w", args.Tmpl, err)
	}
	return nil
}

// templatePack creates a filepath for the embedded templates.
func (args *Args) templatePack() error {
	// sep should be kept as-is, regardless of platform
	const dir, ext, sep = "html", ".gohtml", "/"
	name := args.Layouts.Pack()
	if name == "" {
		return fmt.Errorf("template pack %q: %w", args.Layouts, logs.ErrTmplName)
	}
	filename := name + ext
	args.Pack = strings.Join([]string{dir, filename}, sep)
	return nil
}

// templateData reads and returns an embedded file.
func (args *Args) templateData() (*[]byte, error) {
	b, err := static.Tmpl.ReadFile(args.Pack)
	if err != nil {
		return nil, fmt.Errorf("template data %s: %w", args.Pack, err)
	}
	return &b, nil
}

// fontFamily value for the CSS font face.
func (args *Args) fontFamily() string {
	if args.FontFamily.Flag {
		return args.FontFamily.Value
	}
	return viper.GetString("html.font.family")
}

// metaAuthor content for the meta sauce-data attribute.
func (args *Args) metaAuthor() string {
	if args.Metadata.Author.Flag {
		return args.Metadata.Author.Value
	}
	return viper.GetString("html.meta.author")
}

// metaColorScheme content for the meta color-scheme attribute.
func (args *Args) metaColorScheme() string {
	if args.Metadata.ColorScheme.Flag {
		return args.Metadata.ColorScheme.Value
	}
	return viper.GetString("html.meta.color-scheme")
}

// metaAuthor content for the meta sauce-description attribute.
func (args *Args) metaDesc() string {
	if args.Metadata.Description.Flag {
		return args.Metadata.Description.Value
	}
	return viper.GetString("html.meta.description")
}

// metaKeywords content for the meta keywords attribute.
func (args *Args) metaKeywords() string {
	if args.Metadata.Keywords.Flag {
		return args.Metadata.Keywords.Value
	}
	return viper.GetString("html.meta.keywords")
}

// metaReferrer content for the meta referrer attribute.
func (args *Args) metaReferrer() string {
	if args.Metadata.Referrer.Flag {
		return args.Metadata.Referrer.Value
	}
	return viper.GetString("html.meta.referrer")
}

// metaRobots content for the meta robots attribute.
func (args *Args) metaRobots() string {
	if args.Metadata.Robots.Flag {
		return args.Metadata.Robots.Value
	}
	return viper.GetString("html.meta.robots")
}

// metaRobots content for the meta theme-color attribute.
func (args *Args) metaThemeColor() string {
	if args.Metadata.ThemeColor.Flag {
		return args.Metadata.ThemeColor.Value
	}
	return viper.GetString("html.meta.theme-color")
}

// pageTitle value for the title element.
func (args *Args) pageTitle() string {
	if args.Title.Flag {
		return args.Title.Value
	}
	return viper.GetString("html.title")
}
