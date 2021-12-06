package long

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

var ConfigEdit = fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("Edit the %s configuration file.", meta.Name),
	"To change the editor program, either:",
	fmt.Sprintf("  1. Configure one by creating a %s shell environment variable.",
		str.Example("$EDITOR")),
	"  2. Set an editor in the configuration file:",
	str.Example(fmt.Sprintf("     %s config set --name=editor", meta.Bin)),
)
