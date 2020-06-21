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
package main

// TODO package main // import "retrotxt.com/go" -> godoc.org/github.com/bengarrett/retrotxtgo

import (
	"github.com/bengarrett/retrotxtgo/lib/cmd"
	ver "github.com/bengarrett/retrotxtgo/lib/version"
)

// GoReleaser ldflags flags
var (
	version = "0.0.0"
	commit  = "n/a"
	date    = "n/a"
)

func main() {
	ver.B.Version = version
	ver.B.Commit = commit
	ver.B.Date = date
	cmd.Execute()
}
