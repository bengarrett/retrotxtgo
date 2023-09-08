package convert

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// Encoding is an implementation of the Encoding interface that adds the String
// and ID methods to an existing encoding.
type Encoding struct {
	encoding.Encoding
	Name string
}

var (
	// AsaX34_1963 ASA X3.4 1963.
	AsaX34_1963 encoding.Encoding = &x34_1963 // nolint: gochecknoglobals

	// AsaX34_1965 ASA X3.4 1965.
	AsaX34_1965 encoding.Encoding = &x34_1965 // nolint: gochecknoglobals

	// AnsiX34_1967 ANSI X3.4 1967/77/86.
	AnsiX34_1967 encoding.Encoding = &x34_1967 // nolint: gochecknoglobals

	x34_1963 = Encoding{ // nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1963",
	}
	x34_1965 = Encoding{ // nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1965",
	}
	x34_1967 = Encoding{ // nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ANSI X3.4 1967/77/86",
	}
)

func (e Encoding) String() string {
	return e.Name
}

// AsaX34 returns a named value for the legacy ASA ASCII character encodings.
func AsaX34(e encoding.Encoding) string {
	switch e {
	case AsaX34_1963:
		return ascii63
	case AsaX34_1965:
		return ascii65
	case AnsiX34_1967:
		return ascii67
	}
	return ""
}
