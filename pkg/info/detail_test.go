package info_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/info"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
)

func ExampleDetail_Ctrls() {
	d := info.Detail{}
	_ = d.Ctrls("testdata/example.ans")
	fmt.Printf("%d controls\n", d.Count.Controls)
	// Output: 12 controls
}

func ExampleDetail_Marshal() {
	d := info.Detail{
		Name: "example.ans",
	}
	b := bytes.Buffer{}
	_ = d.Marshal(&b, info.JSON)
	fmt.Printf("%d bytes, is json = %t", b.Len(), json.Valid(b.Bytes()))
	// Output: 2130 bytes, is json = true
}

func sampleFile() string {
	b := []byte(mock.T()["Tabs"]) // Tabs and Unicode glyphs
	const digits = 40
	i := rand.Intn(digits)
	name := fmt.Sprintf("info_test_%d.txt", i)
	path, err := fsys.SaveTemp(name, b...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func TestValidText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"empty", "", false},
		{"image", "image/jpeg", false},
		{"stream", "application/octet-stream", true},
		{"text", "text/plain", true},
		{"js", "text/javascript", true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := info.ValidText(tt.contentType); got != tt.want {
				t.Errorf("ValidText() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestRead(t *testing.T) {
	t.Parallel()
	tmp := sampleFile()
	t.Cleanup(func() {
		os.Remove(tmp)
	})
	fmt.Fprintln(os.Stdout, "path:", tmp)
	var got info.Detail
	if err := got.Read(tmp); err != nil {
		t.Errorf("Read() = %v, want %v", err, nil)
	}
	if got.Size.Bytes != 57 {
		t.Errorf("Read() = %v, want %v", got.Size.Bytes, 57)
	}
	if !strings.Contains(got.Name, "info_test") {
		t.Errorf("Read() = %v, want %v", got.Name, "info_test.txt")
	}
	if !strings.Contains(got.Slug, "info-test") {
		t.Errorf("Read() = %v, want %v", got.Slug, "info-test-txt")
	}
	if got.Mime.Type != "text/plain" {
		t.Errorf("Read() = %v, want %v", got.Mime, "text/plain")
	}
	if got.UTF8 != true {
		t.Errorf("Read() = %v, want %v", got.UTF8, true)
	}
	const want = "883643f5e9ed278732c92d9b6f834b96"
	if got.Sums.MD5 != want {
		t.Errorf("Read() = %v, want %v", got.Sums.MD5, want)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"empty", args{[]byte("")}, 0, false},
		{"string", args{[]byte("hello")}, 5, false},
		{"string", args{[]byte("世界你好")}, 4, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			var got info.Detail
			err := got.Parse("", tt.args.data...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Count.Chars, tt.want) {
				t.Errorf("Parse() = %v, want %v", got.Count.Chars, tt.want)
			}
		}
	})
}

func TestMarshal_json(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		d      info.Detail
		format info.Format
		want   bool
	}{
		{"no indent", info.Detail{}, info.JSONMin, true},
		{"indent", info.Detail{}, info.JSON, true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			j := &bytes.Buffer{}
			_ = tt.d.Marshal(j, tt.format)
			if got := json.Valid(j.Bytes()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() json = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestMarshal_text(t *testing.T) {
	t.Parallel()
	tmp := sampleFile()
	t.Cleanup(func() {
		os.Remove(tmp)
	})
	const want = 835
	var d info.Detail
	if err := d.Read(tmp); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s := &strings.Builder{}
	_ = d.Marshal(s, info.PlainText)
	if got := len(s.String()); got != want {
		t.Errorf("Marshal() text = %v, want %v", got, want)
	}
}
