package version

// Build ..
type Build struct {
	// Date in RFC3339
	Date string
	// Commit git SHA
	Commit string
	// Version of RetroTxt
	Version string
}

// Inf ...
var Inf = Build{
	Date:    "",
	Commit:  "",
	Version: "",
}
