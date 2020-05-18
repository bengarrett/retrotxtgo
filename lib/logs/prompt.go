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
	input, prompts := "", 0
	for {
		input = ""
		prompts++
		if i, err := fmt.Scanln(&input); err != nil {
			log.Fatalf("logs.promptport scanln at %d: %s", i, err)
		}
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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		words := scanner.Text()
		switch words {
		case "":
			os.Exit(0)
		case "-":
			words = ""
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading standard input: %s", err)
	}
	return words
}

// PromptStrings asks for and returns a single choice from a collection of keys.
func PromptStrings(keys []string) (key string) {
	prompts := 0
	for {
		key = ""
		prompts++
		if i, err := fmt.Scanln(&key); err != nil {
			log.Fatalf("logs.promptstrings scanln at %d: %s", i, err)
		}
		if key == "" {
			promptCheck(prompts)
			continue
		}
		var i = sort.SearchStrings(keys, key)
		if i >= len(keys) || keys[i] != key {
			fmt.Printf("%s %v\n", Ce("✗"), key)
			promptCheck(prompts)
			continue
		}
		return key
	}
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
