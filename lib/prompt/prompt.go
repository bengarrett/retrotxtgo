// Package prompt inputs for user interactions.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

type keys []string

// the lowest, largest and recommended network ports to serve HTTP.
const (
	PortMin  uint = 1 // port 0 is not valid as viper treats it as a null value
	PortMax  uint = 65535
	NoChange      = "no changes applied"
)

// CtrlC intercepts any Ctrl-C keyboard input and exits to the shell.
func CtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(os.Stdout, "%s has quit.", meta.Name)
		os.Exit(0)
	}()
}

// IndexStrings asks for a numeric index position and returns a single choice from a string of keys.
func IndexStrings(options *[]string, setup bool) string {
	if options == nil {
		return ""
	}
	var k keys = *options
	return k.prompt(os.Stdin, setup)
}

// SkipSet returns a skipped setting string for the setup walk through.
func SkipSet(setup bool) string {
	if !setup {
		return ""
	}
	return str.Cs("  skipped setting")
}

// Port asks for and returns a HTTP port value.
func Port(validate, setup bool) (port uint) {
	return pport(os.Stdin, validate, setup)
}

// PortValid checks if the network port is within range to serve HTTP.
func PortValid(port uint) (ok bool) {
	if port < PortMin || port > PortMax {
		return false
	}
	return true
}

// ShortStrings asks for and returns a single choice from a string of keys.
// Either the first letter or the full name of the key are accepted.
func ShortStrings(options *[]string) string {
	if options == nil {
		return ""
	}
	var k keys = *options
	return k.shortPrompt(os.Stdin)
}

// String asks for and returns a multi-word string.
// Inputting ⏎ the Enter/Return key exits the program,
// or returns an empty value when in SetupMode.
func String() (words string) {
	return pstring(os.Stdin)
}

// Strings asks for and returns a single choice from a string of keys.
func Strings(options *[]string, setup bool) string {
	if options == nil {
		return ""
	}
	var k keys = *options
	return k.prompt(os.Stdin, setup)
}

// YesNo asks for a yes or no input.
func YesNo(ask string, yesDefault bool) bool {
	y, n := "Y", "n"
	if !yesDefault {
		y, n = "y", "N"
	}
	if !strings.HasSuffix(ask, "?") && !strings.HasSuffix(ask, ")") {
		ask = fmt.Sprintf("%s?", ask)
	}
	yes, no := y, n
	if color.Enable {
		yes, _ = str.UnderlineChar(y)
		no, _ = str.UnderlineChar(n)
	}
	fmt.Printf("%s [%ses/%so] ", ask, yes, no)
	input, err := promptRead(os.Stdin)
	if err != nil {
		logs.SaveFatal(err)
	}
	return parseYN(input, yesDefault)
}

// Check used in scanner Scans to prompt a user for a valid text input.
func check(prompts int) {
	const info, max = 2, 4
	switch {
	case prompts == info:
		if i, err := fmt.Println("Ctrl+C to keep the existing value"); err != nil {
			log.Fatalf("prompt check println at %dB: %s", i, err)
		}
	case prompts >= max:
		fmt.Println(NoChange)
		os.Exit(0)
	}
}

// PPort asks for and validates HTTP ports.
func pport(r io.Reader, validate, setup bool) (port uint) {
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
			check(prompts)
			continue
		}
		value, err := strconv.ParseInt(input, baseTen, 0)
		if err != nil {
			fmt.Printf("%s %v\n", str.Bool(false), input)
			check(prompts)
			continue
		}
		port = uint(value)
		if validate {
			if v := PortValid(port); !v {
				fmt.Printf("%s %v, is out of range\n", str.Bool(false), input)
				check(prompts)
				continue
			}
		}
		return port
	}
	return reset
}

// PromptRead parses a line of text from the reader.
func promptRead(r io.Reader) (input string, err error) {
	reader := bufio.NewReader(r)
	input, err = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && err != io.EOF {
		return input, fmt.Errorf("prompt string reader error: %w", err)
	}
	return input, nil
}

// ParseYN parses the input to boolean value.
func parseYN(input string, yesDefault bool) bool {
	switch input {
	case "":
		if yesDefault {
			return true
		}
	case "yes", "y":
		return true
	}
	return false
}

// pstring parses the reader input for any os exit commands.
func pstring(r io.Reader) (words string) {
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

// Numeric checks input string for a valid int value and returns a matching slice.
func (k keys) numeric(input string) string {
	if input == "" {
		return ""
	}
	i, err := strconv.Atoi(input)
	if err != nil {
		return ""
	}
	if i >= len(k) || i < 0 {
		return ""
	}
	sort.Strings(k)
	return k[i]
}

// Prompt parses the reader input for a valid key.
func (k keys) prompt(r io.Reader, setup bool) string {
	prompts := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key := scanner.Text()
		if key == "-" {
			return ""
		}
		if setup && key == "" {
			return ""
		}
		if n := k.numeric(key); n != "" {
			return n
		}
		if !k.validate(key) {
			check(prompts)
			continue
		}
		return key
	}
	return ""
}

// ShortPrompt parses the reader input for a valid key or alias of the key.
func (k keys) shortPrompt(r io.Reader) string {
	prompts, scanner := 0, bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key := scanner.Text()
		switch key {
		case "", "-":
			return key
		}
		if long := k.shortValidate(key); long != "" {
			return long
		}
		if !k.validate(key) {
			check(prompts)
			continue
		}
		return key
	}
	return ""
}

// ShortValidate validates the key exists in keys.
// Both the first letter of the key and the full name of the key are accepted as valid.
// Whenever the key is valid the full key name will be returned otherwise an empty string
// signifies a false result.
func (k keys) shortValidate(key string) string {
	if key == "" {
		return ""
	}
	l, x := len(k), strings.ToLower(key)
	sort.Strings(k)
	var a, b = make([]string, l), make([]string, l)
	for i, l := range k {
		a[i] = strings.ToLower(string(l[0]))
		b[i] = strings.ToLower(l)
	}
	// match the first letter
	sort.Strings(a)
	var i = sort.SearchStrings(a, x)
	if i >= len(a) || a[i] != x {
		// match the whole word
		sort.Strings(b)
		var j = sort.SearchStrings(b, x)
		if j >= len(b) || b[j] != x {
			return ""
		}
		return k[j]
	}
	return k[i]
}

// Validate that the key exists in the slice of keys.
func (k keys) validate(key string) (ok bool) {
	if key == "" {
		return false
	}
	sort.Strings(k)
	var i = sort.SearchStrings(k, key)
	if i >= len(k) || k[i] != key {
		fmt.Printf("%s %v\n", str.Bool(false), key)
		return false
	}
	return true
}
