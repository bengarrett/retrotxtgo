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

/*
const (
	name = "blob.go"
	dir  = "../../static"
)

//blog.go
//package pack
func init(){
	pack.Add("css/styles.css", []byte{...})
}

x := pack.Get("css/styles.css")
*/

import (
	"retrotxt.com/retrotxt/lib/cmd"
	ver "retrotxt.com/retrotxt/lib/version"
)

// goreleaser generated ldflags containers
// https://goreleaser.com/environment/#using-the-mainversion
var version, commit, date string

func main() {
	if version != "" {
		ver.B.Version = version
	}
	if commit != "" {
		ver.B.Commit = commit
	}
	if date != "" {
		ver.B.Date = date
	}
	cmd.Execute()
}
