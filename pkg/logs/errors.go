package logs

import (
	"errors"
)

var (
	// internal arguments.
	ErrEmpty    = errors.New("value is empty")
	ErrErrorNil = errors.New("error value cannot be nil")
	ErrNameNil  = errors.New("name value cannot be empty")
	ErrShort    = errors.New("word count is too short, less than 3")

	// generic errors.
	ErrEncode    = errors.New("text encoding not known or supported")
	ErrFmt       = errors.New("format is not known")
	ErrHighlight = errors.New("could not format or colorize the element")
	ErrMax       = errors.New("maximum attempts reached")
	ErrTabFlush  = errors.New("tab writer could not write or flush")
	ErrZipFile   = errors.New("could not create the zip archive")

	// logs package errors.
	ErrCmd        = errors.New("the command is invalid")
	ErrFlag       = errors.New("the flag does not work with this command")
	ErrFlagChoice = errors.New("choose a value from the following")
	ErrFlagNil    = errors.New("the flag with a value must be included with this command")
	ErrNotBool    = errors.New("the value must be either true or false")
	ErrNotInt     = errors.New("the value must be a number")
	ErrNotInts    = errors.New("the value must be a single or a list of numbers")
	ErrNotNil     = errors.New("the value cannot be empty")
	ErrSave       = errors.New("save to log file failure")

	// config file.
	ErrConfigName = errors.New("unknown configuration or setting name")
	ErrConfigNew  = errors.New("could not create a configuration file")
	ErrConfigOpen = errors.New("could not open the configuration file")
	ErrConfigRead = errors.New("could not parse or use the configuration file")

	// pipe/stdin.
	ErrPipeEmpty = errors.New("empty text stream from piped stdin (standard input)")
	ErrPipeRead  = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse = errors.New("could not parse the text stream from piped stdin (standard input)")

	// named file.
	ErrFileName  = errors.New("file does not exist")
	ErrFileOpen  = errors.New("could not open the file")
	ErrFileSave  = errors.New("could not save the file")
	ErrFileSaveD = errors.New("could not save the file to the directory")

	// sample file.
	ErrSampleName = errors.New("sample filename does not exist")
	ErrSampleOpen = errors.New("could not open the sample text")
	ErrSampleHTML = errors.New("could not convert the sample text to a HTML document")

	// template file.
	ErrTmplName  = errors.New("the named template layout does not exist")
	ErrTmplIsDir = errors.New("the named template file is a directory")

	// temporary directory.
	ErrTmpRMD   = errors.New("could not cleanup the temporary directory")
	ErrTmpSaveD = errors.New("could not save file to the temporary directory")

	// temporary file.
	ErrTmpClose = errors.New("could not close the temporary file")
	ErrTmpOpen  = errors.New("could not open the temporary file")
	ErrTmpSave  = errors.New("could not save to the temporary file")
)
