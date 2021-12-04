package dirs

import "strings"

// Windows appends Windows style syntax to the directory.
func Windows(i int, p, platform, dir string) (s string, cont bool) {
	if platform == "windows" {
		if len(p) == 2 && p[1:] == ":" {
			dir = strings.ToUpper(p) + "\\"
			return dir, true
		}
		if dir == "" && i > 0 {
			dir = p + "\\"
			return dir, true
		}
	}
	return dir, false
}
