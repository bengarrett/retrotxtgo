package logs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func promptCheck(prompts int) {
	switch {
	case prompts == 2:
		if i, err := fmt.Println("Ctrl+C to keep the existing value"); err != nil {
			log.Fatalf("logs.promptCheck println at %dB: %s", i, err)
		}
	case prompts >= 4:
		os.Exit(1)
	}
}

// PromptPort asks for and returns a HTTP port value.
func PromptPort(validate bool) (port uint) {
	return pport(os.Stdin, validate)
}

func pport(r io.Reader, validate bool) (port uint) {
	input, prompts := "", 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		input = scanner.Text()
		if input == "" {
			promptCheck(prompts)
			continue
		}
		value, err := strconv.ParseInt(input, 10, 0)
		if err != nil {
			fmt.Printf("%s %v\n", Ce("✗"), input)
			promptCheck(prompts)
			continue
		}
		port = uint(value)
		if validate {
			if v := PortValid(port); !v {
				fmt.Printf("%s %v, is out of range\n", Ce("✗"), input)
				promptCheck(prompts)
				continue
			}
		}
		return port
	}
	return 0
}

// PortValid checks if the network port is within range to serve HTTP.
func PortValid(port uint) (ok bool) {
	if port < PortMin || port > PortMax {
		return false
	}
	return true
}

// PromptString asks for and returns a multi-word string.
func PromptString() (words string) {
	return pstring(os.Stdin)
}

func pstring(r io.Reader) (words string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words := scanner.Text()
		switch words {
		case "":
			os.Exit(0)
		case "-":
			words = ""
		}
		return words
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading standard input: %s", err)
	}
	return words
}

type keys []string

// PromptStrings asks for and returns a single choice from a string of keys.
func PromptStrings(options *[]string) (key string) {
	var k keys = *options
	return k.prompt(os.Stdin)
}

func (k keys) validate(key string) (ok bool) {
	if key == "" {
		return false
	}
	sort.Strings(k)
	var i = sort.SearchStrings(k, key)
	if i >= len(k) || k[i] != key {
		fmt.Printf("%s %v\n", Ce("✗"), key)
		return false
	}
	return true
}

func (k keys) prompt(r io.Reader) (key string) {
	prompts := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key = scanner.Text()
		if !k.validate(key) {
			promptCheck(prompts)
			continue
		}
		return key
	}
	return ""
}

// PromptYN asks for a yes or no input.
func PromptYN(ask string, yesDefault bool) bool {
	y, n := "Y", "n"
	if !yesDefault {
		y, n = "y", "N"
	}
	fmt.Printf("%s? [%s/%s] ", ask, y, n)
	input, err := promptRead(os.Stdin)
	Log(err)
	return parseyn(input, yesDefault)
}

func promptRead(stdin io.Reader) (input string, err error) {
	reader := bufio.NewReader(stdin)
	input, err = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && err != io.EOF {
		return input, err
	}
	return input, nil
}

func parseyn(input string, yesDefault bool) bool {
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
