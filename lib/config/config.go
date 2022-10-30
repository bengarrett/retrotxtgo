// Package config handles the user configations.
package config

import (
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
)

// Formats the choices for command flags.
type Formats struct {
	Info [5]string
}

// Format flag choices for the info command.
func Format() Formats {
	return Formats{
		Info: [5]string{"color", "json", "json.min", "text", "xml"},
	}
}

// Tip provides some brief help information of the configurations.
func Tip() get.Hints {
	return get.Tip()
}
