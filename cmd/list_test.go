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

func Test_ListTables(t *testing.T) {
	tests := []struct {
		table      string
		wantHeader string
		wantRunes  []rune
	}{
		{"cp037", "IBM Code Page 037 (US/Canada Latin 1) - EBCDIC", []rune{'␉', '␅', '␖', 'A', '9'}},
		{"cp437", "IBM Code Page 437 (DOS, OEM-US) - Extended ASCII", []rune{'☺', '♪', '╬', '¿'}},
		{"cp850", "850 (DOS, Latin 1)", []rune{'§', '®', '¤'}},
		{"cp852", "852 (DOS, Latin 2)", []rune{'Ľ', 'đ', 'ő'}},
		{"cp855", "855 (DOS, Cyrillic)", []rune{'Ж', 'Ы', 'ю'}},
		{"cp858", "858 (DOS, Western Europe)", []rune{'€', '¹', 'ë'}},
		{"cp860", "860 (DOS, Portuguese)", []rune{'ã', 'õ', 'Ó'}},
		{"cp862", "862 (DOS, Hebrew)", []rune{'א', 'ש'}},
		{"cp863", "863 (DOS, French Canada)", []rune{'‗', '³'}},
		{"cp865", "865 (DOS, Nordic)", []rune{'ø', 'Ø'}},
		{"cp866", "866 (DOS, Cyrillic", []rune{'Б', 'Я'}},
		{"cp1047", "1047 (C programming", []rune{'¬', '␓'}},
		{"cp1140", "1140 (US/Canada Latin 1 plus €", []rune{'€'}},
		{"latin1", "8859-1 (Western European", []rune{'ÿ', '«'}},
		{"latin2", "8859-2 (Central European", []rune{'Ł', 'ű'}},
		{"latin3", "8859-3 (South European", []rune{'Ħ', 'ĝ'}},
		{"latin4", "8859-4 (North European", []rune{'ĸ', 'Æ'}},
		{"cyrillic", "8859-5 (Cyrillic", []rune{'Ђ', 'й'}},
		{"arabic", "8859-6 (Arabic", []rune{'ي', 'ص'}},
		{"iso88596e", "8859-6E", []rune{'ي', 'ص'}},
		{"iso88596i", "8859-6I", []rune{'ي', 'ص'}},
		{"greek", "8859-7 (Greek", []rune{'₯', 'Ξ'}},
		{"hebrew", "8859-8 (Hebrew", []rune{'א', 'ת'}},
		{"iso88598e", "8859-8E", []rune{'א', 'ת'}},
		{"iso88598i", "8859-8I", []rune{'א', 'ת'}},
		{"latin5", "8859-9 (Turkish", []rune{'ğ', 'Ğ'}},
		{"latin6", "8859-10 (Nordic", []rune{'ą', 'ß'}},
		{"iso889511", "", []rune{'ข', '๛', '฿'}}, // TODO: missing in tables
		{"iso889512", "", []rune{}},
		{"iso889513", "8859-13 (Baltic Rim", []rune{'Ø', 'ų'}},                  // TODO: iso889513, iso-8895-13 fails; 13 "ISO 8859-13" works
		{"iso889514", "8859-14 (Celtic", []rune{'Ŵ', 'ẅ'}},                      // TODO: same as iso 13
		{"iso889515", "8859-15 (Western Eruopean, 1999", []rune{'€', 'Š', 'œ'}}, // TODO: same as iso 13
		{"iso889516", "8859-16 (Soouth-Eastern Eruopean", []rune{'đ', 'Œ'}},     // TODO: same as iso 13
		{"koi8r", "KOI8-R (Russian)", []rune{'Ю', '╬'}},
		{"koi8u", "KOI8-U (Ukrainian)", []rune{'є', 'Ґ', '█'}},
		{"mac", "Macintosh (Mac OS Roman", []rune{'∞', '‰', '◊'}},
		{"windows-874", "Windows 874 (Thai", []rune{'€', '“', '”'}},
		{"windows-1250", "Windows 1250 (Central European", []rune{'‡', '“', '”'}},
		{"windows-1251", "Windows 1251 (Cyrillic", []rune{'“', '”', 'Љ'}},
		{"windows-1252", "Windows 1252 (Western European", []rune{'“', '”', 'ƒ'}},
		{"windows-1253", "Windows 1253 (Greek", []rune{'“', '”', '…'}},
		{"windows-1254", "Windows 1254 (Turkish", []rune{'“', '”', 'Ğ'}},
		{"windows-1255", "Windows 1255 (Hebrew", []rune{'“', '”', '₪'}},
		{"windows-1256", "Windows 1256 (Arabic", []rune{'“', '”', 'گ'}},
		{"windows-1257", "Windows 1257 (Baltic Rim", []rune{'“', '”', '€'}},
		{"windows-1258", "Windows 1258 (Vietnamese", []rune{'“', '”', '†'}},
		// TODO: list tables, thai windows 874 is relisted again??
		{"shift_jis", "Shift JIS (Japanese", []rune{'ﾗ', 'ｼ', 'ﾎ'}},
		{"utf-8", "UTF-8 - Unicode", []rune{'␀', '␟', '␗'}},
		// UTF16 & UTF32 tables are not supported; TODO: add *** notice on list codepages command
		// {"utf-16", "", []rune{}},
		// {"utf-16be", "", []rune{}},
		// {"utf-16le", "", []rune{}},
		// {"utf-32", "", []rune{}},
		// {"utf-32be", "", []rune{}},
		// {"utf-32le", "", []rune{}},
		{"ascii-63", "", []rune{'␉', '0', 'A'}},
		{"ascii-65", "", []rune{'␉', '0', 'A', 'a', '{'}},
		{"ascii-67", "", []rune{'␉', '0', 'A', 'a', '~'}},
	}
	// test the list tables command
	args := []string{"tables"}
	gotB, err := listT.tester(args)
	if err != nil {
		t.Error(err)
		return
	}
	if len(gotB) == 0 {
		t.Errorf("TABLES, the result of the llist tables command is empty, or 0 bytes")
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
			if err != nil {
				t.Error(err)
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
