package cmd

import (
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/format"
	"github.com/bengarrett/retrotxtgo/cmd/internal/info"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
)

const iCmdLong = `Discover details and information about any text or text art file.

The info command will return the following information about a text file:

- slug			A URL friendly version of the filename.
- filename		The filename.
- filetype		The file type or function, such as plain text file.
- Unicode		Whether the file is readable as Unicode.
- line break		The line break type in use, such as CRLF.
- characters		The number of characters in the file.
- words			The number of words in the file.
- size			The file size in a human readable format.
- lines			The number of lines in the file determined by the line breaks.
- width			The widest line in the file.
- modified		The date and time the file was last modified.
- media type		The IANA media type, such as text/plain.
- SHA256 check		The SHA256 integrity checksum of the file.
- CRC64			The cyclic redundancy check of the file.
- CRC32			The cyclic redundancy check of the file.
- MD5			The MD5 hash of the file.

Art scene files embedded with SAUCE metadata will also return the 
following information:

- title			The title of the file.
- author		The nickname or handle of the file creator.
- group			The group or company of the author.
- date			The date the file was created.
- original size		The size of the file without the SAUCE metadata.
- file type		The specific type of the file, such as PNG image.
- data type		The media type of the file, such as a bitmap image.
- description		The description of file and data type.
- character width	The number of characters per line.
- number of lines	The number of lines in the file.
- interpretation	Additional information about the file types.
- comments		Additional comments by the author.

The full SAUCE specification can be found at:
https://www.acid.org/info/sauce/sauce.htm`

func InfoCommand() *cobra.Command {
	s := "Information on a text file"
	expl := strings.Builder{}
	example.Info.String(&expl)
	return &cobra.Command{
		Use:     fmt.Sprintf("info %s", example.Filenames),
		Aliases: []string{"i"},
		GroupID: IDfile,
		Short:   s,
		Long:    iCmdLong,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return info.Run(cmd.OutOrStdout(), cmd, args...)
		},
	}
}

func InfoInit() *cobra.Command {
	infoc := InfoCommand()
	infos := format.Format().Info
	s := &strings.Builder{}
	term.Options(s, "print format or syntax", true, true, infos[:]...)
	infoc.Flags().StringVarP(&flag.Info.Format, "format", "f", "color", s.String())
	return infoc
}

func init() {
	Cmd.AddCommand(InfoInit())
}
