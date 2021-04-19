// Copyright © 2020 Ben Garrett. All rights reserved.
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
	"retrotxt.com/retrotxt/lib/cmd"
	v "retrotxt.com/retrotxt/lib/version"
)

// goreleaser generated ldflags containers.
// https://goreleaser.com/environment/#using-the-mainversion
var version, commit, date, builtBy string

func main() {
	if version != "" {
		v.Release.Version = version
	}
	if commit != "" {
		v.Release.GitHash = commit
	}
	if date != "" {
		v.Release.Date = date
	}
	if builtBy != "" {
		v.Release.BuiltBy = builtBy
	}
	cmd.Execute()
}
