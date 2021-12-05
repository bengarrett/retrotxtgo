package config

import "errors"

var (
	ErrEditorNil = errors.New("no suitable text editor can be found")
	ErrEditorRun = errors.New("editor cannot be run")
	ErrLogo      = errors.New("program logo is missing")
	ErrSaveType  = errors.New("save value type is unsupported")
)
