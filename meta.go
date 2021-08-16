package meta

type Release struct {
	// Version of RetroTxt.
	Version string
	// GitHub commit checksum.
	Commit string
	// Date in the RFC3339 format.
	Date string
	// Built by (goreleaser).
	BuiltBy string
}

// App contains the version release and build metadata.
var App = Release{}
