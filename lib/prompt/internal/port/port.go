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
	PortMin  uint = 1 // port 0 is not valid as viper treats it as a null value
	PortMax  uint = 65535
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

// PPort asks for and validates HTTP ports.
func PPort(r io.Reader, validate, setup bool) uint {
	const reset uint = 0
	const baseTen = 10
	input, prompts := "", 0
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
			if v := PortValid(p); !v {
				fmt.Printf("%s %v, is out of range\n", str.Bool(false), input)
				Check(prompts)
				continue
			}
		}
		return p
	}
	return reset
}

// PortValid checks if the network port is within range to serve HTTP.
func PortValid(p uint) bool {
	if p < PortMin || p > PortMax {
		return false
	}
	return true
}
