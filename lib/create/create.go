// Package create makes HTML and other web resources from a text file.
package create

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/assets"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/layout"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

type Args assets.Args

// ColorScheme values for the content attribute of <meta name="color-scheme">.
func ColorScheme() [3]string {
	return [...]string{"normal", "dark light", "only light"}
}

// Referrer values for the content attribute of <meta name="referrer">.
func Referrer() [8]string {
	return [...]string{
		"no-referrer", "origin", "no-referrer-when-downgrade",
		"origin-when-cross-origin", "same-origin", "strict-origin",
		"strict-origin-when-cross-origin", "unsafe-URL",
	}
}

// Robots values for the content attribute of <meta name="robots">.
func Robots() [9]string {
	return [...]string{
		"index", "noindex", "follow", "nofollow", "none",
		"noarchive", "nosnippet", "noimageindex", "nocache",
	}
}

// Layouts are the names of the HTML templates.
func Layouts() []string {
	return []string{
		layout.Standard.String(),
		layout.Inline.String(),
		layout.Compact.String(),
		layout.None.String(),
	}
}

// Create handles the target output command arguments.
func (args *Args) Create(b *[]byte) error {
	var err error
	args.Layouts, err = layout.ParseLayout(args.Layout)
	if err != nil {
		return err
	}
	switch {
	case args.Save.AsFiles:
		if err := args.SaveAssets(b); err != nil {
			// --overwrite hint
			if errors.As(err, &ErrFileExist) {
				fmt.Println(logs.Hint("create [filenames] --overwrite", ErrFileExist))
				fmt.Println(str.Info() + "Use the overwrite flag to replace any existing files.")
				os.Exit(logs.OSErrCode)
			}
			return nil
		}
	case args.Save.Compress:
		const noDestination = ""
		args.ZipAssets(noDestination, b)
	default:
		// print to terminal
		if err := args.Stdout(b); err != nil {
			return err
		}
	}
	return nil
}

func (args *Args) SaveAssets(b *[]byte) error {
	skip := func(c chan error) {
		c <- nil
	}
	if args.Save.Destination == "" {
		dir := []string{viper.GetString("save-directory")}
		var err error
		if args.Save.Destination, err = assets.Destination(dir...); err != nil {
			logs.FatalMark(args.Save.Destination, logs.ErrFileSaveD, err)
		}
	}

	r := bytes.NewReader(*b)
	args.Source.BBSType = bbs.Find(r)

	ch, cnt := make(chan error), 0

	go args.SaveHTML(b, ch)

	if layout.UseCSS(args.Layouts) {
		cnt++
		go args.saveStyles(ch)
	}
	if usePCBoard(args.Source.BBSType) {
		cnt += 2
		go args.saveBBS(ch)
		go args.savePCBoard(ch)
	}
	if layout.UseFontCSS(args.Layouts) {
		cnt++
		go args.saveFont(ch)
	}
	if layout.UseJS(args.Layouts) {
		cnt++
		go args.saveJS(ch)
	}
	if layout.UseIcon(args.Layouts) {
		cnt++
		go args.saveFavIcon(ch)
	}

	const optionalCh = 6
	skips := optionalCh - cnt
	for i := 0; i < skips; i++ {
		go skip(ch)
	}
	return check(<-ch, <-ch, <-ch, <-ch, <-ch, <-ch, <-ch)
}

func check(ch ...error) error {
	var errs error
	for _, err := range ch {
		errs = appendErr(errs, err)
	}
	return errs
}

func appendErr(errs, err error) error {
	// handle first error
	if errs == nil {
		return err
	}
	// skip duplicate errors
	if !errors.Is(err, ErrFileExist) && errors.As(errs, &err) {
		return errs
	}
	return fmt.Errorf("%s;%w", errs, err)
}

func usePCBoard(b bbs.BBS) bool {
	switch b {
	case bbs.PCBoard:
		return true
	case bbs.ANSI, bbs.Celerity, bbs.Renegade, bbs.Telegard, bbs.WWIVHash, bbs.WWIVHeart, bbs.Wildcat:
		return false
	default:
		return false
	}
}

