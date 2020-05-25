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

type (
	// B struct
	B struct{}
)

const (
	_ = 1.0 << (10 * iota) // ignore first value by assigning to blank identifier
	KB
	MB
	GB
	TB
)

var global = New()

// New creates a B instance.
func New() *B {
	return &B{}
}

// Bytes formats bytes integer to localized readable string.
// For example, 31323 bytes will return 30.59KB with language.AmericanEnglish.
func (*B) Bytes(b int64, t language.Tag) string {
	p := message.NewPrinter(t)

	multiple := ""
	value := float64(b)

	switch {
	case b >= TB:
		value /= TB
		multiple = "TB"
	case b >= GB:
		value /= GB
		multiple = "GB"
	case b >= MB:
		value /= MB
		multiple = "MB"
	case b >= KB:
		value /= KB
		multiple = "KB"
	case b == 0:
		return "0"
	default:
		return p.Sprintf("%dB", b)
	}
	return p.Sprintf("%.2f%s", value, multiple)
}

// Bytes formats bytes integer to localized readable string.
// For example, 31323 bytes will return 30.59KB with language.AmericanEnglish.
func Bytes(b int64, t language.Tag) string {
	return global.Bytes(b, t)
}
