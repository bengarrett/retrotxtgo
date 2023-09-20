package logs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

var ErrTest = errors.New("something went wrong")

func ExampleHint() {
	err := errors.New("oops")
	fmt.Print(logs.Hint(err, "helpme"))
	// Output: Problem:
	// oops.
	//  run retrotxt helpme
}

func ExampleSprint() {
	err := errors.New("oops")
	fmt.Print(logs.Sprint(err))
	// Output: Problem:
	// oops.
}

func ExampleSprintCmd() {
	err := errors.New("oops")
	fmt.Print(logs.SprintCmd(err, "helpme"))
	// Output: Problem:
	//  the command helpme does not exist, oops
}

func ExampleSprintFlag() {
	err := errors.New("oops")
	fmt.Print(logs.SprintFlag(err, "error", "err"))
	// Output: Problem:
	//  with the error --err flag, oops
}

func ExampleSprintS() {
	err := errors.New("oops")
	wrap := errors.New("uh-oh")
	fmt.Print(logs.SprintS(err, wrap, "we have some errors"))
	// Output: Problem:
	//  oops "we have some errors": uh-oh
}

func TestHint_String(t *testing.T) {
	t.Parallel()
	color.Enable = false
	err := logs.Hint(ErrTest, "hint")
	assert.Contains(t, err, "Problem:")
	assert.Contains(t, err, "something went wrong")
	err = logs.Hint(nil, "hint")
	assert.Empty(t, err)
}

func TestSprint(t *testing.T) {
	t.Parallel()
	color.Enable = false
	err := logs.Sprint(ErrTest)
	assert.Contains(t, err, "Problem:")
	assert.Contains(t, err, "something went wrong")
	err = logs.Sprint(nil)
	assert.Empty(t, err)
}

func TestSprintCmd(t *testing.T) {
	t.Parallel()
	color.Enable = false
	err := logs.SprintCmd(ErrTest, "hint")
	assert.Contains(t, err, "Problem:")
	assert.Contains(t, err, "something went wrong")
	err = logs.SprintCmd(nil, "hint")
	assert.Empty(t, err)
}

func TestSprintFlag(t *testing.T) {
	t.Parallel()
	color.Enable = false
	err := logs.SprintFlag(ErrTest, "hint", "dosomething")
	assert.Contains(t, err, "Problem:")
	assert.Contains(t, err, "with the hint --dosomething flag, something went wrong")
	err = logs.SprintFlag(nil, "hint", "dosomething")
	assert.Empty(t, err)
}

func TestSprintS(t *testing.T) {
	t.Parallel()
	color.Enable = false
	err := logs.SprintS(ErrTest, ErrTest, "hint")
	assert.Contains(t, err, "Problem:")
	assert.Contains(t, err, "something went wrong")
	err = logs.SprintS(nil, nil, "hint")
	assert.Empty(t, err)
}
