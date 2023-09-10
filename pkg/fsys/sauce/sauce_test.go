package sauce_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/fsys/sauce"
)

func TestSAUCE(t *testing.T) {
	name, err := filepath.Abs("../../static/text/sauce.txt")
	if err != nil {
		t.Error(err)
		return
	}
	f, err := fsys.ReadAllBytes(name)
	if err != nil {
		t.Error(err)
		return
	}
	got := sauce.SAUCE{}
	got.Read(&f)
	if reflect.DeepEqual(got, sauce.SAUCE{}) {
		t.Error("SAUCE result is empty")
		return
	}
	if !got.Use {
		t.Error("SAUCE.Use result is false")
	}
	const wantTitle = "Sauce title"
	if got.Title != wantTitle {
		t.Errorf("SAUCE.Title = %q, want %q", got.Title, wantTitle)
	}
	const wantAuthor = "Sauce author"
	if got.Author != wantAuthor {
		t.Errorf("SAUCE.Author = %q, want %q", got.Title, wantAuthor)
	}
	const wantGroup = "Sauce group"
	if got.Group != wantGroup {
		t.Errorf("SAUCE.Group = %q, want %q", got.Group, wantGroup)
	}
	const wantDesc = "ASCII text file with no formatting codes or color codes."
	if got.Description != wantDesc {
		t.Errorf("SAUCE.Description = %q, want %q", got.Description, wantDesc)
	}
	const wantWidth = 977
	if got.Width != wantWidth {
		t.Errorf("SAUCE.Width = %d, want %d", got.Width, wantWidth)
	}
	const wantLines = 9
	if got.Lines != wantLines {
		t.Errorf("SAUCE.Lines = %d, want %d", got.Lines, wantLines)
	}
}
