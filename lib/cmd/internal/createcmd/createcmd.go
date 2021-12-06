package createcmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrCreate = errors.New("could not convert the text into a HTML document")
	ErrBody   = errors.New("could not parse the body flag")
)

func Run(cmd *cobra.Command, args []string) error {
	f := convert.Flag{
		Controls:  flag.CreateDefaults.Controls,
		SwapChars: flag.CreateDefaults.Swap,
	}
	// handle defaults, use these control codes
	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		f.Controls = []string{"eof", "tab"}
	}
	// handle defaults, swap out these characters with UTF-8 alternatives
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		f.SwapChars = []string{"null", "bar"}
	}
	// handle the defaults for most other flags
	Strings(cmd)
	// handle standard input (stdio)
	if filesystem.IsPipe() {
		return ParsePipe(cmd, f)
	}
	// handle the hidden --body flag value,
	// used for debugging, it ignores most other flags and
	// overrides the <pre></pre> content before exiting
	if body := cmd.Flags().Lookup("body"); body.Changed {
		return ParseBody(cmd)
	}
	if err := flag.PrintUsage(cmd, args...); err != nil {
		return err
	}
	return ParseFiles(cmd, f, args...)
}

// HTML applies a HTML template to the src text.
func HTML(cmd *cobra.Command, flags convert.Flag, src *[]byte) ([]byte, error) {
	var err error
	conv := convert.Convert{
		Flags: flags,
	}
	f := sample.Flags{}
	// encode and convert the source text
	if cp := cmd.Flags().Lookup("encode"); cp != nil {
		name := flag.EncodingDefault
		if cp.Changed {
			name = cp.Value.String()
		}
		if f.From, err = convert.Encoder(name); err != nil {
			return nil, fmt.Errorf("%s: %w", logs.ErrEncode, err)
		}
		conv.Input.Encoding = f.From
	}
	// obtain any appended SAUCE metadata
	flag.SAUCE(src)
	// convert the source text into web friendly UTF8
	var r []rune
	if flag.EndOfFile(conv.Flags) {
		r, err = conv.Text(*src...)
	} else {
		r, err = conv.Dump(*src...)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreate, err)
	}
	return []byte(string(r)), nil
}

// ParsePipe creates HTML content using the standard input (stdin) of the operating system.
func ParsePipe(cmd *cobra.Command, flags convert.Flag) error {
	src, err := filesystem.ReadPipe()
	if err != nil {
		return fmt.Errorf("%s: %w", logs.ErrPipeRead, err)
	}
	b, err := HTML(cmd, flags, &src)
	if err != nil {
		return err
	}
	serve := cmd.Flags().Lookup("serve").Changed
	h, err := ServeBytes(0, serve, &b)
	if err != nil {
		return err
	}
	if !h {
		if err := flag.HTML.Create(&b); err != nil {
			return err
		}
	}
	return nil
}

// hidden is a hidden debugging feature.
// It takes the supplied text and uses for the HTML <pre></pre> elements text content.
func ParseBody(cmd *cobra.Command) error {
	// hidden --body flag ignores most other args
	if body := cmd.Flags().Lookup("body"); body.Changed {
		b := []byte(body.Value.String())
		serve := cmd.Flags().Lookup("serve").Changed
		h, err := ServeBytes(0, serve, &b)
		if err != nil {
			return err
		}
		if !h {
			err := flag.HTML.Create(&b)
			if err != nil {
				return fmt.Errorf("%s: %w", ErrBody, err)
			}
		}
	}
	return nil
}

// ParseFiles parses the flags to create the HTML document or website.
// The generated HTML and associated files will either be served, saved or printed.
func ParseFiles(cmd *cobra.Command, flags convert.Flag, args ...string) error {
	args, conv, samp, err := flag.InitArgs(cmd, args...)
	if err != nil {
		return err
	}
	for i, arg := range args {
		b, err := flag.ReadArg(arg, cmd, conv, samp)
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
			continue
		}
		b, err = HTML(cmd, flags, &b)
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
			continue
		}
		if b == nil {
			continue
		}
		h, err := ServeBytes(i, cmd.Flags().Lookup("serve").Changed, &b)
		if err != nil {
			return err
		}
		if !h {
			if err := flag.HTML.Create(&b); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveDir returns the directory the created HTML and other files will be saved to.
func SaveDir() string {
	var err error
	s := viper.GetString("save-directory")
	if s == "" {
		s, err = os.Getwd()
		if err != nil {
			fmt.Printf("current working directory error: %v\n", err)
		}
	}
	return s
}

// ServeBytes hosts the HTML using an internal HTTP server.
func ServeBytes(i int, changed bool, b *[]byte) (bool, error) {
	if i != 0 {
		return false, nil
		// only ever serve the first file given to the args.
		// in the future, when handling multiple files a dynamic
		// index.html could be generated with links to each of the htmls.
	}
	if changed {
		if err := flag.HTML.Serve(b); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Strings handles the defaults for flags that accept strings.
// These flags are parse to three different states.
// 1) the flag is unchanged, so use the configured viper default.
// 2) the flag has a new value to overwrite viper default.
// 3) a blank flag value is given to overwrite viper default with an empty/disable value.
func Strings(cmd *cobra.Command) {
	changed := func(key string) bool {
		l := cmd.Flags().Lookup(key)
		if l == nil {
			return false
		}
		return l.Changed
	}
	flag.HTML.FontFamily.Flag = changed("font-family")
	flag.HTML.Metadata.Author.Flag = changed("meta-author")
	flag.HTML.Metadata.ColorScheme.Flag = changed("meta-color-scheme")
	flag.HTML.Metadata.Description.Flag = changed("meta-description")
	flag.HTML.Metadata.Keywords.Flag = changed("meta-keywords")
	flag.HTML.Metadata.Referrer.Flag = changed("meta-referrer")
	flag.HTML.Metadata.Robots.Flag = changed("meta-robots")
	flag.HTML.Metadata.ThemeColor.Flag = changed("meta-theme-color")
	flag.HTML.Title.Flag = changed("title")
	ff := cmd.Flags().Lookup("font-family")
	if !ff.Changed {
		flag.HTML.FontFamily.Value = "vga"
	}
	if flag.HTML.FontFamily.Value == "" {
		flag.HTML.FontFamily.Value = ff.Value.String()
	}
}
