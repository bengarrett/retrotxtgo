package create

import (
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
