// nolint:golint,gochecknoglobals,stylecheck
package static

import (
	"embed"
	"errors"
)

//ErrPackGet = errors.New("pack.get name is invalid")
var ErrNotFound = errors.New("static file was not found")

//go:embed font/*
var Font embed.FS

//go:embed img/*
var Image embed.FS

//go:embed text/*
var Text embed.FS

//go:embed js/scripts.js
var Scripts []byte

//go:embed css/styles.css
var Styles []byte
