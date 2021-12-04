package read

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
)

var (
	ErrNoReader = errors.New("reader interface is empty")
	ErrPString  = errors.New("prompt string standard input problem")
)

// Read parses a line of text from the reader.
func Read(r io.Reader) (input string, err error) {
	reader := bufio.NewReader(r)
	input, err = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && err != io.EOF {
		return input, fmt.Errorf("prompt string reader error: %w", err)
	}
	return input, nil
}

// ParseYN parses the input to boolean value.
func ParseYN(input string, yesDefault bool) bool {
	switch strings.ToLower(input) {
	case "":
		if yesDefault {
			return true
		}
	case "yes", "y":
		return true
	}
	return false
}

// Parse parses the reader input for any os exit commands.
func Parse(r io.Reader) (words string) {
	if r == nil {
		logs.FatalWrap(ErrPString, ErrNoReader)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		logs.FatalWrap(ErrPString, err)
	}
	return words
}
