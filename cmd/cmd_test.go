package cmd_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"testing"

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

	static = "../static"
	file1  = "../static/ansi/ansi-cp.ans"
	file2  = "../static/bbs/SHEET.ANS"
)

func init() {
	color.Enable = false
}

// tester initialises, runs and returns the results of the a Cmd package command.
// args are the command line arguments to test with the command.
func (t cmdT) tester(args []string) ([]byte, error) {
	c := &cobra.Command{}
	b := &bytes.Buffer{}
	switch t {
	case infoT:
		c = cmd.InfoInit()
	// case listT:
	// 	c = cmd.ListInit()
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
	out, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Test_InfoText(t *testing.T) { //nolint:paralleltest
	t.Run("info format text", func(t *testing.T) {
		err := filepath.Walk(static,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				gotB, err := infoT.tester([]string{"--format", "text", path})
				if err != nil {
					t.Error(err)
					return nil
				}
				if len(gotB) == 0 {
					t.Errorf("info --format=text %s returned nothing", path)
				}
				if !bytes.Contains(gotB, []byte(info.Name())) {
					t.Errorf("could not find filename in the info result, want: %q", info.Name())
				}
				mod := info.ModTime().UTC().Format("2 Jan 2006")
				if !bytes.Contains(gotB, []byte(mod)) {
					t.Errorf("could not find the modified time in the info result, want: %q", mod)
				}
				return nil
			})
		if err != nil {
			t.Errorf("walk error: %s", err)
		}
	})
}

func Test_InfoData(t *testing.T) {
	t.Parallel()
	type Sizes struct {
		Bytes int `json:"bytes" xml:"bytes"`
	}
	type response struct {
		Name string `json:"filename" xml:"name"`
		Size Sizes  `json:"size"     xml:"size"`
	}
	t.Run("info format json/xml", func(t *testing.T) {
		t.Parallel()
		err := filepath.Walk(static,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				// test --format=json
				gotJSON, err := infoT.tester([]string{"--format", "json", path})
				if err != nil {
					t.Error(err)
					return nil
				}
				res := response{}
				if err := json.Unmarshal(gotJSON, &res); err != nil {
					t.Error(err)
				}
				if res.Name != info.Name() {
					t.Errorf("could not find filename in the json result, want: %q", info.Name())
				}
				if int64(res.Size.Bytes) != info.Size() {
					t.Errorf("could not find file size in the json result, want: %q", info.Size())
				}
				// test --format=xml
				gotXML, err := infoT.tester([]string{"--format", "xml", path})
				if err != nil {
					t.Error(err)
					return nil
				}
				res = response{}
				if err := xml.Unmarshal(gotXML, &res); err != nil {
					t.Error(err)
				}
				if res.Name != info.Name() {
					t.Errorf("could not find filename in the xml result, want: %q", info.Name())
				}
				if int64(res.Size.Bytes) != info.Size() {
					t.Errorf("could not find file size in the xml result, want: %q", info.Size())
				}
				return nil
			})
		if err != nil {
			t.Errorf("walk error: %s", err)
		}
	})
}
