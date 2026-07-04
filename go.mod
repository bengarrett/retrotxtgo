module github.com/bengarrett/retrotxtgo

go 1.26.3

// When updating go version, you must also update
//  - .github/workflows/goreleaser.yml for GitHub
//  - Formula/retrotxt.rb version must also be updated to the current tag

require (
	github.com/gookit/color v1.6.1
	github.com/karrick/godirwalk v1.17.0
	github.com/mattn/go-isatty v0.0.22
	github.com/mozillazg/go-slugify v0.2.0
	github.com/mozillazg/go-unidecode v0.2.0 // indirect
	github.com/muesli/go-app-paths v0.2.2
	github.com/spf13/cobra v1.10.2
	github.com/zRedShift/mimemagic v1.2.0
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0
	golang.org/x/text v0.38.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require (
	github.com/bengarrett/bbs v1.0.7
	github.com/bengarrett/sauce v1.2.7
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/nalgeon/be v0.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/x/ansi v0.11.7 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.15 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.23 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/nilaway v0.0.0-20251021214447-34f56b8c16b9 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
)

tool go.uber.org/nilaway/cmd/nilaway
