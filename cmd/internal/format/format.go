// Package format handles the text output, syntax options.
package format

// Syntax choices for the input format flag.
type Syntax struct {
	Info [5]string
}

// Format flag choices for the info command.
func Format() Syntax {
	return Syntax{
		Info: [5]string{"color", "json", "json.min", "text", "xml"},
	}
}
