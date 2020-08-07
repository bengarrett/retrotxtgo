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

var (
	// DMY12 12-hour day month year.
	DMY12 = fmt.Sprintf("%s %s", DMY, H12)
	// DMY24 24-hour day month year.
	DMY24 = fmt.Sprintf("%s %s", DMY, H24)
	// YMD12 12-hour year month day.
	YMD12 = fmt.Sprintf("%s %s", YMD, H12)
	// YMD24 24-hour year month day.
	YMD24 = fmt.Sprintf("%s %s", YMD, H24)
	// MDY12 12-hour month day year.
	MDY12 = fmt.Sprintf("%s %s", MDY, H12)
	// MDY24 24-hour month day year.
	MDY24 = fmt.Sprintf("%s %s", MDY, H24)
)

// Date returns a formatted date string.
func Date(format string, t time.Time) string {
	format = strings.ToUpper(format)
	if format == "" {
		format = "DMY"
	}
	switch format {
	case "DMY":
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
	format = strings.ToUpper(format)
	if format == "" {
		format = "H24"
	}
	switch format {
	case "H12":
		return t.Format(H12)
	case "H24":
		return t.Format(H24)
	}
	return ""
}

// Datetime returns a formatted date and time string.
func Datetime(format string, t time.Time) string {
	format = strings.ToUpper(format)
	if format == "" {
		format = "DMY24"
	}
	switch format {
	case "DMY12":
		return t.Format(DMY12)
	case "YMD12":
		return t.Format(YMD12)
	case "MDY12":
		return t.Format(MDY12)
	case "DMY24":
		return t.Format(DMY24)
	case "YMD24":
		return t.Format(YMD24)
	case "MDY24":
		return t.Format(MDY24)
	}
	return ""
}
