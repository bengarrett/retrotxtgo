package create

import (
	"fmt"
	"testing"
)

func ExampleServe() {
	// Initalize the bare minimum configuration.
	b := []byte("hello world")
	args := Args{}
	args.Layout = "standard"
	args.Port = 8080

	// The test argument will immediately shutdown
	// the server after it successfully starts.
	args.test = true

	// Run the HTTP server
	i, a, err := args.Serve(&b)
	if err != nil {
		fmt.Print(i, a, err)
	}
	// Output:Server example was successful
}

func TestPort(t *testing.T) {
	tests := []struct {
		name string
		port uint
		want bool
	}{
		{"empty", 0, true},
		{"www", 80, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Port(tt.port); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
