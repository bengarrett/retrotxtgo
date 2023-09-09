package sgr

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/bg"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/fg"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/xterm"
)

// Attributes are values for applying different CSS styles
// and colors to the Bytes of text.
// TODO: consolidate these bools? Change Bold/Faint bool to Intensity uint (normal 0, bold 1, faint 2)?
// TODO: then change SetBold() SetFaint() to SetIntensity().
type Attributes struct {
	Bytes         []byte      // Bytes of text for display.
	Foreground    xterm.Color // Text color 8 bit value.
	ForegroundRGB string      // Text override color as a RGB value.
	Background    xterm.Color // Background color 8 bit value.
	BackgroundRGB string      // Background override color as a RGB value.
	Font          Fonts       // Font selection.
	Bold          bool        // Increased font intensity.
	Faint         bool        // Decreased font intensity.
	Italic        bool        // Use an italic type font.
	Underline     bool        // Apply an underline text decoration.
	Blink         bool        // Animate the text with a standard speed blink.
	BlinkFast     bool        // Animate the text with a rapid blink.
	Inverse       bool        // Invert the color of the text.
	Conceal       bool        // Sets the text color to match the background color.
	StrikeThrough bool        // Apply a horizontal line decoration through the middle of the text.
	Underline2x   bool        // Apply a doubled underline text decoration.
	Framed        bool        // Framed text decoration.
	Encircled     bool        // Encircled text decoration.
	Overlined     bool        // Overlined text decoration.
}

// String returns the Attributes as CSS class names.
// TODO: background, foreground ... BOOL FONTS RGB colors.
// TODO: also return (css, style string).
func (a Attributes) String() string {
	var cls []string
	// Conceal hides the text so do not return any other styles.
	// TODO: a conceal func to make foreground with background??
	if a.Conceal {
		cls = append(cls, a.Foreground.String(), a.Background.String(), Conceal.String())
		return strings.Join(cls, " ")
	}
	// Colors
	// TODO: if a.Bytes is empty or only whitespace,
	// set Foreground to match Background to reduce artifacts.
	cls = append(cls, a.Foreground.String(), a.Background.String())
	// Bold or faint intensity
	if a.Bold {
		cls = append(cls, Bold.String())
	} else if a.Faint {
		cls = append(cls, Faint.String())
	}
	// Italic type
	if a.Italic {
		cls = append(cls, Italic.String())
	}
	// Underline style
	if a.Underline2x {
		cls = append(cls, Underline2x.String())
	} else if a.Underline {
		cls = append(cls, Underline.String())
	}
	// Animation
	if a.Blink {
		cls = append(cls, Blink.String())
	} else if a.BlinkFast {
		cls = append(cls, BlinkFast.String())
	}
	// Reverse color
	if a.Inverse {
		cls = append(cls, Inverse.String())
	}
	// Line strike-through style
	if a.StrikeThrough {
		cls = append(cls, StrikeThrough.String())
	}
	// Encircled, framed or overlined style
	if a.Encircled {
		cls = append(cls, Encircled.String())
	} else if a.Framed {
		cls = append(cls, Framed.String())
	} else if a.Overlined {
		cls = append(cls, Overlined.String())
	}
	return strings.Join(cls, " ")
}

func (a *Attributes) DataStream(b []byte) {
	if string(b) == "" {
		return
	}
	i := bytes.IndexByte(b, 'm')
	if i == -1 {
		a.Bytes = b
		return
	}

	v := b[0:i]
	bs := bytes.Split(v, []byte(";"))
	if len(bs) == 0 {
		a.Bytes = b
		return
	}

	var ext Extensions
	ext.Reset()
	for _, ps := range bs {
		rs := bytes.Runes(ps)
		for _, r := range rs {
			if !unicode.IsDigit(r) {
				a.Bytes = b
				return
			}
		}
		s, err := strconv.Atoi(string(ps))
		if err != nil {
			a.Bytes = b
			return
		}
		if !Ps(s).Valid() {
			a.Bytes = b
			return
		}
		if cont := ext.Scan(s); cont {
			// TODO: background and foreground.
			if ext.Color > -1 {
				// todo bg/fg
				a.Background = xterm.Color(ext.Color)
			}
			continue
		}
		a.Set(Ps(s))
	}
	if string(b[i+1:]) != "" {
		a.Bytes = b[i+1:]
	}
}

func (a *Attributes) Set(ps Ps) {
	a.SetForeground(ps)
	a.SetBackground(ps)
	a.SetFont(ps)
	a.SetBold(ps)
	a.SetFaint(ps)
	a.SetItalic(ps)
	a.SetUnderline(ps)
	a.SetBlink(ps)
	a.SetBlinkFast(ps)
	a.SetInverse(ps)
	a.SetConceal(ps)
	a.SetStrike(ps)
	a.SetUnderline2x(ps)
	a.SetFramed(ps)
	a.SetEncircled(ps)
	a.SetOverlined(ps)
}

