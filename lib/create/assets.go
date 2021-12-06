package create

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/bengarrett/sauce/humanize"
	"github.com/gookit/color"
	"golang.org/x/text/language"
)

// Asset filenames.
type Asset int

//nolint:stylecheck,revive
const (
	HTML     Asset = iota // Index html.
	FontCss               // CSS containing fonts.
	StyleCss              // CSS containing styles and colors.
	Scripts               // JS scripts.
	FavIco                // Favorite icon.
	BbsCss                // Other BBS CSS.
	PcbCss                // PCBoard BBS CSS.
)

func (a Asset) Write() string {
	// do not change the order of this array, they must match the Asset iota values.
	return [...]string{
		// core assets
		"index.html",
		"font.css",
		"styles.css",
		"scripts.js",
		"favicon.ico",
		// dynamic assets
		"text_bbs.css",
		"text_pcboard.css",
	}[a]
}

// saveStyles creates and save the styles CSS file.
func (args *Args) saveStyles(c chan error) {
	c <- args.saveCSSFile(static.CSSStyles, StyleCss)
}

func (args *Args) saveBBS(c chan error) {
	c <- args.saveCSSFile(static.CSSBBS, BbsCss)
}

func (args *Args) savePCBoard(c chan error) {
	c <- args.saveCSSFile(static.CSSPCBoard, PcbCss)
}

func (args *Args) saveCSSFile(src []byte, a Asset) error {
	s, err := args.destination(a.Write())
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	if len(src) == 0 {
		return fmt.Errorf("%s, %w", a.Write(), static.ErrNotFound)
	}
	nn, _, err := filesystem.Write(s, src...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	return nil
}

// saveFavIcon read and save the favorite icon to a file.
func (args *Args) saveFavIcon(c chan error) {
	const faviconSrc = "img/retrotxt_16.png"
	s, err := args.destination(FavIco.Write())
	if err != nil {
		c <- fmt.Errorf("%w: %s", err, s)
	}
	b, err := static.Image.ReadFile(faviconSrc)
	if err != nil {
		c <- fmt.Errorf("saveFavIcon, %w", err)
	}
	nn, _, err := filesystem.Write(s, b...)
	if err != nil {
		c <- err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	c <- nil
}

// saveFont read and save the font to a file.
func (args *Args) saveFont(c chan error) {
	if !args.FontEmbed {
		f := Family(args.FontFamily.Value)
		if f.String() == "" {
			c <- fmt.Errorf("saveFont %s, %w", args.FontFamily.Value, ErrFont)
			return
		}
		path := "font/" + f.File()
		if err := args.saveFontWoff2(f.File(), path); err != nil {
			c <- err
		}
	}
	if err := args.saveFontCSS(FontCss.Write()); err != nil {
		c <- err
	}
	c <- nil
}

// saveFontCSS creates and save the font styles CSS file.
func (args *Args) saveFontCSS(name string) error {
	s, err := args.destination(name)
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	f := Family(args.FontFamily.Value).String()
	if f == "" {
		return fmt.Errorf("saveFontCSS %s, %w", name, ErrFont)
	}
	b, err := FontCSS(f, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return err
	}
	nn, _, err := filesystem.Write(s, b...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	return nil
}

// saveFontWoff2 read and save the WOFF2 font binary to a file.
func (args *Args) saveFontWoff2(name, packName string) error {
	s, err := args.destination(name)
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	b, err := static.Font.ReadFile(packName)
	if err != nil {
		return fmt.Errorf("saveFontWoff2 %q, %w", args.Pack, err)
	}
	nn, _, err := filesystem.Write(s, b...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	return nil
}

// saveJS creates and save the JS file.
func (args *Args) saveJS(c chan error) {
	s, err := args.destination(Scripts.Write())
	if err != nil {
		c <- fmt.Errorf("%w: %s", err, s)
	}
	b := static.Scripts
	if len(b) == 0 {
		c <- fmt.Errorf("saveJS, %w", static.ErrNotFound)
	}
	nn, _, err := filesystem.Write(s, b...)
	if err != nil {
		c <- err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	c <- nil
}

// SaveHTML creates and save the HTML file.
func (args *Args) SaveHTML(b *[]byte, c chan error) {
	s, err := args.destination(HTML.Write())
	if err != nil {
		c <- fmt.Errorf("%w: %s", err, s)
	}
	if s == "" {
		c <- ErrFileNil
	}
	// check directory
	file, err := os.Create(s)
	if err != nil {
		c <- err
	}
	defer func() {
		cerr := file.Close()
		c <- cerr
	}()
	buf, err := args.marshalTextTransform(b)
	if err != nil {
		c <- err
	}
	w := bufio.NewWriter(file)
	nn, err := w.Write(buf.Bytes())
	if err != nil {
		c <- err
	}
	if !args.Test {
		fmt.Println(Stats(s, nn))
	}
	if err := w.Flush(); err != nil {
		c <- err
	}
	c <- nil
}

// Stats humanizes, colorizes and prints the filename and size.
func Stats(name string, nn int) string {
	const kB = 1000
	if nn == 0 {
		return color.OpFuzzy.Sprintf("saved to %s (zero-byte file)", name)
	}
	h := humanize.Decimal(int64(nn), language.AmericanEnglish)
	s := color.OpFuzzy.Sprintf("saved to %s", name)
	switch {
	case nn < kB:
		s += color.OpFuzzy.Sprintf(", %s", h)
	default:
		s += color.OpFuzzy.Sprintf(", %s (%d)", h, nn)
	}
	return s
}

// destination validate and returns the path of the named file.
func (args *Args) destination(name string) (string, error) {
	dir := filesystem.DirExpansion(args.Save.Destination)
	stat, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("destination %s: %w", dir, err)
	}
	dst := ""
	if stat.IsDir() {
		dst = filepath.Join(dir, name)
	}
	fs, err := os.Stat(dst)
	if os.IsNotExist(err) {
		// expected, dst doesn't exist
	} else if err != nil {
		// unexpected, some other system err
		return "", err
	} else if fs.Size() > 0 && !args.Save.OW {
		// unexpected, dst does exist
		switch name {
		case FavIco.Write(), Scripts.Write(), "vga.woff2", "mona.woff2":
			// existing static files can be ignored
			return dst, nil
		}
		return dst, ErrFileExist
	}
	// create an empty zero byte file
	if os.IsNotExist(err) {
		empty := []byte{}
		if _, _, err = filesystem.Write(dst, empty...); err != nil {
			return "", fmt.Errorf("create %s, %w", dst, err)
		}
	}
	return dst, nil
}
