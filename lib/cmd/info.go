// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"

	internal "retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/info"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/pack"
	"retrotxt.com/retrotxt/lib/str"
)

var infoFlag struct {
	format string
}

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:     "info [filenames]",
	Aliases: []string{"i"},
	Short:   "Information on a text file",
	Example: "  retrotxt info text.asc logo.jpg\n  retrotxt info file.txt --format=json",
	Run: func(cmd *cobra.Command, args []string) {
		// piped input from other programs
		if filesystem.IsPipe() {
			infoPipe()
		}
		checkUse(cmd, args...)
		var n info.Names
		n.Length = len(args)
		for i, arg := range args {
			n.Index = i + 1
			// internal, packed example file
			filename, err := infoPackage(arg)
			if err != nil {
				logs.Println("pack", arg, err)
				continue
			}
			if filename != "" {
				defer os.Remove(filename)
				arg = filename
			}
			if err := n.Info(arg, infoFlag.format); err.Err != nil {
				if errors.Is(err.Err, info.ErrNoFile) {
					if n.Length <= 1 {
						err.Fatal()
					}
					logs.Println("pack", arg, info.ErrNoFile)
					continue
				}
				if err := cmd.Usage(); err != nil {
					logs.Println("command", "usage", err)
				}
				err.Fatal()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	i := config.Format().Info
	infoCmd.Flags().StringVarP(&infoFlag.format, "format", "f", "color",
		str.Options("output format", true, i[:]...))
}

func infoPackage(name string) (filename string, err error) {
	var s = strings.ToLower(name)
	if _, err = os.Stat(s); !os.IsNotExist(err) {
		return "", nil
	}
	pkg, exist := pack.Map()[s]
	if !exist {
		return "", nil
	}
	b := internal.Get(pkg.Name)
	if b == nil {
		return "", fmt.Errorf("view package %q: %w", pkg.Name, ErrPackGet)
	}
	file, err := ioutil.TempFile("", fmt.Sprintf("retrotxt_%s.*.txt", s))
	if err != nil {
		return "", fmt.Errorf("view package %q: %w", pkg.Name, ErrTempOpen)
	}
	if _, err = file.Write(b); err != nil {
		return "", fmt.Errorf("view package %q: %w", pkg.Name, ErrTempWrite)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("view package %q: %w", pkg.Name, ErrTempClose)
	}
	return file.Name(), nil
}

func infoPipe() {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.Fatal("info", "read stdin", err)
	}
	if err = info.Stdin(infoFlag.format, b...); err != nil {
		logs.Fatal("info", "parse stdin", err)
	}
	os.Exit(0)
}
