package create

import "errors"

var (
	ErrCleanPath = errors.New("cleanup temporary path match failed")
	ErrFileExist = errors.New("filename already exists")
	ErrFileNil   = errors.New("filename cannot be empty")
	ErrFont      = errors.New("unknown font name or family")
	ErrNilByte   = errors.New("cannot convert a nil byte value")
	ErrPorts     = errors.New("cannot run the http server using these ports")
)