func (a *Attributes) SetForeground(ps Ps) {
	switch {
	case ps == Normal:
		a.Foreground = xterm.Foreground(fg.White)
	case
		ps >= Black && ps <= White,
		ps >= BoldBlack && ps <= BoldWhite:
		a.Foreground = xterm.Foreground(fg.Colors(ps))
	}
}

func (a *Attributes) SetBackground(ps Ps) {
	switch {
	case ps == Normal:
		a.Background = xterm.Background(bg.Black)
	case
		ps >= BlackB && ps <= WhiteB,
		ps >= BrightBlack && ps <= BrightWhite:
		a.Background = xterm.Background(bg.Colors(ps))
	}
}

func (a *Attributes) SetFont(ps Ps) {
	switch ps {
	case Normal, NotItalicFraktur:
		a.Font = Primary
	case
		Font0, Font1, Font2, Font3, Font4, Font5,
		Font6, Font7, Font8, Font9, Fraktur:
		a.Font = Fonts(ps)
	}
}

func (a *Attributes) SetBold(ps Ps) {
	switch ps {
	case Normal, NotBoldFaint:
		a.Bold = false
	case Bold:
		a.Bold = true
	}
}

func (a *Attributes) SetFaint(ps Ps) {
	switch ps {
	case Normal, NotBoldFaint:
		a.Faint = false
	case Faint:
		a.Faint = true
	}
}

func (a *Attributes) SetItalic(ps Ps) {
	switch ps {
	case Normal, NotItalicFraktur:
		a.Italic = false
	case Italic:
		a.Italic = true
	}
}

func (a *Attributes) SetUnderline(ps Ps) {
	switch ps {
	case Normal, NotUnderline:
		a.Underline = false
	case Underline:
		a.Underline = true
	}
}

func (a *Attributes) SetBlink(ps Ps) {
	switch ps {
	case Normal:
		a.Blink = false
	case Blink:
		a.Blink = true
	}
}

func (a *Attributes) SetBlinkFast(ps Ps) {
	switch ps {
	case Normal:
		a.BlinkFast = false
	case BlinkFast:
		a.BlinkFast = true
	}
}

func (a *Attributes) SetInverse(ps Ps) {
	switch ps {
	case Normal:
		a.Inverse = false
	case Inverse:
		a.Inverse = true
	}
}

func (a *Attributes) SetConceal(ps Ps) {
	switch ps {
	case Normal, RevertB:
		a.Conceal = false
	case Conceal:
		a.Conceal = true
	}
}

func (a *Attributes) SetStrike(ps Ps) {
	switch ps {
	case Normal, NotStrikeThrough:
		a.StrikeThrough = false
	case StrikeThrough:
		a.StrikeThrough = true
	}
}

func (a *Attributes) SetUnderline2x(ps Ps) {
	switch ps {
	case Normal, NotUnderline:
		a.Underline2x = false
	case Underline2x:
		a.Underline2x = true
	}
}

func (a *Attributes) SetFramed(ps Ps) {
	switch ps {
	case Normal:
		a.Framed = false
	case Framed:
		a.Framed = true
	}
}

func (a *Attributes) SetEncircled(ps Ps) {
	switch ps {
	case Normal:
		a.Encircled = false
	case Framed:
		a.Encircled = true
	}
}

func (a *Attributes) SetOverlined(ps Ps) {
	switch ps {
	case Normal:
		a.Overlined = false
	case Framed:
		a.Overlined = true
	}
}

func DataStream(b []byte) Attributes {
	if string(b) == "" {
		return Attributes{}
	}
	a := Attributes{}
	i := bytes.IndexByte(b, 'm')
	if i == -1 {
		return Attributes{Bytes: b}
	}
	v := b[0:i]
	bs := bytes.Split(v, []byte(";"))
	if len(bs) == 0 {
		return Attributes{Bytes: b}
	}

	var ext Extensions
	ext.Reset()
	for _, ps := range bs {
		rs := bytes.Runes(ps)
		for _, r := range rs {
			if !unicode.IsDigit(r) {
				return Attributes{Bytes: b}
			}
		}
		s, err := strconv.Atoi(string(ps))
		if err != nil {
			return Attributes{Bytes: b}
		}
		if !Ps(s).Valid() {
			return Attributes{Bytes: b}
		}
		if cont := ext.Scan(s); cont {
			if ext.Color > -1 {
				// todo bg/fg
				a.Background = xterm.Color(ext.Color)
			}
			continue
		}
		fmt.Println(ps)
		// if color, l := Colors(Ps(s), bs[i:]); color > -1 && l > 0 {
		// 	//cnt = cnt + l
		// 	continue
		// 	//vals[i] = Ps(color)
		// 	//break // assume no other styling if special colors are used
		// }
		// vals[i] = Ps(s)

		a.Set(Ps(s))
	}
	if string(b[i+1:]) != "" {
		a.Bytes = b[i+1:]
	}
	return a
}

// NewSGR creates a new Attributes struct with default values.
func NewSGR() Attributes {
	return Attributes{
		Background: xterm.Color(bg.System),
		Foreground: xterm.Color(fg.System),
		Font:       Primary,
	}
}
