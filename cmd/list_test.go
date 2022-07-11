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

func Test_ListTable(t *testing.T) {
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
		{"arabic", "", []rune{'ي', 'ص'}},
		{"iso88596e", "", []rune{'ي', 'ص'}},
		{"iso88596i", "", []rune{'ي', 'ص'}},
		{"greek", "", []rune{'₯', 'Ξ'}},
		{"hebrew", "", []rune{'א', 'ת'}},
		{"iso88598e", "", []rune{'א', 'ת'}},
		{"iso88598i", "", []rune{'א', 'ת'}},
		{"latin5", "", []rune{'ğ', 'Ğ'}},
		{"latin6", "", []rune{'ą', 'ß'}},
		{"iso889511", "", []rune{'ข', '๛', '฿'}},
		{"iso889512", "", []rune{}},
		{"iso889513", "", []rune{'š', '¾'}},
		{"iso889514", "", []rune{'Ŵ', 'ẅ'}},
		{"iso889515", "", []rune{'€', 'Š', 'œ'}},
		{"iso889516", "", []rune{'đ', 'Œ'}},
		{"koi8r", "", []rune{'Ю', '╬'}},
		{"koi8u", "", []rune{'є', 'Ґ', '█'}},
		{"mac", "", []rune{'∞', '‰', '◊'}},
		{"windows-874", "", []rune{'€', '“', '”'}},
		{"windows-1250", "", []rune{'‡', '“', '”'}},
		{"windows-1251", "", []rune{'“', '”', 'Љ'}},
		{"windows-1252", "", []rune{'“', '”', 'ƒ'}},
		{"windows-1253", "", []rune{'“', '”', '…'}},
		{"windows-1254", "", []rune{'“', '”', 'Ğ'}},
		{"windows-1255", "", []rune{'“', '”', '₪'}},
		{"windows-1256", "", []rune{'“', '”', 'گ'}},
		{"windows-1257", "", []rune{'“', '”', '€'}},
		{"windows-1258", "", []rune{'“', '”', '†'}},
		{"shift_jis", "", []rune{'ﾗ', 'ｼ', 'ﾎ'}},
		{"utf-8", "", []rune{'␀', '␟', '␗'}},
		{"utf-16", "", []rune{}},
		{"utf-16be", "", []rune{}},
		{"utf-16le", "", []rune{}},
		{"utf-32", "", []rune{}},
		{"utf-32be", "", []rune{}},
		{"utf-32le", "", []rune{}},
		{"ascii-63", "", []rune{'␉', '0', 'A'}},
		{"ascii-65", "", []rune{'␉', '0', 'A', 'a', '{'}},
		{"ascii-67", "", []rune{'␉', '0', 'A', 'a', '~'}},
	}
	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			args := []string{"table", tt.table}
			gotB, err := listT.tester(args)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Contains(gotB, []byte(tt.wantHeader)) {
				t.Errorf("could not find %q header in the table", tt.wantHeader)
			}
			for i, runes := range tt.wantRunes {
				if !bytes.ContainsRune(gotB, runes) {
					t.Errorf("%d. result doesn't include the expected rune, %q got:\n%s", i, runes, gotB)
				}
			}
		})
	}
}
