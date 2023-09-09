package static

import (
	"embed"
	"errors"
)

var ErrNotFound = errors.New("internal embed file is not found")

//go:embed *
var File embed.FS

//go:embed ansi/*
var ANSI embed.FS

//go:embed text/*
var Text embed.FS
