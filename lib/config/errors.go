package config

import "errors"

var (
	ErrEditor   = errors.New("no suitable text editor can be found")
	ErrFilename = errors.New("filename cannot be empty")
	ErrLogo     = errors.New("retrotxt logo is missing")
	ErrName     = errors.New("unknown configuration setting name")
	ErrNameNil  = errors.New("name cannot be empty")

	// key types
	ErrBool   = errors.New("key is not a boolean (true/false) value")
	ErrString = errors.New("key is not a string (text) value")
	ErrUint   = errors.New("key is not a absolute number")
)
