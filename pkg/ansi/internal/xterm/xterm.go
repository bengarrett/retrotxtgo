package xterm

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/bg"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/fg"
)

const (
	Control = 5
)

type Color uint

func (c Color) String() string {
	return fmt.Sprintf("SGR%d", c)
}

func Foreground(c fg.Colors) Color {
	if c >= fg.IbmAIX {
		const base = fg.White - fg.System
		return Color(c - fg.IbmAIX + base)
	}
	return Color(c)
	// return Color(fg.System + (c - fg.System))
}

func Background(c bg.Colors) Color {
	if c >= bg.IbmAIX {
		const base = bg.White - bg.System
		return Color(c - bg.IbmAIX + base)
	}
	return Color(c)
	// return Color(bg.System + (c - bg.System))
}
