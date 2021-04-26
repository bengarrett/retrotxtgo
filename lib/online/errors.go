package online

import "errors"

var (
	ErrJSON = errors.New("cannot understand the response body as the syntax is not json")
	ErrMash = errors.New("cannot unmarshal the json response body")
)
