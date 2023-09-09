package cmd_test

import (
	"bytes"
	"testing"
)

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

func Test_ListExamples(t *testing.T) {
	tests := []struct {
		name       string
		wantFormal string
	}{
		{"037", "EBCDIC 037 IBM mainframe test"},
		{"1252", "Windows-1252 English test"},
		{"ansi", "RetroTxtGo 256 color ANSI logo"},
		{"ansi.rgb", "ANSI RGB 24-bit color sheet"},
		{"sauce", "SAUCE metadata test"},
		{"utf8", "UTF-8 test with no Byte Order Mark"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"examples"}
			gotB, err := listT.tester(args)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Contains(gotB, []byte(tt.wantFormal)) {
				t.Errorf("could not find %q example in the examples", tt.wantFormal)
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
		{"cp037", "IBM Code Page 037 (US/Canada Latin 1) - EBCDIC", []rune{'␉', '␅', '␖', 'A', '9'}, false},
		{"cp437", "IBM Code Page 437 (DOS, OEM-US) - Extended ASCII", []rune{'☺', '♪', '╬', '¿'}, false},
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
		{"ascii-65", "", []rune{'␉', '0', 'A', 'a', '{'}, false},
		{"ascii-67", "", []rune{'␉', '0', 'A', 'a', '~'}, false},
		// UTF16 & UTF32 tables are not supported
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
