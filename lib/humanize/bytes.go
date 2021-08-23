// Package humanize parses data to a human readable format.
package humanize

// The code on this page is derived from labstack/gommon, Common packages for Go
// https://github.com/labstack/gommon.
//
// The MIT License (MIT) Copyright (c) 2018 labstack
// https://github.com/labstack/gommon/blob/master/LICENSE

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	_ = 1.0 << (binaryBase * iota)
	kib
	mib
	gib
	tib

	oneDecimalPoint  = "%.1f %s"
	twoDecimalPoints = "%.2f %s"
	binaryBase       = 10
	kb               = 1000
	mb               = kb * kb
	gb               = mb * kb
	tb               = gb * kb
)

// Binary formats bytes integer to localized readable string.
func binary(b int64, t language.Tag) string {
	p := message.NewPrinter(t)
	multiple, value := "", float64(b)
	switch {
	case b >= tib:
		value /= tib
		multiple = "TiB"
	case b >= gib:
		value /= gib
		multiple = "GiB"
	case b >= mib:
		value /= mib
		multiple = "MiB"
	case b >= kib:
		value /= kib
		return p.Sprintf(oneDecimalPoint, value, "KiB")
	case b == 0:
		return "0"
	default:
		return p.Sprintf("%dB", b)
	}
	return p.Sprintf(twoDecimalPoints, value, multiple)
}

// Decimal formats bytes integer to localized readable string.
func decimal(b int64, t language.Tag) string {
	p := message.NewPrinter(t)
	multiple, value := "", float64(b)
	switch {
	case b >= tb:
		value /= tb
		multiple = "TB"
	case b >= gb:
		value /= gb
		multiple = "GB"
	case b >= mb:
		value /= mb
		multiple = "MB"
	case b >= kb:
		value /= kb
		return p.Sprintf(oneDecimalPoint, value, "kB")
	case b == 0:
		return "0"
	default:
		return p.Sprintf("%dB", b)
	}
	return p.Sprintf(twoDecimalPoints, value, multiple)
}

// Binary formats bytes integer to localized readable string.
func Binary(b int64, t language.Tag) string {
	return binary(b, t)
}

// Decimal formats bytes integer to localized readable string.
func Decimal(b int64, t language.Tag) string {
	return decimal(b, t)
}
