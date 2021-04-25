package logs

import (
	"errors"
)

var (
	// generic errors
	ErrConfigFile = errors.New("could not open the configuration file")
	ErrEncode     = errors.New("text encoding not known or supported")
	ErrHighlight  = errors.New("could not format or colorize the element")
	ErrOpenFile   = errors.New("could not open the file")
	ErrPipe       = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse  = errors.New("could not parse the text stream from piped stdin (standard input)")
	ErrSaveDir    = errors.New("could not save file to the directory")
	ErrSaveFile   = errors.New("could not save the file")
	ErrSampFile   = errors.New("unknown sample filename")
	ErrTmpClean   = errors.New("could not cleanup the temporary directory")
	ErrTmpDir     = errors.New("could not save file to the temporary directory")
	ErrTmpSave    = errors.New("could not save the temporary file")
	ErrTabFlush   = errors.New("tab writer could not write or flush")
	ErrZipFile    = errors.New("could not create the zip archive")
	// command (cmd library) and argument errors
	ErrHelp        = errors.New("command help could not display")
	ErrMarkRequire = errors.New("command flag could not be marked as required")
	ErrUsage       = errors.New("command usage could not display")
	// type errors
	ErrCmd     = errors.New("choose a command from the list available")
	ErrNewCmd  = errors.New("choose another command from the available commands")
	ErrNoCmd   = errors.New("invalid command")
	ErrEmpty   = errors.New("value is empty")
	ErrFlag    = errors.New("use a flag from the list of flags")
	ErrSyntax  = errors.New("flags can only be in -s (short) or --long (long) form")
	ErrNoFlagx = errors.New("cannot be empty and requires a value")
	ErrReqFlag = errors.New("you must include this flag in your command")
	ErrSlice   = errors.New("invalid option choice")
	ErrShort   = errors.New("word count is too short, less than 3")
	ErrVal     = errors.New("value is not a valid choice")
)
