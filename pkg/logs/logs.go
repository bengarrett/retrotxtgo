// Package logs handles the printing and saving of errors.
package logs

import (
	"fmt"
	"log"
	"os"
)

const (
	Panic = false
)

// Fatal saves the error to the logfile and exits.
func Fatal(err error) {
	if err == nil {
		return
	}
	// print error
	switch Panic {
	case true:
		log.Printf("error type: %T\tmsg: %v\n", err, err)
		log.Panic(err)
	default:
		fmt.Fprintln(os.Stderr, Sprint(err))
		os.Exit(OSErrCode)
	}
}
