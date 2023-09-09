package sgr

import (
	"errors"
	"math"
	"strconv"
)

const (
	// Color16bit is the maximum value for a 16bit color.
	Color16bit = 65536
	// Color8bit is the maximum value for a 8bit color.
	Color8bit = 256
)

var ErrNotRGB = errors.New("int is not a valid red-green-blue value")

// RGB reads b and returns the color value.
func RGB(b [][]byte) int {
	const expected = 3
	if len(b) < expected {
		return -1
	}
	var red, green, blue uint8
	for count, v := range b {
		if count > expected {
			break
		}
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return -1
		}
		if !Uint8(i) {
			return -1
		}
		const (
			r = 0
			g = 1
			b = 2
		)
		switch count {
		case r:
			red = uint8(i)
		case g:
			green = uint8(i)
		case b:
			blue = uint8(i)
		}
	}
	return RGBDecimal(red, green, blue)
}

// XTerm256 reads b and returns the color value.
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

// RGBDecimal converts red, green and blue values to a decimal color value.
func RGBDecimal(r uint8, g uint8, b uint8) int {
	red := int(r)
	green := int(g)
	blue := int(b)
	return red*Color16bit + green*Color8bit + blue
}

// DecimalRGB converts a decimal color value to red, green and blue values.
func DecimalRGB(f float64) (uint8, uint8, uint8, error) {
	red := math.Floor(f / (Color8bit * Color8bit))
	green := math.Mod(math.Floor(f/Color8bit), Color8bit)
	blue := math.Mod(f, Color8bit)
	if !Uint8(int(red)) {
		return 0, 0, 0, ErrNotRGB
	}
	return uint8(red), uint8(green), uint8(blue), nil
}
