package cmd

import "errors"

var (
	ErrHelp        = errors.New("command help could not display")
	ErrHide        = errors.New("could not hide the flag")
	ErrMarkRequire = errors.New("command flag could not be marked as required")
	ErrTable       = errors.New("could not display the table")
	ErrUsage       = errors.New("command usage could not display")

	ErrCreate = errors.New("could not convert the text into a HTML document")
	ErrEncode = errors.New("could not convert the text into the requested encoding")
	ErrInfo   = errors.New("could not any obtain information")
	ErrIANA   = errors.New("could not work out the IANA index or MIME type")
	ErrUTF8   = errors.New("could not convert the text into a UTF-8 encoding")

	ErrBody    = errors.New("could not parse the body flag")
	ErrServeIn = errors.New("could not serve stdin over HTTP")

	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	ErrCacheData = errors.New("set cache cannot create a data path")
	ErrCacheSave = errors.New("set cache cannot save data")

	ErrFlagE     = errors.New("ignoring encode flag")
	ErrFilenames = errors.New("ignoring [filenames]")
)
