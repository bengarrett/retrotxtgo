package logs

import (
	"fmt"
	"os"
	"strings"

	"retrotxt.com/retrotxt/lib/str"
)

func CmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the command %s does not exist, %s", alert, name, err)
}

func CmdProblemFatal(name, flag string, err error) {
	fmt.Println(FlagProblem(name, flag, err))
	os.Exit(1)
}

// rename to Fatal
func ErrorFatal(err error) {
	fmt.Println(Errorf(err))
	os.Exit(1)
}

func Errorf(err error) string {
	return fmt.Sprintf("%s%s", str.Alert(), err) // TODO: change to string return?
}

func FlagProblem(name, flag string, err error) string {
	alert, toggle := str.Alert(), "--"
	fmt.Println("FLAG:", flag)
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s", alert, name, toggle, flag, err)
}

func MarkProblem(value string, new, err error) string {
	v := value
	a := str.Alert()
	n := str.Cf(fmt.Sprintf("%v", new))
	e := str.Cf(fmt.Sprintf("%v", err))

	return fmt.Sprintf("%s %s %q: %s", a, n, v, e)
}

func MarkProblemFatal(value string, new, err error) {
	fmt.Println(MarkProblem(value, new, err))
	os.Exit(1)
}

// ok
func Problemln(new, err error) {
	e := fmt.Errorf("%s: %w", new, err)
	fmt.Printf("%s%s\n", str.Alert(), e)
}

// ok
func ProblemFatal(new, err error) {
	Problemln(new, err)
	os.Exit(1)
}

func SubCmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the subcommand %s does not exist, %s", alert, name, err)
}

func Hint(s string, err error) string {
	return fmt.Sprintf("%s\n         run %s", Errorf(err), str.Example("retrotxt "+s))
}
