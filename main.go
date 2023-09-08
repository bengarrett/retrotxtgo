// Copyright © 2020-2023 Ben Garrett. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
// nolint:gochecknoglobals
package main

import (
	"github.com/bengarrett/retrotxtgo/cmd"
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
	cmd.Execute()
}
