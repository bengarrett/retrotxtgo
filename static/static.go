//nolint:golint,gochecknoglobals
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

// CSS

//go:embed css/styles.css
var CSSStyles []byte

//go:embed css/text_bbs.css
var CSSBBS []byte

//go:embed css/text_blink.css
var CSSBlink []byte

//go:embed css/text_pcboard.css
var CSSPCBoard []byte
