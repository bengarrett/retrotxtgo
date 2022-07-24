package read

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrNoChange = errors.New("no changes applied")
	ErrNoReader = errors.New("reader interface is empty")
	ErrPString  = errors.New("prompt string standard input problem")
)

const (
	alertLoop = 1
	maxLoops  = 2
)

// Check used by the scanner iterations to prompt for valid stdin.
func Check(w io.Writer, prompts int) error {
	switch {
	case prompts == alertLoop:
		fmt.Fprint(w, "\r  Press enter to keep the existing value\n")
	case prompts >= maxLoops:
		return ErrNoChange
	}
	return nil
}

// Read parses a line of text from the reader.
func Read(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && !errors.Is(err, io.EOF) {
		return input, fmt.Errorf("prompt string reader error: %w", err)
	}
	return input, nil
}

// ParseYN parses the input to boolean value.
func ParseYN(input string, suggestion bool) bool {
	switch strings.ToLower(input) {
	case "":
		if suggestion {
			return true
		}
	case "yes", "y":
		return true
	}
	return false
}

// Parse parses the reader input for any os exit commands.
func Parse(r io.Reader) (string, error) {
	if r == nil {
		return "", fmt.Errorf("%w: %s", ErrPString, ErrNoReader)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("%w: %s", ErrPString, err)
	}
	return "", nil
}
