package tests

import (
	"fmt"
	"log"
	"testing"
)

func ExampleUTF8() {
	result, _, err := UTF8(utf8)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-utf8.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%dB %s", len(t), t)
	// Output: 13B [☠|☮|♺]
}

func ExampleUTF16LE() {
	result, _, err := UTF16LE(utf8)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-utf16le.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(t))
	// Output: 16
}

func ExampleUTF16BE() {
	result, _, err := UTF16BE(utf8)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-utf16be.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(t))
	// Output: 16
}

func ExampleCP437Out() {
	result, err := CP437Out(toCP437)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-cp437.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(t))
	// Output: 7
}

func ExampleCP437In() {
	result, err := CP437In(fromCP437)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-cp437In.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(t)
	// Output: ═╣▓╠═
}

func BenchmarkCP437In(b *testing.B) {
	result, err := CP437In(fromCP437)
	if err != nil {
		log.Fatal(err)
	}
	name := "sample-cp437In.txt"
	SaveFile(result, name)
	t, err := ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(t)
}

func ExampleEmbedText() {
	fmt.Printf("%s\n", EmbedText())
	// Output: xxxsdd
}

func ExampleLogoASCII() {
	bin, err := BinaryOut(LogoASCII)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(bin))
	// Output: x
}

func ExampleLogoANSI() {
	bin, err := BinaryOut(LogoANSI)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(bin))
	// Output: x
}
