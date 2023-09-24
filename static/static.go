// Package static provides the embedded static files for
// text, ANSI and BBS used for testing and examples.
package static

import (
	"embed"
	"errors"
)

var ErrNotFound = errors.New("internal embed file is not found")

// File is the embedded file system with all the static files.
//
//go:embed *
var File embed.FS

// ANSI is the embedded file system with the ansi subdirectory.
//
//go:embed ansi/*
var ANSI embed.FS

// Text is the embedded file system with the text subdirectory.
//
//go:embed text/*
var Text embed.FS
