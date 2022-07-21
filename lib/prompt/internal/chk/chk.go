package chk

import (
	"errors"
	"fmt"
	"io"
)

var ErrNoChange = errors.New("no changes applied")

const (
	alertLoop = 1
	maxLoops  = 2
)

// Check used in scanner Scans to prompt for valid stdin.
func Check(w io.Writer, prompts int) error {
	switch {
	case prompts == alertLoop:
		fmt.Fprint(w, "\r  Press enter to keep the existing value\n")
	case prompts >= maxLoops:
		return ErrNoChange
	}
	return nil
}
