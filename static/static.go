// nolint:golint,gochecknoglobals
package static

import (
	"embed"
	"errors"
)

var ErrNotFound = errors.New("static file was not found")

//go:embed *
var File embed.FS

//go:embed font/*
var Font embed.FS

//go:embed img/*
var Image embed.FS

//go:embed html/*
var Tmpl embed.FS

//go:embed text/*
var Text embed.FS

//go:embed js/scripts.js
var Scripts []byte

//go:embed css/styles.css
var Styles []byte
