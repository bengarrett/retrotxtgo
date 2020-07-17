//Package convert is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package convert

import (
	"golang.org/x/text/encoding/unicode"
)

const (
	bom  = unicode.UseBOM
	_bom = unicode.IgnoreBOM
	be   = unicode.BigEndian
	le   = unicode.LittleEndian
)
