// nolint:dupl
package humanize

import (
	"testing"
	"time"
)

var (
	nyd      = time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC) //nolint:gochecknoglobals
	midnight = time.Date(2020, 1, 1, 24, 0, 0, 0, time.UTC) //nolint:gochecknoglobals
)

func TestDate(t *testing.T) {
	type args struct {
		format string
		t      time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nyd empty", args{"", nyd}, "1 Jan 2020"},
		{"nyd invalid", args{"???", nyd}, ""},
		{"nyd DMY", args{"DMY", nyd}, "1 Jan 2020"},
		{"nyd YMD", args{"YMD", nyd}, "2020 Jan 1"},
		{"nyd MDY", args{"MDY", nyd}, "Jan 1 2020"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Date(tt.args.format, tt.args.t); got != tt.want {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		format string
		t      time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nyd empty", args{"", nyd}, "12:00"},
		{"nyd 12", args{"h12", nyd}, "12:00 pm"},
		{"nyd 24", args{"h24", nyd}, "12:00"},
		{"midnight 12", args{"h12", midnight}, "12:00 am"},
		{"midnight 24", args{"h24", midnight}, "00:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.format, tt.args.t); got != tt.want {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatetime(t *testing.T) {
	type args struct {
		format string
		t      time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nyd empty", args{"", nyd}, "1 Jan 2020 12:00"},
		{"dmy 12", args{"DMY12", nyd}, "1 Jan 2020 12:00 pm"},
		{"ymd 12", args{"YMD12", nyd}, "2020 Jan 1 12:00 pm"},
		{"mdy 12", args{"MDY12", nyd}, "Jan 1 2020 12:00 pm"},
		{"dym 24", args{"DMY24", nyd}, "1 Jan 2020 12:00"},
		{"ymd 24", args{"ymd24", midnight}, "2020 Jan 2 00:00"},
		{"mdy 24", args{"MDY24", midnight}, "Jan 2 2020 00:00"},
		{"error", args{"xxx", midnight}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Datetime(tt.args.format, tt.args.t); got != tt.want {
				t.Errorf("Datetime() = %v, want %v", got, tt.want)
			}
		})
	}
}
