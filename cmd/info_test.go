package cmd_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	static = "../static"
	file1  = "../static/ansi/ansi-cp.ans"
	file2  = "../static/bbs/SHEET.ANS"
	file3  = "../static/font/mona.woff2"
)

func Test_InfoErrDir(t *testing.T) {
	t.Run("info dir", func(t *testing.T) {
		const invalid = static + "invalid_path"
		gotB, err := infoT.tester([]string{"--format", "text", invalid})
		if err == nil {
			t.Errorf("invalid file or directory path did not return an error: %s", invalid)
			t.Error(gotB)
		}
	})
}

func Test_InfoFiles(t *testing.T) {
	t.Run("info multiple files", func(t *testing.T) {
		gotB, err := infoT.tester([]string{"--format", "color", file1, file2, file3})
		if err != nil {
			t.Errorf("info arguments threw an unexpected error: %s", err)
		}
		files := []string{filepath.Base(file1), filepath.Base(file2), filepath.Base(file3)}
		for _, f := range files {
			if !bytes.Contains(gotB, []byte(f)) {
				t.Errorf("could not find filename in the info result, want: %q", f)
			}
		}
	})
}

func Test_InfoText(t *testing.T) {
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
	type Sizes struct {
		Bytes int `json:"bytes" xml:"bytes"`
	}
	type response struct {
		Name string `json:"filename" xml:"name"`
		Size Sizes  `json:"size" xml:"size"`
	}
	t.Run("info format json/xml", func(t *testing.T) {
		err := filepath.Walk(static,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				// test --format=json
				gotJson, err := infoT.tester([]string{"--format", "json", path})
				if err != nil {
					t.Error(err)
					return nil
				}
				res := response{}
				if err := json.Unmarshal(gotJson, &res); err != nil {
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
					fmt.Print(string(gotXML))
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
