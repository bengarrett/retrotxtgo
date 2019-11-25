package filesystem

import (
	"bytes"
	"testing"
)

func TestReadAllBytesEmpty(t *testing.T) {
	r, err := ReadAllBytes("")
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %s but got %s", e, r)
	}
	if err == nil {
		t.Fatalf("Expected err but got nil")
	}
}

func TestReadAllBytes(t *testing.T) {
	r, err := ReadAllBytes("../textfiles/hi.txt")
	var e = []byte("Hello world ☺\n☺ ሰላም ልዑል")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err != nil {
		t.Fatalf("Expected nil error, %q", err)
	}
}

func TestTailBytesEmpty(t *testing.T) {
	r, err := TailBytes("", 0)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %s but got %s", e, r)
	}
	if err == nil {
		t.Fatalf("Expected err but got nil")
	}
}

func TestTailBytes0(t *testing.T) {
	r, err := TailBytes("../textfiles/hi.txt", 0)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err != nil {
		t.Fatalf("Expected nil error, %q", err)
	}
}

func TestTailBytesPositive(t *testing.T) {
	r, err := TailBytes("../textfiles/hi.txt", 10)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err == nil {
		t.Fatalf("Expected an error, EOF")
	}
}

func TestTailBytesRange(t *testing.T) {
	_, err := TailBytes("../textfiles/hi.txt", -999999999)
	if err == nil {
		t.Fatalf("Expected an error, offset: value is too large")
	}
}

func TestTailBytesNegative(t *testing.T) {
	r, err := TailBytes("../textfiles/hi.txt", -9)
	var e = []byte("ልዑል")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err != nil {
		t.Fatalf("Expected nil error, %q", err)
	}
}
