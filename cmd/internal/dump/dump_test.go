package dump_test

// The dump package has complex dependencies that make it difficult to test
// without extensive mocking. For now, we'll just test the basic error handling.
// This test is skipped to avoid panics from missing command configuration.

import (
	"errors"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/dump"
	"github.com/nalgeon/be"
)

func TestErrPipeRead(t *testing.T) {
	t.Parallel()
	be.Equal(t, dump.ErrPipeRead.Error(), "could not read text stream from piped stdin (standard input)")
}

func TestPipe(t *testing.T) {
	t.Parallel()
	// NOTE: This will try to read from stdin and fail, but it should not panic
	err := dump.Pipe(nil)
	// should fail because there's no actual pipe
	be.True(t, err != nil)
	// error should be related to pipe reading
	be.True(t, errors.Is(err, dump.ErrPipeRead))
}
