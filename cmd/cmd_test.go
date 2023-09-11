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
	out, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}

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
	samplers := []string{"037", "ansi.aix", "shiftjis", "utf8"}
	wants := []string{
		"EBCDIC encoded text document",
		"Text document with ANSI controls",
		"plain text document",
		"UTF-8 compatible",
	}
	t.Run("info multiple samples", func(t *testing.T) {
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

func Test_InfoData(t *testing.T) { //nolint:gocognit
	type Sizes struct {
		Bytes int `json:"bytes" xml:"bytes"`
	}
	type response struct {
		Name string `json:"filename" xml:"name"`
		Size Sizes  `json:"size"     xml:"size"`
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

func Test_ListCommand(t *testing.T) {
	tests := []struct {
		name       string
		wantFormal string
	}{
		{"cp037", "IBM Code Page 037"},
		{"japanese", "Shift JIS"},
		{"utf-8", "UTF-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"codepages"}
			gotB, err := listT.tester(args)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Contains(gotB, []byte(tt.wantFormal)) {
				t.Errorf("could not find %q codepage in codepages", tt.wantFormal)
			}
		})
	}
}

func Test_ListTables(t *testing.T) { //nolint:gocognit,funlen
	tests := []struct {
		table      string
		wantHeader string
		wantRunes  []rune
		wantErr    bool
	}{
		{
			"cp037", "IBM Code Page 037 (US/Canada Latin 1) - EBCDIC",
			[]rune{'␉', '␅', '␖', 'A', '9'},
			false,
		},
		{
			"cp437", "IBM Code Page 437 (DOS, OEM-US) - Extended ASCII",
			[]rune{'☺', '♪', '╬', '¿'},
			false,
		},
		{"cp850", "850 (DOS, Latin 1)", []rune{'§', '®', '¤'}, false},
		{"cp852", "852 (DOS, Latin 2)", []rune{'Ľ', 'đ', 'ő'}, false},
		{"cp855", "855 (DOS, Cyrillic)", []rune{'Ж', 'Ы', 'ю'}, false},
		{"cp858", "858 (DOS, Western Europe)", []rune{'€', '¹', 'ë'}, false},
		{"cp860", "860 (DOS, Portuguese)", []rune{'ã', 'õ', 'Ó'}, false},
		{"cp862", "862 (DOS, Hebrew)", []rune{'א', 'ש'}, false},
		{"cp863", "863 (DOS, French Canada)", []rune{'‗', '³'}, false},
		{"cp865", "865 (DOS, Nordic)", []rune{'ø', 'Ø'}, false},
		{"cp866", "866 (DOS, Cyrillic", []rune{'Б', 'Я'}, false},
		{"cp1047", "1047 (C programming", []rune{'¬', '␓'}, false},
		{"cp1140", "1140 (US/Canada Latin 1 plus €", []rune{'€'}, false},
		{"latin1", "8859-1 (Western European", []rune{'ÿ', '«'}, false},
		{"latin2", "8859-2 (Central European", []rune{'Ł', 'ű'}, false},
		{"latin3", "8859-3 (South European", []rune{'Ħ', 'ĝ'}, false},
		{"latin4", "8859-4 (North European", []rune{'ĸ', 'Æ'}, false},
		{"cyrillic", "8859-5 (Cyrillic", []rune{'Ђ', 'й'}, false},
		{"arabic", "8859-6 (Arabic", []rune{'ي', 'ص'}, false},
		{"iso88596e", "8859-6E", []rune{'ي', 'ص'}, false},
		{"iso88596i", "8859-6I", []rune{'ي', 'ص'}, false},
		{"greek", "8859-7 (Greek", []rune{'₯', 'Ξ'}, false},
		{"hebrew", "8859-8 (Hebrew", []rune{'א', 'ת'}, false},
		{"iso88598e", "8859-8E", []rune{'א', 'ת'}, false},
		{"iso88598i", "8859-8I", []rune{'א', 'ת'}, false},
		{"latin5", "8859-9 (Turkish", []rune{'ğ', 'Ğ'}, false},
		{"latin6", "8859-10 (Nordic", []rune{'ą', 'ß'}, false},
		{"iso885911", "8859-11 (Thai", []rune{'ข', '๛', '฿'}, false},
		{"iso885912", "", []rune{}, true},
		{"iso885913", "8859-13 (Baltic Rim", []rune{'Ø', 'ų'}, false},
		{"iso885914", "8859-14 (Celtic", []rune{'Ŵ', 'ẅ'}, false},
		{"iso885915", "8859-15 (Western European, 1999", []rune{'€', 'Š', 'œ'}, false},
		{"iso885916", "8859-16 (South-Eastern European", []rune{'đ', 'Œ'}, false},
		{"koi8r", "KOI8-R (Russian)", []rune{'Ю', '╬'}, false},
		{"koi8u", "KOI8-U (Ukrainian)", []rune{'є', 'Ґ', '█'}, false},
		{"mac", "Macintosh (Mac OS Roman", []rune{'∞', '‰', '◊'}, false},
		{"windows-874", "Windows 874 (Thai", []rune{'€', '“', '”'}, false},
		{"windows-1250", "Windows 1250 (Central European", []rune{'‡', '“', '”'}, false},
		{"windows-1251", "Windows 1251 (Cyrillic", []rune{'“', '”', 'Љ'}, false},
		{"windows-1252", "Windows 1252 (Western European", []rune{'“', '”', 'ƒ'}, false},
		{"windows-1253", "Windows 1253 (Greek", []rune{'“', '”', '…'}, false},
		{"windows-1254", "Windows 1254 (Turkish", []rune{'“', '”', 'Ğ'}, false},
		{"windows-1255", "Windows 1255 (Hebrew", []rune{'“', '”', '₪'}, false},
		{"windows-1256", "Windows 1256 (Arabic", []rune{'“', '”', 'گ'}, false},
		{"windows-1257", "Windows 1257 (Baltic Rim", []rune{'“', '”', '€'}, false},
		{"windows-1258", "Windows 1258 (Vietnamese", []rune{'“', '”', '†'}, false},
		{"shift_jis", "Shift JIS (Japanese", []rune{'ﾗ', 'ｼ', 'ﾎ'}, false},
		{"utf-8", "UTF-8 - Unicode", []rune{'␀', '␟', '␗'}, false},
		// ASA/early ASCII
		{"ascii-63", "", []rune{'␉', '0', 'A'}, false},
		// TODO: confirm fix
		{"ascii-65", "", []rune{'␉', '0', 'A', 'a', '{'}, false},
		{"ascii-67", "", []rune{'␉', '0', 'A', 'a', '~'}, false},
		// // UTF16 & UTF32 tables are not supported
		{"utf-16", "", []rune{}, true},
		{"utf-16be", "", []rune{}, true},
		{"utf-16le", "", []rune{}, true},
		{"utf-32", "", []rune{}, true},
		{"utf-32be", "", []rune{}, true},
		{"utf-32le", "", []rune{}, true},
	}
	// test the list tables command
	args := []string{"tables"}
	gotB, err := listT.tester(args)
	if err != nil {
		t.Error(err)
		return
	}
	if len(gotB) == 0 {
		t.Errorf("TABLES, the result of the list tables command is empty, or 0 bytes")
	}
	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			if !bytes.Contains(gotB, []byte(tt.wantHeader)) {
				t.Errorf("TABLES, could not find %q header in the tables output", tt.wantHeader)
				return
			}
			// find duplicate tables
			if len(tt.wantHeader) == 0 {
				return
			}
			const expected = 1
			if finds := bytes.Count(gotB, []byte(tt.wantHeader)); finds > expected {
				t.Errorf("TABLES, %d instances of the %q table were displayed in the tables output", finds, tt.wantHeader)
			}
			// don't check for individual runes,
			// as all the tables are dumped together, there are too many duplicates
		})
	}
	// test the list table command
	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			args := []string{"table", tt.table}
			gotB, err := listT.tester(args)
			if gotE := (err != nil); gotE != tt.wantErr {
				t.Errorf("TABLE %q, returned error %v, wanted %v", tt.table, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			// find empty results
			if len(gotB) == 0 {
				t.Errorf("TABLE %q, the returned table is empty, or 0 bytes", tt.table)
				return
			}
			// find invalid or unexpected results
			if !bytes.Contains(gotB, []byte(tt.wantHeader)) {
				t.Errorf("TABLE, could not find %q header in the table", tt.wantHeader)
			}
			// confirm a couple of unique runes displayed in the table
			for i, runes := range tt.wantRunes {
				if !bytes.ContainsRune(gotB, runes) {
					t.Errorf("%d. result doesn't include the expected rune, %q got:\n%s", i, runes, gotB)
				}
			}
		})
	}
}
