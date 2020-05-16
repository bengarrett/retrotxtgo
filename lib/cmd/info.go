package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const infoFormats = "color, json, json.min, text, xml"

var (
	infoFmt  string
	fileName string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on a text file",
	Run: func(cmd *cobra.Command, args []string) {
		if fileName == "" {
			if cmd.Flags().NFlag() == 0 {
				fmt.Printf("%s\n\n", cmd.Short)
				err := cmd.Usage()
				CheckErr(err)
				os.Exit(0)
			}
			err := cmd.Usage()
			CheckErr(err)
			FileMissingErr()
		}
		if err := infoPrint(fileName, infoFmt); err != nil {
			if fmt.Sprint(err) == "format:invalid" {
				logs.ChkArg("format", config.Format.Info)
			} else {
				logs.ChkErr(fmt.Sprintf("--name=%s is invalid,", fileName), err)
			}
		}
	},
}

func init() {
	InitDefaults()
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&fileName, "name", "n", "", cp("text file to analyse")+" (required)\n")
	infoCmd.Flags().StringVarP(&infoFmt, "format", "f", viper.GetString("info.format"), "output format \noptions: "+ci(infoFormats))
	err := viper.BindPFlag("info.format", infoCmd.Flags().Lookup("format"))
	CheckErr(err)
	infoCmd.Flags().SortFlags = false
}

func infoPrint(filename, format string) (err error) {
	d, err := info.File(filename)
	if err != nil {
		return err
	}
	switch format {
	case "color", "c", "":
		fmt.Printf("%s", d.Text(true))
	case "json", "j":
		fmt.Printf("%s\n", d.JSON(true))
	case "json.min", "jm":
		fmt.Printf("%s\n", d.JSON(false))
	case "text", "t":
		fmt.Printf("%s", d.Text(false))
	case "xml", "x":
		data, _ := d.XML()
		fmt.Printf("%s\n", data)
	default:
		return errors.New("format:invalid")
	}
	return err
}
