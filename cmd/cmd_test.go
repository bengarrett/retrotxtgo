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

// tester initialises, runs and returns the results of the a Cmd package command.
// args are the command line arguments to test with the command.
func (t cmdT) tester(args []string) ([]byte, error) {
	color.Enable = false
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

func Test_InfoErrDir(t *testing.T) {
	t.Parallel()
	t.Run("info dir", func(t *testing.T) {
		t.Parallel()
		const invalid = static + "invalid_path"
		gotB, err := infoT.tester([]string{"--format", "text", invalid})
		if err == nil {
			t.Errorf("invalid file or directory path did not return an error: %s", invalid)
			t.Error(gotB)
		}
	})
}

func Test_InfoFiles(t *testing.T) {
	t.Parallel()
	t.Run("info multiple files", func(t *testing.T) {
		t.Parallel()
		gotB, err := infoT.tester([]string{"--format", "color", file1, file2})
		if err != nil {
			t.Errorf("info arguments threw an unexpected error: %s", err)
		}
		files := []string{filepath.Base(file1), filepath.Base(file2)}
		for _, f := range files {
			if !bytes.Contains(gotB, []byte(f)) {
				t.Errorf("could not find filename in the info result, want: %q", f)
			}
		}
	})
}

func Test_InfoSamples(t *testing.T) {
	t.Parallel()
	samplers := []string{"037", "ansi.aix", "shiftjis", "utf8"}
	wants := []string{
		"EBCDIC encoded text document",
		"Text document with ANSI controls",
		"plain text document",
		"UTF-8 compatible",
	}
	t.Run("info multiple samples", func(t *testing.T) {
		t.Parallel()
		for i, sample := range samplers {
			gotB, err := infoT.tester([]string{"--format", "text", sample})
			if err != nil {
				t.Error(err)
			}
			if !bytes.Contains(gotB, []byte(wants[i])) {
				t.Errorf("sample %s result does not contain: %s", sample, wants[i])
			}
		}
	})
}

func Test_InfoText(t *testing.T) {
	t.Parallel()
	t.Run("info format text", func(t *testing.T) {
		t.Parallel()
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
