package logs

import (
	"errors"
)

var (
	// config file
	ErrCfgCreate = errors.New("could not create a configuration file")
	ErrCfgFile   = errors.New("could not open the configuration file")

	// pipe/stdin
	ErrPipe      = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse = errors.New("could not parse the text stream from piped stdin (standard input)")

	// named file
	ErrFileOpen = errors.New("could not open the file")
	ErrFileSave = errors.New("could not save the file")
	ErrDirSave  = errors.New("could not save the file to the directory")

	// sample file
	ErrSampHTML = errors.New("could not convert the sample text into a HTML document")
	ErrSampFile = errors.New("unknown sample filename")
	ErrSampView = errors.New("could not view the sample text")

	// temporary directory
	ErrTmpClean = errors.New("could not cleanup the temporary directory")
	ErrTmpDir   = errors.New("could not save file to the temporary directory")

	// temporary file
	ErrTmpClose = errors.New("could not close temporary file")
	ErrTmpOpen  = errors.New("could not create temporary file")
	ErrTmpSave  = errors.New("could not save to temporary file")

	// generic errors
	ErrEncode    = errors.New("text encoding not known or supported")
	ErrHighlight = errors.New("could not format or colorize the element")
	ErrTabFlush  = errors.New("tab writer could not write or flush")
	ErrZipFile   = errors.New("could not create the zip archive")

	// logs errors
	ErrEmpty = errors.New("value is empty")
	ErrShort = errors.New("word count is too short, less than 3")
)
