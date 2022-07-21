package port

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/chk"
	"github.com/bengarrett/retrotxtgo/lib/str"
)

var ErrNoChange = errors.New("no changes applied")

const (
	Min uint = 1 // Min 0 is not valid as Viper treats it as a null value
	Max uint = 65535
)

// Port asks for and validates HTTP ports.
func Port(w io.Writer, r io.Reader, validate, setup bool) uint {
	const reset uint = 0
	const baseTen = 10
	var (
		input   string
		prompts int
	)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		input = scanner.Text()
		if setup && input == "" {
			break
		}
		if input == "" {
			if err := chk.Check(w, prompts); err != nil {
				return reset
			}
			continue
		}
		value, err := strconv.ParseInt(input, baseTen, 0)
		if err != nil {
			fmt.Printf("%s %v\n", str.Bool(false), input)
			if err := chk.Check(w, prompts); err != nil {
				return reset
			}
			continue
		}
		p := uint(value)
		if validate {
			if v := Valid(p); !v {
				fmt.Printf("%s %v, is out of range\n", str.Bool(false), input)
				if err := chk.Check(w, prompts); err != nil {
					return reset
				}
				continue
			}
		}
		return p
	}
	return reset
}

// Valid checks if the network port is within range to serve HTTP.
func Valid(p uint) bool {
	if p < Min || p > Max {
		return false
	}
	return true
}
