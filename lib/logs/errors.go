package logs

import (
	"errors"
)

var (
	// logs package errors.

	ErrCmdExist   = errors.New("the command is invalid")
	ErrFlag       = errors.New("the flag does not work with this command")
	ErrFlagChoice = errors.New("choose a value from the following")
	ErrFlagNil    = errors.New("the flag with a value must be included with this command")
	ErrNil        = errors.New("error value cannot be nil")
	ErrNotBool    = errors.New("the value must be either true or false")
	ErrNotInt     = errors.New("the value must be a number")
	ErrNotInts    = errors.New("the value must be a single or a list of numbers")
	ErrNotNil     = errors.New("the value cannot be empty")
	ErrLogSave    = errors.New("save fatal logs failed")
)

var (
	// config file.

	ErrCfg       = errors.New("could not parse or use the configuration file")
	ErrCfgCreate = errors.New("could not create a configuration file")
	ErrCfgFile   = errors.New("could not open the configuration file")
	ErrCfgName   = errors.New("unknown configuration or setting name")

	// inputs.

	ErrNameNil = errors.New("name cannot be empty")

	// pipe/stdin.

	ErrPipe      = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse = errors.New("could not parse the text stream from piped stdin (standard input)")

	// named file.

	ErrFileOpen = errors.New("could not open the file")
	ErrFileNil  = errors.New("file does not exist")
	ErrFileSave = errors.New("could not save the file")
	ErrDirSave  = errors.New("could not save the file to the directory")

	// sample file.

	ErrSampHTML = errors.New("could not convert the sample text into a HTML document")
	ErrSampFile = errors.New("unknown sample filename")
	ErrSampView = errors.New("could not view the sample text")

	// template file.

	ErrTmplDir = errors.New("the named template file is a directory")
	ErrTmplNil = errors.New("the named template layout does not exist")

	// temporary directory.

	ErrTmpClean = errors.New("could not cleanup the temporary directory")
	ErrTmpDir   = errors.New("could not save file to the temporary directory")

	// temporary file.

	ErrTmpClose = errors.New("could not close temporary file")
	ErrTmpOpen  = errors.New("could not create temporary file")
	ErrTmpSave  = errors.New("could not save to temporary file")

	// generic errors.

	ErrEncode    = errors.New("text encoding not known or supported")
	ErrFmt       = errors.New("format is not known")
	ErrHighlight = errors.New("could not format or colorize the element")
	ErrMax       = errors.New("maximum attempts reached")
	ErrTabFlush  = errors.New("tab writer could not write or flush")
	ErrZipFile   = errors.New("could not create the zip archive")

	// logs errors.

	ErrEmpty = errors.New("value is empty")
	ErrShort = errors.New("word count is too short, less than 3")
)
