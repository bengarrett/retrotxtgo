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
// TODO: then change SetBold() SetFaint() to SetIntensity()
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
// TODO: background, foreground ... BOOL FONTS RGB colors
// TODO: also return (css, style string)
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

func (t *Attributes) DataStream(b []byte) {
	if string(b) == "" {
		return
	}
	i := bytes.IndexByte(b, 'm')
	if i == -1 {
		t.Bytes = b
		return
	}

	v := b[0:i]
	bs := bytes.Split(v, []byte(";"))
	if len(bs) == 0 {
		t.Bytes = b
		return
	}

	var ext Extensions
	ext.Reset()
	for _, ps := range bs {
		rs := bytes.Runes(ps)
		for _, r := range rs {
			if !unicode.IsDigit(r) {
				t.Bytes = b
				return
			}
		}
		s, err := strconv.Atoi(string(ps))
		if err != nil {
			t.Bytes = b
			return
		}
		if !Ps(s).Valid() {
			t.Bytes = b
			return
		}
		if cont := ext.Scan(s); cont {
			// TODO: background and foreground
			if ext.Color > -1 {
				// todo bg/fg
				t.Background = xterm.Color(ext.Color)
			}
			continue
		}
		t.Set(Ps(s))
	}
	if string(b[i+1:]) != "" {
		t.Bytes = b[i+1:]
	}
}

func (t *Attributes) Set(ps Ps) {
	t.SetForeground(ps)
	t.SetBackground(ps)
	t.SetFont(ps)
	t.SetBold(ps)
	t.SetFaint(ps)
	t.SetItalic(ps)
	t.SetUnderline(ps)
	t.SetBlink(ps)
	t.SetBlinkFast(ps)
	t.SetInverse(ps)
	t.SetConceal(ps)
	t.SetStrike(ps)
	t.SetUnderline2x(ps)
	t.SetFramed(ps)
	t.SetEncircled(ps)
	t.SetOverlined(ps)
}

func (t *Attributes) SetForeground(ps Ps) {
	switch {
	case ps == Normal:
		t.Foreground = xterm.Foreground(fg.White)
	case
		ps >= Black && ps <= White,
		ps >= BoldBlack && ps <= BoldWhite:
		t.Foreground = xterm.Foreground(fg.Colors(ps))
	}
}

func (t *Attributes) SetBackground(ps Ps) {
	switch {
	case ps == Normal:
		t.Background = xterm.Background(bg.Black)
	case
		ps >= BlackB && ps <= WhiteB,
		ps >= BrightBlack && ps <= BrightWhite:
		t.Background = xterm.Background(bg.Colors(ps))
	}
}

func (t *Attributes) SetFont(ps Ps) {
	switch ps {
	case Normal, NotItalicFraktur:
		t.Font = Primary
	case
		Font0, Font1, Font2, Font3, Font4, Font5,
		Font6, Font7, Font8, Font9, Fraktur:
		t.Font = Fonts(ps)
	}
}

func (t *Attributes) SetBold(ps Ps) {
	switch ps {
	case Normal, NotBoldFaint:
		t.Bold = false
	case Bold:
		t.Bold = true
	}
}

func (t *Attributes) SetFaint(ps Ps) {
	switch ps {
	case Normal, NotBoldFaint:
		t.Faint = false
	case Faint:
		t.Faint = true
	}
}

func (t *Attributes) SetItalic(ps Ps) {
	switch ps {
	case Normal, NotItalicFraktur:
		t.Italic = false
	case Italic:
		t.Italic = true
	}
}

func (t *Attributes) SetUnderline(ps Ps) {
	switch ps {
	case Normal, NotUnderline:
		t.Underline = false
	case Underline:
		t.Underline = true
	}
}

func (t *Attributes) SetBlink(ps Ps) {
	switch ps {
	case Normal:
		t.Blink = false
	case Blink:
		t.Blink = true
	}
}

func (t *Attributes) SetBlinkFast(ps Ps) {
	switch ps {
	case Normal:
		t.BlinkFast = false
	case BlinkFast:
		t.BlinkFast = true
	}
}

func (t *Attributes) SetInverse(ps Ps) {
	switch ps {
	case Normal:
		t.Inverse = false
	case Inverse:
		t.Inverse = true
	}
}

func (t *Attributes) SetConceal(ps Ps) {
	switch ps {
	case Normal, RevertB:
		t.Conceal = false
	case Conceal:
		t.Conceal = true
	}
}

func (t *Attributes) SetStrike(ps Ps) {
	switch ps {
	case Normal, NotStrikeThrough:
		t.StrikeThrough = false
	case StrikeThrough:
		t.StrikeThrough = true
	}
}

func (t *Attributes) SetUnderline2x(ps Ps) {
	switch ps {
	case Normal, NotUnderline:
		t.Underline2x = false
	case Underline2x:
		t.Underline2x = true
	}
}

func (t *Attributes) SetFramed(ps Ps) {
	switch ps {
	case Normal:
		t.Framed = false
	case Framed:
		t.Framed = true
	}
}

func (t *Attributes) SetEncircled(ps Ps) {
	switch ps {
	case Normal:
		t.Encircled = false
	case Framed:
		t.Encircled = true
	}
}

func (t *Attributes) SetOverlined(ps Ps) {
	switch ps {
	case Normal:
		t.Overlined = false
	case Framed:
		t.Overlined = true
	}
}

func DataStream(b []byte) Attributes {
	if string(b) == "" {
		return Attributes{}
	}
	var t Attributes
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
				t.Background = xterm.Color(ext.Color)
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

		t.Set(Ps(s))
	}
	if string(b[i+1:]) != "" {
		t.Bytes = b[i+1:]
	}
	return t
}

// NewSGR creates a new Attributes struct with default values.
func NewSGR() Attributes {
	return Attributes{
		Background: xterm.Color(bg.System),
		Foreground: xterm.Color(fg.System),
		Font:       Primary,
	}
}
