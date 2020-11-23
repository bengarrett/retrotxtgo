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
	oneDecimalPoint  = "%.1f %s"
	twoDecimalPoints = "%.2f %s"
	binaryBase       = 10
	kB               = 1000
	mB               = kB * kB
	gB               = mB * kB
	tB               = gB * kB
)
const (
	_ = 1.0 << (binaryBase * iota) // ignore first value by assigning to blank identifier
	kiB
	miB
	giB
	tiB
)

// Binary formats bytes integer to localized readable string.
func binary(b int64, t language.Tag) string {
	p := message.NewPrinter(t)
	multiple, value := "", float64(b)
	switch {
	case b >= tiB:
		value /= tiB
		multiple = "TiB"
	case b >= giB:
		value /= giB
		multiple = "GiB"
	case b >= miB:
		value /= miB
		multiple = "MiB"
	case b >= kiB:
		value /= kiB
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
	case b >= tB:
		value /= tB
		multiple = "TB"
	case b >= gB:
		value /= gB
		multiple = "GB"
	case b >= mB:
		value /= mB
		multiple = "MB"
	case b >= kB:
		value /= kB
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
