package cmd

import "errors"

var (
	ErrHelp        = errors.New("command help could not display")
	ErrHideCreate  = errors.New("could not hide the create flag")
	ErrMarkRequire = errors.New("command flag could not be marked as required")
	ErrShell       = errors.New("could not generate shell completion")
	ErrTable       = errors.New("could not display the table")
	ErrUsage       = errors.New("command usage could not display")

	ErrCreate = errors.New("could not convert the text into a HTML document")
	ErrEncode = errors.New("could not convert the text into the requested encoding")
	ErrInfo   = errors.New("could not any obtain information")
	ErrIANA   = errors.New("could not work out the IANA index or MIME type")
	ErrUTF8   = errors.New("could not convert the text into a UTF-8 encoding")
)
