package convert

import "errors"

var (
	ErrChainANSI = errors.New("ansi() is a chain method that is to be used in conjunction with swap: c.swap().ansi()")
	ErrChainWrap = errors.New("wrapWidth() is a chain method that is to be used in conjunction with swap: c.swap().wrapWidth()")
	ErrBytes     = errors.New("cannot transform an empty byte slice")
	ErrEncoding  = errors.New("no encoding provided")
	ErrName      = errors.New("encoding cannot match name or alias")
	ErrUTF8      = errors.New("string cannot encode to utf-8")
	ErrUTF16     = errors.New("utf-16 table encodings are not supported")
	ErrUTF32     = errors.New("utf-32 table encodings are not supported")
	ErrWidth     = errors.New("cannot determine the number columns from using line break")
)
