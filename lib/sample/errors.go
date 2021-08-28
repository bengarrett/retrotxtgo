package sample

import "errors"

var (
	ErrEncode  = errors.New("no encoding provided")
	ErrConvert = errors.New("unknown convert method")
	ErrConvNil = errors.New("conv argument cannot be empty")
)
