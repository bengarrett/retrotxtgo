package main

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/logs"
	"github.com/bengarrett/retrotxtgo/meta"
)

// goreleaser generated ldflags containers.
// https://goreleaser.com/cookbooks/using-main.version
var (
	version = meta.GoBuild
	commit  = meta.Placeholder
	date    = meta.Placeholder
	builtBy = "go builder"
)

func main() {
	meta.App.Version = version
	meta.App.Commit = commit
	meta.App.Date = date
	meta.App.BuiltBy = builtBy
	if err := cmd.Execute(); err != nil {
		if s := logs.Execute(err, false); s != "" {
			fmt.Fprintln(os.Stderr, s)
			os.Exit(logs.OSErr)
		}
	}
}
