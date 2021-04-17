package humanize

import (
	"fmt"
	"strings"
	"time"
)

// 2 Jan 15:04 2006.
const (
	DMY = "2 Jan 2006"
	YMD = "2006 Jan 2"
	MDY = "Jan 2 2006"
	H12 = "3:04 pm"
	H24 = "15:04"
)

// DMY12 12-hour day month year.
func DMY12() string { return fmt.Sprintf("%s %s", DMY, H12) }

// DMY24 24-hour day month year.
func DMY24() string { return fmt.Sprintf("%s %s", DMY, H24) }

// YMD12 12-hour year month day.
func YMD12() string { return fmt.Sprintf("%s %s", YMD, H12) }

// YMD24 24-hour year month day.
func YMD24() string { return fmt.Sprintf("%s %s", YMD, H24) }

// MDY12 12-hour month day year.
func MDY12() string { return fmt.Sprintf("%s %s", MDY, H12) }

// MDY24 24-hour month day year.
func MDY24() string { return fmt.Sprintf("%s %s", MDY, H24) }

// Date returns a formatted date string.
func Date(format string, t time.Time) string {
	const dmy = "DMY"
	format = strings.ToUpper(format)
	if format == "" {
		format = dmy
	}
	switch format {
	case dmy:
		return t.Format(DMY)
	case "YMD":
		return t.Format(YMD)
	case "MDY":
		return t.Format(MDY)
	}
	return ""
}

// Time returns a formatted time string.
func Time(format string, t time.Time) string {
	const def = "H24"
	format = strings.ToUpper(format)
	if format == "" {
		format = def
	}
	switch format {
	case "H12":
		return t.Format(H12)
	case def:
		return t.Format(H24)
	}
	return ""
}

// Datetime returns a formatted date and time string.
func Datetime(format string, t time.Time) string {
	const def = "DMY24"
	format = strings.ToUpper(format)
	if format == "" {
		format = def
	}
	switch format {
	case "DMY12":
		return t.Format(DMY12())
	case "YMD12":
		return t.Format(YMD12())
	case "MDY12":
		return t.Format(MDY12())
	case def:
		return t.Format(DMY24())
	case "YMD24":
		return t.Format(YMD24())
	case "MDY24":
		return t.Format(MDY24())
	}
	return ""
}
