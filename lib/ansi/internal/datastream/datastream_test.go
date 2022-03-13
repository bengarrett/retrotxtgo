package datastream

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/ansi/internal/sgr"
)

func TestNumber(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want bool
	}{
		{"empty", nil, false},
		{"1", []byte("1"), true},
		{"1A", []byte("1A"), false},
		{"leading 0", []byte("0941"), true},
		{"float", []byte("9.41"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Number(tt.b); got != tt.want {
				t.Errorf("Number() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCUU(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want uint8
	}{
		{"empty", nil, 0},
		{"random", []byte("the quick brown..."), 0},
		{"false1", []byte("A"), 0},
		{"false2", []byte("abcde!A"), 0},
		{"false3", []byte("123pA"), 0},
		{"0", []byte("0A"), 0},
		{"5", []byte("5A"), 5},
		{"too big", []byte("257A"), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CUU(tt.b); got != tt.want {
				t.Errorf("CUU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCUB(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want uint8
	}{
		{"lowercase", []byte("10d"), 0},
		{"0", []byte("0D"), 0},
		{"5", []byte("5D"), 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CUB(tt.b); got != tt.want {
				t.Errorf("CUB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCUP(t *testing.T) {
	tests := []struct {
		name     string
		b        []byte
		wantLine uint8
		wantCol  uint8
	}{
		{"empty", nil, 0, 0},
		{"random", []byte("the quick brown..."), 0, 0},
		{"false1", []byte("abcde!H"), 0, 0},
		{"false2", []byte("123pH"), 0, 0},
		{"home", []byte("H"), 1, 1},
		{"line 6", []byte("6H"), 6, 1},
		{"col 12", []byte(";12H"), 1, 12},
		{"invalid col", []byte(";H"), 0, 0},
		{"line 6, col 12", []byte("6;12H"), 6, 12},
		{"end of page", []byte("99;99H"), 99, 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotCol := CUP(tt.b)
			if gotLine != tt.wantLine {
				t.Errorf("CUP() gotLine = %v, want %v", gotLine, tt.wantLine)
			}
			if gotCol != tt.wantCol {
				t.Errorf("CUP() gotCol = %v, want %v", gotCol, tt.wantCol)
			}
		})
	}
}

func TestED(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want Erase
	}{
		{"empty", nil, -1},
		{"0", []byte("0J"), 0},
		{"1", []byte("1J"), 1},
		{"2", []byte("2J"), 2},
		{"invalid 3", []byte("3J"), -1},
		{"lower", []byte("2j"), -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPs := ED(tt.b); gotPs != tt.want {
				t.Errorf("ED() = %v, want %v", gotPs, tt.want)
			}
		})
	}
}

func TestSGR(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []sgr.Ps
	}{
		{"empty", nil, nil},
		{"reset", []byte("0m"), []sgr.Ps{0}},
		{"invalid", []byte("1234m"), nil},
		{"bold bg/fg", []byte("37;40;1"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SGR(tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SGR() = %v, want %v", got, tt.want)
			}
		})
	}
}
