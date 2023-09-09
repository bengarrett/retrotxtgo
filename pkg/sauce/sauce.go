package sauce

import (
	"github.com/bengarrett/sauce"
)

type SAUCE struct {
	Use         bool
	Title       string
	Author      string
	Group       string
	Description string
	Width       uint
	Lines       uint
}

// Read returns any SAUCE metadata that is attached to the byte array.
func (s *SAUCE) Read(b *[]byte) {
	sr := sauce.Decode(*b)
	if !sr.Valid() {
		return
	}
	s = &SAUCE{
		Use:         true,
		Title:       sr.Title,
		Author:      sr.Author,
		Group:       sr.Group,
		Description: sr.Desc,
		Width:       uint(sr.Info.Info1.Value),
		Lines:       uint(sr.Info.Info2.Value),
	}
}
