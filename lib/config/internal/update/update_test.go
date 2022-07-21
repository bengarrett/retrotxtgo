package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	example := filepath.Join(os.TempDir(), "test.cfg")
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{""}, true},
		{"okay", args{example}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if err := Config(os.Stdout, tt.args.name); (err != nil) != tt.wantErr {
			// 	t.Errorf("Config() error = %v, wantErr %v", err, tt.wantErr)
			// }
			// if tt.args.name == example {
			// 	defer os.Remove(example)
			// }
		})
	}
}
