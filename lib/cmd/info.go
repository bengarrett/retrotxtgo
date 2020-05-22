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

var (
	infoFormat   string
	infoFilename string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on a text file",
	Run: func(cmd *cobra.Command, args []string) {
		if infoFilename == "" {
			if cmd.Flags().NFlag() == 0 {
				fmt.Printf("%s\n\n", cmd.Short)
				err := cmd.Usage()
				logs.ReCheck(err)
				os.Exit(0)
			}
			err := cmd.Usage()
			logs.ReCheck(err)
			logs.FileMissingErr()
		}
		if err := infoPrint(infoFilename, infoFormat); err != nil {
			if fmt.Sprint(err) == "format:invalid" {
				logs.CheckFlag("format", infoFilename, config.Format.Info)
			} else {
				//logs.CheckCmd(err)
			}
			// else {
			// 	logs.Check(fmt.Sprintf("--name=%s is invalid,", infoFilename), err)
			// }
		}
	},
}

func init() {
	config.InitDefaults()
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoFilename, "name", "n", "",
		logs.Cp("text file to analyse")+" (required)\n")
	infoCmd.Flags().StringVarP(&infoFormat, "format", "f", viper.GetString("style.yaml"),
		"output format \noptions: "+logs.Ci(config.Format.String("info")))
	err := viper.BindPFlag("style.yaml", infoCmd.Flags().Lookup("format"))
	logs.ReCheck(err)
	infoCmd.Flags().SortFlags = false
}

func infoPrint(filename, format string) (err error) {
	var d info.Detail
	if err := d.Read(filename); err != nil {
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
