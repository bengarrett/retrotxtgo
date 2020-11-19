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

type B struct{}

const (
	_ = 1.0 << (10 * iota) // ignore first value by assigning to blank identifier
	kiB
	miB
	giB
	tiB
)

const (
	byte = 1
	kB   = 1000 * byte
	mB   = kB * kB
	gB   = mB * kB
	tB   = gB * kB
)

// New creates a B instance.
func New() *B {
	return &B{}
}

// Binary formats bytes integer to localized readable string.
func (*B) Binary(b int64, t language.Tag) string {
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
		return p.Sprintf("%.1f %s", value, "KiB")
	case b == 0:
		return "0"
	default:
		return p.Sprintf("%dB", b)
	}
	return p.Sprintf("%.2f %s", value, multiple)
}

// Decimal formats bytes integer to localized readable string.
func (*B) Decimal(b int64, t language.Tag) string {
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
		return p.Sprintf("%.1f %s", value, "KB")
	case b == 0:
		return "0"
	default:
		return p.Sprintf("%dB", b)
	}
	return p.Sprintf("%.2f %s", value, multiple)
}

// Bytes formats bytes integer to localized readable string.
func Bytes(b int64, t language.Tag) string {
	return New().Binary(b, t)
}

// Binary formats bytes integer to localized readable string.
func Binary(b int64, t language.Tag) string {
	return New().Binary(b, t)
}

// Decimal formats bytes integer to localized readable string.
func Decimal(b int64, t language.Tag) string {
	return New().Decimal(b, t)
}
