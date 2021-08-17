package config

import "errors"

var (
	ErrEditorNil = errors.New("no suitable text editor can be found")
	ErrEditorRun = errors.New("editor cannot be run")
	ErrLogo      = errors.New("program logo is missing")
	ErrBool      = errors.New("key is not a boolean (true/false) value")
	ErrString    = errors.New("key is not a string (text) value")
	ErrUint      = errors.New("key is not a absolute number")
)
