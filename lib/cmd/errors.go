package cmd

import "errors"

var (
	ErrCfgCreate  = errors.New("could not create a configuration file")
	ErrCmdUsage   = errors.New("command usage could not display")
	ErrCreateHTML = errors.New("could not convert the text into a HTML document")
	ErrEncode     = errors.New("could not convert the text into the requested encoding")
	ErrInfo       = errors.New("could not any obtain information")
	ErrIANA       = errors.New("could not work out the IANA index or MIME type")
	ErrSampleHTML = errors.New("could not convert the sample text into a HTML document")
	ErrSampleView = errors.New("could not view the sample text")
	ErrTable      = errors.New("could not display the table")
	ErrViewUTF8   = errors.New("could not convert the text into a UTF-8 encoding")

	ErrShellCompletion = errors.New("could not generate completion for")
	ErrHideCreate      = errors.New("could not hide the create flag")
	ErrIntpr           = errors.New("the interpreter is not supported")
	ErrPackValue       = errors.New("unknown package convert value")

	ErrTempClose = errors.New("could not close temporary file")
	ErrTempOpen  = errors.New("could not create temporary file")
	ErrTempWrite = errors.New("could not write to temporary file")
)
