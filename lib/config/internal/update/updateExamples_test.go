package update_test

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/update"
	"github.com/gookit/color"
)

func Example_updateBool() {
	color.Enable = false
	fmt.Print(update.Bool(false, "example"))
	// Output: The example is not in use.
}

func Example_updateString() {
	color.Enable = false
	update.String(os.Stdout, "", "example", "")
	update.String(os.Stdout, "x", get.SaveDir, "")
	// Output: The empty example setting is not in use.
	//
	//   The save_directory is set to "". ✗
}
