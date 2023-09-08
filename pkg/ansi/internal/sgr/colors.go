package sgr

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

var ErrNotRGB = errors.New("int is not a valid red-green-blue value")

// Colors is currently unused.
func Colors(p Ps, next [][]byte) (value int, length int) {
	const xterm256, rgb = 5, 2
	const xtermLen, rgbLen = 2, 5
	if p != 38 && p != 48 {
		return -1, 0
	}
	if len(next) < 2 {
		return -1, 0
	}
	s, err := strconv.Atoi(string(next[1]))
	if err != nil {
		return -1, 0
	}

	fmt.Println("ps", p, len(next), "s:", s)
	switch s {
	case xterm256:
		value = XTerm256(next[2:])
		if value != -1 {
			return value, xtermLen
		}
	case rgb:
		value = RGB(next[2:])
		if value != -1 {
			return value, rgbLen
		}
	}
	return -1, 0
}

// RGB reads b
func RGB(b [][]byte) int {
	if len(b) < 3 {
		return -1
	}
	var red, green, blue uint8
	for count, v := range b {
		if count > 3 {
			break
		}
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return -1
		}
		if !Uint8(i) {
			return -1
		}
		switch count {
		case 0:
			red = uint8(i)
		case 1:
			green = uint8(i)
		case 2:
			blue = uint8(i)
		}
	}
	return RGBDecimal(red, green, blue)
}

// XTerm256
func XTerm256(b [][]byte) int {
	if len(b) < 1 {
		return -1
	}
	i, err := strconv.Atoi(string(b[0]))
	if err != nil {
		return -1
	}
	switch Uint8(i) {
	case true:
		return i
	default:
		return -1
	}
}

// 38;5;0…255 Extension
// 48;5;0…255 ExtensionB
// 48;2;R;G;B;

// RGBDecimal
func RGBDecimal(r uint8, g uint8, b uint8) int {
	red := int(r)
	green := int(g)
	blue := int(b)
	return red*65536 + green*256 + blue
}

// DecimalRGB
func DecimalRGB(f float64) (r uint8, g uint8, b uint8, err error) {
	red := math.Floor(f / (256 * 256))
	green := math.Mod(math.Floor(f/256), 256)
	blue := math.Mod(f, 256)
	if !Uint8(int(red)) {
		return 0, 0, 0, ErrNotRGB
	}
	return uint8(red), uint8(green), uint8(blue), nil
}
