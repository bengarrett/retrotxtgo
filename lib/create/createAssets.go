package create

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"golang.org/x/text/language"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/humanize"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/static"
)

const (
	nameCSS  = "styles.css"
	nameFont = "font.css"
	nameHTML = "index.html"
	nameJS   = "scripts.js"
	nameFav  = "favicon.ico"
)

// saveCSS creates and saves the styles stylesheet to the Destination argument.
func (args *Args) saveCSS(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
	}
	name, err := args.destination(nameCSS)
	if err != nil {
		c <- err
	}
	b := static.Styles
	if len(b) == 0 {
		c <- fmt.Errorf("create.saveCSS %q: %w", args.pack, static.ErrNotFound)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

func (args *Args) saveFavIcon(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
	}
	name, err := args.destination(nameFav)
	if err != nil {
		c <- err
	}
	b, err := static.Image.ReadFile("retrotxt_16.png")
	if err != nil {
		c <- fmt.Errorf("create.saveFavIcon %q: %w", args.pack, err)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

// saveFont unpacks and saves the font binary to the Destination argument.
func (args *Args) saveFont(c chan error) {
	if !args.FontEmbed {
		f := Family(args.FontFamily.Value)
		if f.String() == "" {
			c <- fmt.Errorf("save font, could not save %q: %w", args.FontFamily.Value, ErrUnknownFF)
			return
		}
		if err := args.saveFontWoff2(f.File(), "font/"+f.File()); err != nil {
			c <- err
		}
	}
	switch args.layout {
	case Standard:
		if err := args.saveFontCSS("font.css"); err != nil {
			c <- err
		}
	case Compact, Inline, None:
	}
	c <- nil
}

func (args *Args) saveFontCSS(name string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	f := Family(args.FontFamily.Value).String()
	if f == "" {
		return fmt.Errorf("create.saveFontCSS %q: %w", name, ErrUnknownFF)
	}
	b, err := FontCSS(f, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return err
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	bytesStats(name, nn)
	return nil
}

func (args *Args) saveFontWoff2(name, packName string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	b, err := static.Font.ReadFile(packName)
	if err != nil {
		return fmt.Errorf("create.saveFontWoff2 %q: %w", args.pack, err)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	bytesStats(name, nn)
	return nil
}

func (args *Args) saveJS(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
		return
	}
	name, err := args.destination(nameJS)
	if err != nil {
		c <- err
	}
	b := static.Scripts
	if len(b) == 0 {
		c <- fmt.Errorf("create.saveJS %q: %w", args.pack, static.ErrNotFound)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

// SaveHTML creates and saves the html template to the Destination argument.
func (args *Args) saveHTML(b *[]byte, c chan error) {
	name, err := args.destination(nameHTML)
	if err != nil {
		c <- err
	}
	if name == "" {
		c <- ErrEmptyName
	}
	file, err := os.Create(name)
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
	bytesStats(name, nn)
	if err := w.Flush(); err != nil {
		c <- err
	}
	c <- nil
}

func bytesStats(name string, nn int) {
	const kB = 1000
	h := humanize.Decimal(int64(nn), language.AmericanEnglish)
	color.OpFuzzy.Printf("saved to %s", name)
	switch {
	case nn == 0:
		color.OpFuzzy.Printf("saved to %s (zero-byte file)", name)
	case nn < kB:
		color.OpFuzzy.Printf(", %s", h)
	default:
		color.OpFuzzy.Printf(", %s (%d)", h, nn)
	}
	fmt.Print("\n")
}

func (args *Args) destination(name string) (string, error) {
	dir := filesystem.DirExpansion(args.Save.Destination)
	path := dir
	stat, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("args destination directory failed %q: %w", dir, err)
	}
	if stat.IsDir() {
		path = filepath.Join(dir, name)
	}
	_, err = os.Stat(path)
	if !args.Save.OW && !os.IsNotExist(err) {
		switch name {
		case nameFav, nameJS, "vga.woff2", "mona.woff2":
			// existing static files can be ignored
			return path, nil
		}
		logs.Println("file exists", path, ErrReqOW)
		return path, nil
	}
	if os.IsNotExist(err) {
		empty := []byte{}
		if _, _, err = filesystem.Save(path, empty...); err != nil {
			return "", fmt.Errorf("args destination path failed %q: %w", path, err)
		}
	}
	return path, nil
}
