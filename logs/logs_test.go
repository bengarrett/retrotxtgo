package logs_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/logs"
	"github.com/gookit/color"
	"github.com/nalgeon/be"
)

var ErrTest = errors.New("something went wrong")

func init() {
	color.Enable = false
}

func ExampleHint() {
	err := errors.New("oops") //nolint:goerr113
	fmt.Print(logs.Hint(err, "helpme"))
	// Output: Problem:
	// oops.
	//  run retrotxt helpme
}

func ExampleSprint() {
	err := errors.New("oops") //nolint:goerr113
	fmt.Print(logs.Sprint(err))
	// Output: Problem:
	// oops.
}

func ExampleSprintCmd() {
	err := errors.New("oops") //nolint:goerr113
	fmt.Print(logs.SprintCmd(err, "helpme"))
	// Output: Problem:
	//  the command helpme does not exist, oops
}

func ExampleSprintFlag() {
	err := errors.New("oops") //nolint:goerr113
	fmt.Print(logs.SprintFlag(err, "error", "err"))
	// Output: Problem:
	//  with the error --err flag, oops
}

func ExampleSprintS() {
	err := errors.New("oops")   //nolint:goerr113
	wrap := errors.New("uh-oh") //nolint:goerr113
	fmt.Print(logs.SprintS(err, wrap, "we have some errors"))
	// Output: Problem:
	//  oops "we have some errors": uh-oh
}

func TestHint_String(t *testing.T) {
	t.Parallel()
	err := logs.Hint(ErrTest, "hint")
	find := strings.Contains(err, "Problem:")
	be.True(t, find)
	find = strings.Contains(err, "something went wrong")
	be.True(t, find)
	err = logs.Hint(nil, "hint")
	be.True(t, err == "")
}

func TestSprint(t *testing.T) {
	t.Parallel()
	err := logs.Sprint(ErrTest)
	find := strings.Contains(err, "Problem:")
	be.True(t, find)
	find = strings.Contains(err, "something went wrong")
	be.True(t, find)
	err = logs.Sprint(nil)
	be.True(t, err == "")
}

func TestSprintCmd(t *testing.T) {
	t.Parallel()
	err1 := logs.SprintCmd(ErrTest, "hint")
	find := strings.Contains(err1, "Problem:")
	be.True(t, find)
	find = strings.Contains(err1, "something went wrong")
	be.True(t, find)
	err2 := logs.SprintCmd(nil, "hint")
	be.True(t, err2 == "")
}

func TestSprintFlag(t *testing.T) {
	t.Parallel()
	err := logs.SprintFlag(ErrTest, "hint", "dosomething")
	find := strings.Contains(err, "Problem:")
	be.True(t, find)
	find = strings.Contains(err, "with the hint --dosomething flag, something went wrong")
	be.True(t, find)
	err = logs.SprintFlag(nil, "hint", "dosomething")
	be.True(t, err == "")
}

func TestSprintS(t *testing.T) {
	t.Parallel()
	err := logs.SprintS(ErrTest, ErrTest, "hint")
	find := strings.Contains(err, "Problem:")
	be.True(t, find)
	find = strings.Contains(err, "something went wrong")
	be.True(t, find)
	err = logs.SprintS(nil, nil, "hint")
	be.True(t, err == "")
}

func TestInvalid(t *testing.T) {
	t.Parallel()
	s := logs.Invalid(nil, "", "")
	be.Equal(t, s, "")
	s = logs.Invalid(nil, "", "", "param1")
	be.Equal(t, s, "")
}
