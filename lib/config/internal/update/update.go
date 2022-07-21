package update

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

const namedFile = "config.yaml"

// Bool the boolean value of the named setting.
func Bool(b bool, name string) string {
	switch b {
	case true:
		return fmt.Sprintf("\n  The %s is enabled.\n", str.ColFuz(name))
	default:
		return fmt.Sprintf("\n  The %s is not in use.\n", str.ColFuz(name))
	}
}

// String the string value of the named setting.
func String(w io.Writer, s, name, value string) {
	const sd = get.SaveDir
	switch s {
	case "":
		fmt.Fprintf(w, "\n  The empty %s setting is not in use.\n",
			str.ColFuz(name))
		if name == sd {
			fmt.Fprintf(w, "  Files created by %s will always be saved to the active directory.\n",
				meta.Name)
		}
	default:
		fmt.Fprintf(w, "\n  The %s is set to %q.", str.ColFuz(name), value)
		// print the operating system's ability to use the existing set values
		// does the 'editor' exist in the env path, does the save_directory exist?
		switch name {
		case "editor":
			_, err := exec.LookPath(value)
			fmt.Fprint(w, " ", str.Bool(err == nil))
		case sd:
			f := false
			if _, err := os.Stat(value); !os.IsNotExist(err) {
				f = true
			}
			fmt.Fprint(w, " ", str.Bool(f))
		default:
		}
		fmt.Fprintln(w)
	}
}