// zipAssets compresses all assets into a single zip archive.
// An empty destination directory argument will save the zip file to the user working directory.
func (args *Args) ZipAssets(destDir string, b *[]byte) {
	defer func() {
		dir := args.Save.Destination
		m, err := filepath.Match(filepath.Join(os.TempDir(), "*"), dir)
		if err != nil {
			logs.FatalMark("*", logs.ErrTmpSaveD, err)
		}
		if m {
			if err = os.RemoveAll(dir); err != nil {
				logs.FatalMark(dir, logs.ErrTmpRMD, err)
			}
		}
	}()
	var err error
	args.Save.Destination, err = ioutil.TempDir(os.TempDir(), "*-zip")
	if err != nil {
		logs.FatalMark("temporary", logs.ErrFileSaveD, err)
	}
	if err = args.SaveAssets(b); err != nil {
		fmt.Println(logs.SprintWrap(logs.ErrFileSave, err))
		return
	}
	name := layout.ZipName
	if destDir != "" {
		name = filepath.Join(destDir, layout.ZipName)
	}
	zip := filesystem.Zip{
		Name:      name,
		Root:      args.Save.Destination,
		Comment:   "",
		Overwrite: args.Save.OW,
		Quiet:     args.Test,
	}
	if err = zip.Create(); err != nil {
		logs.FatalMark(name, logs.ErrZipFile, err)
	}
}

// Stdout creates and prints the HTML template.
func (args *Args) Stdout(b *[]byte) error {
	// html
	html, err := args.marshalTextTransform(b)
	if err != nil {
		return fmt.Errorf("stdout: %w", err)
	}
	// font css
	ff := args.FontFamily.Value
	f := Family(ff).String()
	if f == "" {
		return fmt.Errorf("create.saveFontCSS %q: %w", ff, ErrFont)
	}
	font, err := FontCSS(f, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return err
	}
	// print assets
	if errj := args.printJS(&static.Scripts); errj != nil {
		return errj
	}
	if errc := args.printCSS(&static.CSSStyles); errc != nil {
		return errc
	}
	if errf := args.printFontCSS(f, &font); errf != nil {
		return errf
	}
	// always print the HTML
	fmt.Printf("\nHTML file: %s\n\n", HTML.Write())
	if err = str.Highlight(html.String(), "html", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout html highlight: %w", err)
	}
	return nil
}

func colorSyntax(s string) bool {
	switch s {
	case "", "none":
		return false
	}
	return str.Valid(s)
}

func (args *Args) printCSS(b *[]byte) error {
	if !layout.UseCSS(args.Layouts) {
		return nil
	}
	fmt.Printf("\nCSS file: %s\n\n", StyleCss.Write())
	if !colorSyntax(args.Syntax) {
		fmt.Println(string(*b))
		return nil
	}
	if err := str.Highlight(string(*b), "css", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout css highlight: %w", err)
	}
	return nil
}

func (args *Args) printFontCSS(name string, b *[]byte) error {
	if !layout.UseFontCSS(args.Layouts) {
		return nil
	}
	fmt.Printf("\nCSS for %s font file: %s\n\n", name, FontCss.Write())
	if err := str.Highlight(string(*b), "css", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout font css highlight: %w", err)
	}
	return nil
}

func (args *Args) printJS(b *[]byte) error {
	if !layout.UseJS(args.Layouts) {
		return nil
	}
	fmt.Printf("\nJS file: %s\n\n", Scripts.Write())
	if !colorSyntax(args.Syntax) {
		fmt.Println(string(*b))
		return nil
	}
	if err := str.Highlight(string(*b), "js", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout js highlight: %w", err)
	}
	return nil
}

// Normalize runes into bytes by making adjustments to text control codes.
func Normalize(e encoding.Encoding, r ...rune) []byte {
	switch e {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		s, _, err := transform.String(replaceNELs(), string(r))
		if err != nil {
			return []byte(string(r))
		}
		return []byte(s)
	}
	return []byte(string(r))
}

// replaceNELs replace EBCDIC newlines with Unicode linefeeds.
func replaceNELs() runes.Transformer {
	return runes.Map(func(r rune) rune {
		if r == filesystem.NextLine {
			return filesystem.Linefeed
		}
		return r
	})
}
