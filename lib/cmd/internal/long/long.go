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

var Root = fmt.Sprintf(`Turn many pieces of ANSI art, ASCII and NFO texts into HTML5 using %s.
It is the platform agnostic tool that takes nostalgic text files and stylises
them into a more modern, useful format to view or copy in a web browser.`, meta.Name)
