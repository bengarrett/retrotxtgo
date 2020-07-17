package prompt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

type keys []string

// the lowest, largest and recommended network ports to serve HTTP.
const (
	PortMin uint = 0
	PortMax uint = 65535
	PortRec uint = 8080
)

// SetupMode changes the behavor of prompts to exit on empty inputs.
var SetupMode = false

// IndexStrings asks for a numeric index position and returns a single choice from a string of keys.
func IndexStrings(options *[]string) (key string) {
	var k keys = *options
	return k.prompt(os.Stdin)
}

// Port asks for and returns a HTTP port value.
func Port(validate bool) (port uint) {
	return pport(os.Stdin, validate)
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
func ShortStrings(options *[]string) (key string) {
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
func Strings(options *[]string) (key string) {
	var k keys = *options
	return k.prompt(os.Stdin)
}

// YesNo asks for a yes or no input.
func YesNo(ask string, yesDefault bool) bool {
	y, n := "Y", "n"
	if !yesDefault {
		y, n = "y", "N"
	}
	fmt.Printf("%s? [%s/%s] ", ask, y, n)
	input, err := promptRead(os.Stdin)
	if err != nil {
		logs.LogFatal(err)
	}
	return parseYN(input, yesDefault)
}

func check(prompts int) {
	switch {
	case prompts == 2:
		if i, err := fmt.Println("Ctrl+C to keep the existing value"); err != nil {
			log.Fatalf("prompt check println at %dB: %s", i, err)
		}
	case prompts >= 4:
		os.Exit(1)
	}
}

func pport(r io.Reader, validate bool) (port uint) {
	input, prompts := "", 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		input = scanner.Text()
		if SetupMode && input == "" {
			break
		}
		if input == "" {
			check(prompts)
			continue
		}
		value, err := strconv.ParseInt(input, 10, 0)
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
	return 0
}

func promptRead(stdin io.Reader) (input string, err error) {
	reader := bufio.NewReader(stdin)
	input, err = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && err != io.EOF {
		return input, fmt.Errorf("prompt string reader error: %s", err)
	}
	return input, nil
}

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

func pstring(r io.Reader) (words string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words := scanner.Text()
		switch words {
		case "-":
			return "-"
		case "":
			if SetupMode {
				return ""
			}
			os.Exit(0)
		default:
			return words
		}
	}
	if err := scanner.Err(); err != nil {
		logs.Fatal("prompt string scanner", "stdin", err)
	}
	return words
}

// numeric checks input string for a valid int value and returns a matching slice.
func (k keys) numeric(input string) (key string) {
	if input == "" {
		return key
	}
	i, err := strconv.Atoi(input)
	if err != nil {
		return key
	}
	if i >= len(k) || i < 0 {
		return key
	}
	sort.Strings(k)
	return k[i]
}

func (k keys) prompt(r io.Reader) (key string) {
	prompts := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key = scanner.Text()
		if key == "-" {
			fmt.Println("OKAY")
			return ""
		}
		if SetupMode && key == "" {
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

func (k keys) shortPrompt(r io.Reader) (key string) {
	prompts := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key = scanner.Text()
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

// shortValidate validates the key exists in keys.
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
