package sample

import "errors"

var (
	ErrConvert = errors.New("unknown convert method")
	ErrConvNil = errors.New("conv argument cannot be empty")
)
