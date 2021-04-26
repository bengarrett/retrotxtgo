package prompt

import "errors"

var (
	ErrNoReader = errors.New("reader interface is empty")
	ErrPString  = errors.New("prompt string standard input problem")
)
