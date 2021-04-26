package create

import "errors"

var (
	ErrName      = errors.New("font name is not known")
	ErrPack      = errors.New("font pack is not found")
	ErrEmptyName = errors.New("filename is empty")
	ErrFileExist = errors.New("filename already exists")
	ErrReqOW     = errors.New("include an -o flag to overwrite")
	ErrUnknownFF = errors.New("unknown font family")
	ErrNilByte   = errors.New("cannot convert a nil byte value")
	ErrTmplDir   = errors.New("the path to the template file is a directory")
	ErrNoLayout  = errors.New("layout template does not exist")
	ErrLayout    = errors.New("unknown layout template")
	ErrTmpDir    = errors.New("temporary directory match")
	ErrRmTmpDir  = errors.New("temporary directory removal")
	ErrPort      = errors.New("tried and failed to serve using these ports")
)
