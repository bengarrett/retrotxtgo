package config

import "errors"

var (
	ErrCfgName = errors.New("unknown setting name")

	ErrCFG     = errors.New("unknown configuration name")
	ErrEnv     = errors.New("set one by creating an $EDITOR environment variable in your shell configuration")
	ErrKey     = errors.New("configuration key is invalid")
	ErrNoName  = errors.New("name cannot be empty")
	ErrNoFName = errors.New("filename cannot be empty")
	ErrSetting = errors.New("configuration setting name is not known")
	// key types
	ErrBool   = errors.New("key is not a boolean (true/false) value")
	ErrString = errors.New("key is not a string (text) value")
	ErrUint   = errors.New("key is not a absolute number")
)
