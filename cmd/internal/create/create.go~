package create

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrBody   = errors.New("could not parse the body flag")
	ErrCreate = errors.New("could not convert the text into a HTML document")
	ErrFlags  = errors.New("flags cannot be nil")
	ErrSrcNil = errors.New("src pointer cannot be nil")
)

func Run(cmd *cobra.Command, args []string) error {
	f := convert.Flag{
		Controls:  flag.Create().Controls,
		SwapChars: flag.Create().Swap,
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
	flag.Build = Strings(cmd, flag.Build)
	// handle standard input (stdio)
	serve := cmd.Flags().Lookup("serve").Changed
	if filesystem.IsPipe() {
		encode := flag.EncodingDefault
		if cp := cmd.Flags().Lookup("encode"); cp != nil {
			if cp.Changed {
				encode = cp.Value.String()
			}
		}
		return Pipe(encode, serve, f)
	}
	// handle the hidden --body flag value,
	// used for debugging, it ignores most other flags and
	// overrides the <pre></pre> content before exiting
	if body := cmd.Flags().Lookup("body"); body.Changed {
		b := []byte(body.Value.String())
		return Body(serve, &b)
	}
	if err := flag.Help(cmd, args...); err != nil {
		return err
	}
	return Files(cmd, f, args...)
}

// Body takes the supplied text and uses for the HTML <pre></pre> elements text content.
// Body is intended as a hidden debugging feature.
func Body(serve bool, b *[]byte) error {
	if !serve {
		err := flag.Build.Create(b)
		if err != nil {
			return fmt.Errorf("%s: %w", ErrBody, err)
		}
		return nil
	}
	if err := Serve(0, b); err != nil {
		return err
	}
	return nil
}

// Files parses the flags to create the HTML document or website.
// The generated HTML and associated files will either be served, saved or printed.
func Files(cmd *cobra.Command, flags convert.Flag, args ...string) error {
	args, conv, samp, err := flag.Args(cmd, args...)
	if err != nil {
		return err
	}
	for i, arg := range args {
		b, err := flag.ReadArgument(arg, cmd, conv, samp)
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
			continue
		}
		encode := flag.EncodingDefault
		if cp := cmd.Flags().Lookup("encode"); cp != nil {
			if cp.Changed {
				encode = cp.Value.String()
			}
		}
		r, err := Runes(encode, flags, &b)
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
			continue
		}
		if r == nil {
			continue
		}
		b2 := []byte(string(r))
		b, r = nil, nil
		serve := cmd.Flags().Lookup("serve").Changed
		if !serve {
			if err := flag.Build.Create(&b2); err != nil {
				return err
			}
			continue
		}
		if err := Serve(i, &b2); err != nil {
			return err
		}
	}
	return nil
}

// Pipe creates HTML content using the standard input (stdin) of the operating system.
func Pipe(encode string, serve bool, flags convert.Flag) error {
	src, err := filesystem.ReadPipe()
	if err != nil {
		return fmt.Errorf("%s: %w", logs.ErrPipeRead, err)
	}
	r, err := Runes(encode, flags, &src)
	if err != nil {
		return err
	}
	b := []byte(string(r))
	r = nil
	if !serve {
		if err := flag.Build.Create(&b); err != nil {
			return err
		}
		return nil
	}
	if err := Serve(0, &b); err != nil {
		return err
	}
	return nil
}

// Runes converts the src into UTF runes.
func Runes(encode string, flags convert.Flag, src *[]byte) ([]rune, error) {
	if src == nil {
		return nil, ErrSrcNil
	}
	conv := convert.Convert{Flags: flags}
	f := sample.Flags{}
	var err error
	// encode and convert the source text
	if encode != "" {
		if f.From, err = convert.Encoder(encode); err != nil {
			return nil, fmt.Errorf("%s: %w", logs.ErrEncode, err)
		}
		conv.Input.Encoding = f.From
	}
	// obtain any appended SAUCE metadata
	if flag.Build.SauceData.Use {
		flag.Build.SauceData = flag.SAUCE(src)
	}
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
	return r, nil
}

// SaveDst returns the directory the created HTML and other files will be saved to.
func SaveDst() (string, error) {
	var err error
	s := viper.GetString("save_directory")
	if s == "" {
		s, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("current working directory error: %v", err)
		}
	}
	return s, nil
}

// Serve hosts the HTML using an internal HTTP server.
func Serve(arg int, b *[]byte) error {
	if arg != 0 {
		// only ever serve the first file given to the args.
		// in the future, when handling multiple files a dynamic
		// index.html could be generated with links to each of the htmls.
		return nil
	}
	if err := flag.Build.Serve(b); err != nil {
		return err
	}
	return nil
}

// Strings handles the defaults for flags that accept strings.
// These flags are parse to three different states.
// 1) the flag is unchanged, so use the configured viper default.
// 2) the flag has a new value to overwrite viper default.
// 3) a blank flag value is given to overwrite viper default with an empty/disable value.
func Strings(cmd *cobra.Command, args create.Args) create.Args {
	changed := func(key string) bool {
		l := cmd.Flags().Lookup(key)
		if l == nil {
			return false
		}
		return l.Changed
	}
	args.FontFamily.Flag = changed("font-family")
	args.Metadata.Author.Flag = changed("meta-author")
	args.Metadata.ColorScheme.Flag = changed("meta-color-scheme")
	args.Metadata.Description.Flag = changed("meta-description")
	args.Metadata.Keywords.Flag = changed("meta-keywords")
	args.Metadata.Referrer.Flag = changed("meta-referrer")
	args.Metadata.Robots.Flag = changed("meta-robots")
	args.Metadata.ThemeColor.Flag = changed("meta-theme-color")
	args.Title.Flag = changed("title")
	ff := cmd.Flags().Lookup("font-family")
	if !ff.Changed {
		args.FontFamily.Value = "vga"
	}
	if args.FontFamily.Value == "" {
		args.FontFamily.Value = ff.Value.String()
	}
	return args
}
