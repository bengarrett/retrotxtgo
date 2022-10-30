package cmd_test

import (
	"bytes"
	"io/ioutil"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type cmdT int

const (
	confT cmdT = iota
	creaT
	infoT
	listT
	viewT
)

// tester initialises, runs and returns the results of the a Cmd package command.
// args are the command line arguments to test with the command.
func (t cmdT) tester(args []string) ([]byte, error) {
	color.Enable = false
	var c *cobra.Command
	b := bytes.NewBufferString("")
	switch t {
	case confT:
		c = cmd.ConfigInit()
	// case creaT:
	// 	c = cmd.CreateInit()
	case infoT:
		c = cmd.InfoInit()
	case listT:
		c = cmd.ListInit()
	case viewT:
		c = cmd.ViewInit()
	default:
	}
	c = cmd.Tester(c)
	c.SetOut(b)
	if len(args) > 0 {
		c.SetArgs(args)
	}
	if err := c.Execute(); err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}
