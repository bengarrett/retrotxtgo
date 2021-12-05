package upd_test

import (
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/upd"
	"github.com/gookit/color"
)

func Example_updateBool() {
	color.Enable = false
	upd.Bool(false, "example")
	// Output: The example is not in use.
}

func Example_updateString() {
	color.Enable = false
	upd.String("", "example", "")
	upd.String("x", get.SaveDir, "")
	// Output: The empty example setting is not in use.
	//
	//   The save-directory is set to "". ✗
}
