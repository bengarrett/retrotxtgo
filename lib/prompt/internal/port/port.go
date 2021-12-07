package port

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/bengarrett/retrotxtgo/lib/str"
)

const (
	Min      uint = 1 // Min 0 is not valid as Viper treats it as a null value
	Max      uint = 65535
	NoChange      = "no changes applied"
)

// Check used in scanner Scans to prompt a user for a valid text input.
func Check(prompts int) {
	const info, max = 2, 4
	switch {
	case prompts == info:
		if i, err := fmt.Print("\r  Ctrl+C to keep the existing value\n"); err != nil {
			log.Fatalf("prompt check println at %dB: %s", i, err)
		}
	case prompts >= max:
		fmt.Printf("\r  %s", NoChange)
		os.Exit(0)
	}
}

// Port asks for and validates HTTP ports.
func Port(r io.Reader, validate, setup bool) uint {
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
			Check(prompts)
			continue
		}
		value, err := strconv.ParseInt(input, baseTen, 0)
		if err != nil {
			fmt.Printf("%s %v\n", str.Bool(false), input)
			Check(prompts)
			continue
		}
		p := uint(value)
		if validate {
			if v := Valid(p); !v {
				fmt.Printf("%s %v, is out of range\n", str.Bool(false), input)
				Check(prompts)
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
