
## Developer notes

These will eventually be removed.

### References

#### Text files
- [SAUCE](http://www.acid.org/info/sauce/sauce.htm) (ANSI/ASCII metadata)

#### Go libraries
- [Cobra](https://pkg.go.dev/github.com/spf13/cobra), CLI interface
- [Viper](https://pkg.go.dev/mod/github.com/spf13/viper), configuration file
- [chroma](https://github.com/alecthomas/chroma), general purpose syntax highlighter
- [gookit/color](https://github.com/gookit/color), terminal color rendering library
- [godirwalk](https://github.com/karrick/godirwalk), fast directory traversal, last update Aug, 2020
- [go-app-paths](https://github.com/muesli/go-app-paths), retrieve platform-specific paths
- [minify](https://github.com/tdewolff/minify), minifiers for web formats
- [mimemagic](https://github.com/zRedShift/mimemagic), powerful and versatile MIME sniffing package

#### Go packages

#### [Standard library](https://pkg.go.dev/std)

- [Sub-repositories](https://pkg.go.dev/golang.org/x/)

#### Unicode
- [unicode](https://golang.org/pkg/unicode/) Package unicode provides data and functions to test some properties of Unicode code points.
- [utf8](https://golang.org/pkg/unicode/utf8/) Package utf8 implements functions and constants to support text encoded in UTF-8. It includes functions to translate between runes and UTF-8 byte sequences.
- [utf16](https://pkg.go.dev/unicode/utf16) Package utf16 implements encoding and decoding of UTF-16 sequences.
- [x-utf32](https://pkg.go.dev/golang.org/x/text/encoding/unicode/utf32) Package utf32 provides the UTF-32 Unicode encoding.
---
- [x-unicode/norm](https://pkg.go.dev/golang.org/x/text/unicode/norm) Package norm contains types and functions for normalizing Unicode strings.
- [x-unicode/rangetable](https://pkg.go.dev/golang.org/x/text/unicode/rangetable) Package rangetable provides utilities for creating and inspecting unicode.RangeTables.
- [x-runes](https://pkg.go.dev/golang.org/x/text/runes) Package runes provide transforms for UTF-8 encoded text.
- [x-unicode/runenames](https://pkg.go.dev/golang.org/x/text/unicode/runenames) Package runenames provides rune names from the Unicode Character Database. For example, the name for '\u0100' is "LATIN CAPITAL LETTER A WITH MACRON".


#### Text
- [text](https://pkg.go.dev/golang.org/x/text) Text is a repository of text-related packages related to internationalization (i18n) and localization (l10n), such as character encodings, text transformations, and locale-specific text handling.
- [htmlindex](https://pkg.go.dev/golang.org/x/text/encoding/htmlindex) Package htmlindex maps character set encoding names to Encodings as recommended by the W3C for use in HTML 5.
- [ianaindex](https://pkg.go.dev/golang.org/x/text/encoding/ianaindex) Package ianaindex maps names to Encodings as specified by the IANA registry. This includes both the MIME and IANA names.
- [japanese](https://pkg.go.dev/golang.org/x/text/encoding/japanese) Package japanese provides Japanese encodings such as EUC-JP and Shift JIS.
---
- [number](https://pkg.go.dev/golang.org/x/text/number) Package number formats numbers according to the customs of different locales.
- [search](https://pkg.go.dev/golang.org/x/text/search) Package search provides language-specific search and string matching. Natural language matching can be intricate. For example, Danish will insist "Århus" and "Aarhus" are the same name and Turkish will match I to ı (note the lack of a dot) in a case-insensitive match. This package handles such language-specific details.
- [width](https://pkg.go.dev/golang.org/x/text/width) Package width provides functionality for handling different widths in text. Wide characters behave like ideographs; they tend to allow line breaks after each character and remain upright in vertical text layout. Narrow characters are kept together in words or runs that are rotated sideways in vertical text layout.

#### Encoding
- [x-encoding](https://pkg.go.dev/golang.org/x/encoding)
- [x-charmap](https://pkg.go.dev/golang.org/x/encoding/charmap) Package charmap provides simple character encodings such as IBM Code Page 437 and Windows 1252.

#### Language
- [plural](https://pkg.go.dev/golang.org/x/text/feature/plural) Package plural provides utilities for handling linguistic plurals in text.
- [display](https://pkg.go.dev/golang.org/x/text/language/display) Package display provides display names for languages, scripts and regions in a requested language.

##### Sys
- [execabs](https://pkg.go.dev/golang.org/x/sys/execabs) Package execabs is a drop-in replacement for os/exec that requires PATH lookups to find absolute paths.

#### Functions

##### [strconv](https://pkg.go.dev/strconv)

- AppendQuote
- AppendQuoteRune
- AppendQuoteRuneToASCII
- FormatBool
- FormatFloat
- FormatInt/FormatUint
- ParseBool(str) (bool,error)
- ParseFloat(s, bitSize int) (float64, error)
- ParseInt(s, base int, bitSize int) (int64, error)
- Quote(s) string
- QuoteRune(r) string
- QuoteRuneToASCII(r) '\u263a'
- Unquote(s string) (string, error)