package create

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/humanize"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/gookit/color"
	"golang.org/x/text/language"
)

type asset int

const (
	htmlFn asset = iota
	fontFn
	cssFn
	jsFn
	favFn
	bbsFn
	pcbFn

	faviconSrc = "img/retrotxt_16.png"
)

func (a asset) write() string {
	// do not change the order of this array,
	// they must match the Fn asset iota values.
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
	c <- args.saveCSSFile(static.CSSStyles, cssFn)
}

func (args *Args) saveBBS(c chan error) {
	c <- args.saveCSSFile(static.CSSBBS, bbsFn)
}

func (args *Args) savePCBoard(c chan error) {
	c <- args.saveCSSFile(static.CSSPCBoard, pcbFn)
}

func (args *Args) saveCSSFile(src []byte, a asset) error {
	s, err := args.destination(a.write())
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	if len(src) == 0 {
		return fmt.Errorf("%s, %w", a.write(), static.ErrNotFound)
	}
	nn, _, err := filesystem.Save(s, src...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(bytesStats(s, nn))
	}
	return nil
}

// saveFavIcon read and save the favorite icon to a file.
func (args *Args) saveFavIcon(c chan error) {
	s, err := args.destination(favFn.write())
	if err != nil {
		c <- fmt.Errorf("%w: %s", err, s)
	}
	b, err := static.Image.ReadFile(faviconSrc)
	if err != nil {
		c <- fmt.Errorf("saveFavIcon, %w", err)
	}
	nn, _, err := filesystem.Save(s, b...)
	if err != nil {
		c <- err
	}
	if !args.Test {
		fmt.Println(bytesStats(s, nn))
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
	if err := args.saveFontCSS(fontFn.write()); err != nil {
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
	nn, _, err := filesystem.Save(s, b...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(bytesStats(s, nn))
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
		return fmt.Errorf("saveFontWoff2 %q, %w", args.pack, err)
	}
	nn, _, err := filesystem.Save(s, b...)
	if err != nil {
		return err
	}
	if !args.Test {
		fmt.Println(bytesStats(s, nn))
	}
	return nil
}

// saveJS creates and save the JS file.
func (args *Args) saveJS(c chan error) {
	s, err := args.destination(jsFn.write())
	if err != nil {
		c <- fmt.Errorf("%w: %s", err, s)
	}
	b := static.Scripts
	if len(b) == 0 {
		c <- fmt.Errorf("saveJS, %w", static.ErrNotFound)
	}
	nn, _, err := filesystem.Save(s, b...)
	if err != nil {
		c <- err
	}
	if !args.Test {
		fmt.Println(bytesStats(s, nn))
	}
	c <- nil
}

// saveHTML creates and save the HTML file.
func (args *Args) saveHTML(b *[]byte, c chan error) {
	s, err := args.destination(htmlFn.write())
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
		fmt.Println(bytesStats(s, nn))
	}
	if err := w.Flush(); err != nil {
		c <- err
	}
	c <- nil
}

// bytesStats humanizes, colorizes and prints the filename and size.
func bytesStats(name string, nn int) string {
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
		case favFn.write(), jsFn.write(), "vga.woff2", "mona.woff2":
			// existing static files can be ignored
			return dst, nil
		}
		return dst, ErrFileExist
	}
	// create an empty zero byte file
	if os.IsNotExist(err) {
		empty := []byte{}
		if _, _, err = filesystem.Save(dst, empty...); err != nil {
			return "", fmt.Errorf("create %s, %w", dst, err)
		}
	}
	return dst, nil
}
