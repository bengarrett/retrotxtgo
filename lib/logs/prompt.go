package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// PromptPort ...
func PromptPort(validate bool) (port uint) {
	var input string
	cnt := 0
	for {
		input = ""
		cnt++
		fmt.Scanln(&input)
		if input == "" {
			promptCheck(cnt)
			continue
		}
		i, err := strconv.ParseInt(input, 10, 0)
		if err != nil {
			fmt.Printf("%s %v\n", Ce("✗"), input)
			promptCheck(cnt)
			continue
		}
		if validate {
			if v := PortValid(uint(i)); !v {
				fmt.Printf("%s %v, is out of range\n", Ce("✗"), input)
				promptCheck(cnt)
				continue
			}
		}
		return uint(i)
	}
}

// PortValid ...
func PortValid(p uint) bool {
	if p < PortMin || p > PortMax {
		return false
	}
	return true
}

// PromptString ...
func PromptString() (value string) {
	// allow multiple words
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		switch txt {
		case "":
			os.Exit(0)
		case "-":
			value = ""
		default:
			value = txt
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	return value
}

// PromptStrings ...
func PromptStrings(keys []string) (value string) {
	cnt := 0
	for {
		value = ""
		cnt++
		fmt.Scanln(&value)
		if value == "" {
			promptCheck(cnt)
			continue
		}
		var i = sort.SearchStrings(keys, value)
		if i >= len(keys) || keys[i] != value {
			fmt.Printf("%s %v\n", Ce("✗"), value)
			promptCheck(cnt)
			continue
		}
		return value
	}
}

func promptCheck(cnt int) {
	switch {
	case cnt == 2:
		fmt.Println("Ctrl+C to keep the existing value")
	case cnt >= 4:
		os.Exit(1)
	}
}

// PromptYN asks the user for a yes or no input.
func PromptYN(query string, yesDefault bool) bool {
	var y, n string = "Y", "n"
	if !yesDefault {
		y, n = "y", "N"
	}
	fmt.Printf("%s? [%s/%s] ", query, y, n)
	input, err := promptRead(os.Stdin)
	Log(err)
	return promptyn(input, yesDefault)
}

func promptyn(input string, yesDefault bool) bool {
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

func promptRead(stdin io.Reader) (input string, err error) {
	reader := bufio.NewReader(stdin)
	input, err = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil && err != io.EOF {
		return input, err
	}
	return input, nil
}
