package bg

type Colors uint

const (
	System Colors = 40
	IbmAIX Colors = 100
)

// System colors
const (
	Black Colors = iota + System
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// IBM AIX bright colors
const (
	BrightBlack Colors = iota + IbmAIX
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)
