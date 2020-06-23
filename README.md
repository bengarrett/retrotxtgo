![Build](https://github.com/bengarrett/retrotxtgo/workflows/Go/badge.svg) ![Tests](https://github.com/bengarrett/retrotxtgo/workflows/Go%20tests/badge.svg)

# RetroTxt on Go

[RetroTxt](https://github.com/bengarrett/retrotxt) for the command line

### _α_ work-in-progress, feature incomplete & is not in a usable state

---

## Example

```sh
cat samples/ascii-logos.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

```sh
retrotxt create ascii-logos.txt # creates an index.html file in the home directory
```

```html
cat ~/index.html

<!DOCTYPE html>

<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>RetroTXT | ASCII logos</title>
    <meta name="description" content="The RetroTxt logo" />
    <meta name="author" content="Ben Garrett" />
    <meta name="keywords" content="retrotxt,ansi,ascii" />
    <meta
      name="generator"
      content="RetroTxt v; 0001-01-01 00:00:00 &#43;0000 UTC"
    />
    <link rel="stylesheet" href="static/css/styles.css" />
    <script src="static/js/scripts.js" defer></script>
  </head>

  <body>
    <main>
      <pre>
██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝</pre
      >
    </main>
  </body>
</html>
```

---

# Features

- [ ] Convert ASCII text to HTML.
- [ ] Convert ANSI text to HTML.
- [ ] Convert BBS text to HTML.
- [x] List or export (json, text, xml) meta details of a text file.
- [ ] List or export SAUCE metadata of a file.
- [ ] Transform the encoding of a text file. CP437 -> UTF8, UTF8 -> ISO8859-1 ...
- [ ] View any legacy encoded text file in a UTF8 terminal by converting on the fly.
- [x] Extensive customisations through command flags or a configuration file with a setup.
- [ ] ANSI compatibility tests, output 16, 88, 256, high and true-color tables.
- [x] Multiplatform support including Windows, macOS, Linux, Raspbian and FreeBSD.

---

## Install

There are [downloads](https://github.com/bengarrett/retrotxtgo/releases/latest/) available for
[Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_windows.zip),
[macOS](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macos.zip),
[Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.tar.gz),
[FreeBSD](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_freebsd.tar.gz) as well as for the
[Raspberry Pi/ARM family](https://github.com/bengarrett/retrotxtgo/releases/latest/).

Otherwise these operating system specific install methods are available.

### Windows

#### ~~[Chocolatey](https://chocolatey.org/)~~ \*

```ps
choco install retrotxt
retrotxt version
```

#### [Scoop](https://scoop.sh/)

```ps
scoop bucket add retrotxt https://github.com/bengarrett/retrotxtgo.git
scoop install retrotxt
retrotxt version
```

### macOS

#### ~~[Homebrew](https://brew.sh/)~~ \*

```sh
brew cask install retrotxt
retrotxt
```

### Linux

#### ~~[Linux Snap](https://snapcraft.io/)~~ \*

```sh
snap install retrotxt
retrotxt version
```

#### Raspberry Pi OS / Raspbian

Download the **deb** package for either the
[Raspberry Pi](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_raspberry_pi.deb) <small>(ARMv7)</small>
or the [Zero family](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_raspberry_pi-zero.deb) <small>(ARMv6)</small>
and install using `dpkg -i`.

```sh
# Raspberry Pi example
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_raspberry_pi.deb
dpkg -i retrotxt_raspberry_pi.deb
retrotxt version
```

#### Deb - Ubuntu, Mint, Debian

```sh
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.deb
dpkg -i retrotxt_linux.deb
retrotxt version
```

#### RPM - Fedora, OpenSUSE, CentOS, RHEL

```sh
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.rpm
rpm -i retrotxt_linux.rpm
retrotxt version
```

\* not implemented

## Build using Go

RetroTxt on Go requires [Go v1.13+](https://github.com/golang/go/wiki/MinimumRequirements).
Assuming [Go](https://golang.org/) and and the relevant build-tools are already installed.

Standard operating systems.

```sh
git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
# recommended, run the package and import tests on your distribution/platform
go test ./...
go build -i -o $GOBIN/retrotxt
retrotxt version
```

Windows PowerShell

```powershell
git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
go build -i -o $Env:GOBIN\retrotxt.exe
retrotxt version
```

The binary will be installed either at `$GOPATH/bin` or `$HOME/go/bin`

## Compile to other platforms using Go

Go supports a maylard of operating systems and platforms that can be built using any other supported platform.
So if you were on Linux and needed to compile a 32-bit version of RetroTxt to target Windows XP.

```sh
# to see a list of supported platforms
go tool dist list

git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
env GOOS=windows GOARCH=386 go build -o retrotxt.exe
# optional, compress the binary
upx retrotxt.exe
# test the binary
file retrotxt.exe
```

---

### Why Go?

- native [Unicode](https://golang.org/pkg/unicode/) and UTF8/16/32 support
- [native legacy text encodings support](golang.org/x/text/encoding/charmap)
- creates a standalone binary with no dependencies
- [wide OS and CPU support](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)
- simple and fast compiling
- I know it 😉

---

### From here onwards the following text are developer notes

### Shrink the binary

```sh
# strip the DWARF debugging information
go build -ldflags '-w'

# strip the symbol table
go build -ldflags '-d'

go build --race .
# on Ubuntu in Go v1.14 results in a 29MB binary

go build .
# on Ubuntu in Go v1.14 results in a 21MB binary

go build -ldflags '-w -d' .
# on Ubuntu in Go v1.14 results in a 17MB binary
```

[UPX compression](https://upx.github.io/)

```sh
# upx 3.94 on ubuntu 18.04
upx retrotxtgo # 17 -> 6.5MB file
upx --best retrotxtgo # (slow) 17MB --> 6.45MB file (not worth the saving)
upx --brute retrotxtgo # (very slow) 17MB --> 4.8MB file
```

#### References

- [SAUCE](http://www.acid.org/info/sauce/sauce.htm)
- packages: [Cobra](https://pkg.go.dev/github.com/spf13/cobra)/[Viper](https://pkg.go.dev/mod/github.com/spf13/viper)
- go pkg: [utf8](https://golang.org/pkg/unicode/utf8/)/[unicode](https://golang.org/pkg/unicode/),
  [x-text](https://pkg.go.dev/golang.org/x/text@v0.3.2?tab=subdirectories), [x-charmap](https://pkg.go.dev/golang.org/x/text@v0.3.2/encoding/charmap?tab=doc), [x-encoding](https://pkg.go.dev/golang.org/x/text@v0.3.2/encoding?tab=doc)

### Go libraries

- [Package xstrings: A collection of useful string functions in Go](https://github.com/huandu/xstrings)
  includes center, capitalize, justify, reverse text.

- [A collection of common regular expressions for Go](https://github.com/mingrammer/commonregex)
  use for dynamically hyperlinking emails, links. wrap hash values around `<code>` tags. parse known ports.

- [Devd - A local webserver for developers](https://github.com/cortesi/devd)
  replacement for the current internal webserver?

---

### SAUCE implementation

`fmt.Sprintf()`

#### Integer

- %b base 2 (binary)
- %c unicode code point
- %d base 10
- %o base 8
- %0 base 8 with 0o prefix
- %x base 16 lower-case %X upper-case
- %U unicode format

```go
\+ always print a sign (%+q)
\- pad with spaces to the right rather than the left
\# alternate formats:
add leading 0b for binary (%#b)
add leading 0 for octal (%#o)
add 0x or 0X for hex (%#x)
suppress 0x for %p (%#p)
always print a decimal point for floats (%#f)...
0 pad with leading zeros rather than spaces (%0d)...
```

Fprint, Fprintln... `Fprintf(w io.Writer, format, ...inferface{})`
Formats according to a format specifier and writes to w.

`constant.BitLen(x Value) int`
... BoolVal(x) Int64Val(x) Uint64Val(x)
... StringVal(x) string

`constant.Val(x) interface{}`

```go
b := constant.Make(false)
fmt.Printf("%v\n", constant.Val(b))
```

`Shift(x Value, op, s uint)`
op token list > https://golang.org/pkg/go/token/#Token

---

### Future CLI commands

- [ ] `html` - rename `create`: to encode text to HTML. `retrotxt html -n sometext.txt`
- [ ] `save` - to encode text to another format. `retrotxt save -n sometext.txt -e cp437`
- [ ] add piping support - `cat sometext.txt | retrotxt html`
- - `retrotxt save -sometext.txt # save to utf8`
- [ ] add optional argument for destination (dir or file) that overrides the dir configuration.
- - `retrotxt create -n somefile.txt . # would create index.html in the current directory`
- - `retrotxt create -n somefile.txt somefile.html # create sometfile.htm`
- [ ] add a flag to --export `create` command results to a tar or zip archive.

---

## Possible --create flags

- font choice (family)
- font size
- font format
- - base64
- - woff2

- code-page

- input (overwrite for internal use)
- - ascii, ansi, etc

- quiet boolean

---

### Embedding static files such as JS, CSS and fonts.

[go-bindata](https://github.com/go-bindata/go-bindata) Turn data file into go code.

Another option, store fonts as Base64 encoded text in .go files and export+save them as binary data when required?

`go generate` [create our own packer?](https://dev.to/koddr/the-easiest-way-to-embed-static-files-into-a-binary-file-in-your-golang-app-no-external-dependencies-43pc)

---

### TODOs - changes to the existing code

- [x] **Remove all string references to `retrotxtgo`.**
- [x] Directory settings, change `.` shortcut to always use the current working directory.
- [x] When fetching github release data using the `version` command.
      Use HTTP e-tags cache and save the values to reduce the bandwidth usage.
- [x] config command should support the global --config flag.
- [ ] config shell should have a `--append/source/or` flag to save shell auto-completion?
- [ ] scan for supported but current shell configuration.
- [x] when using `create` detect any out of range or unsafe unicode encoding and assume cp437.
- [ ] in `create` HTML insert a header comment with source file, command flags. encoding info (cp437 -> utf8).
- [ ] generator meta tag should have a working date.
- [ ] newline scanner to determine the maxWidth of the text.
- [ ] reverse scan of file looking for EOF, SAUCE00 & COMNTT.
- [ ] scan for unique color codes like 24-bit colors.
- [ ] implement ASCII font and CSS support.
- [ ] scan and linkify any http/s, ftp, mailto links in HTML.

```sh
// TODO: env
// Indicates which language, character set, and sort order to use for messages, datatype conversions, and datetime formats.
// 1. LC_NUMERIC="en_GB.UTF-8"
// 1. LC_TIME="en_GB.UTF-8"
// 2. LC_ALL=""
// 3. LANG=""
// 4. LANGUAGE=""
// 4. US
```

---

### Submission to distribution package managers

Windows

[Chocolatey](https://chocolatey.org/docs/createpackages)

macOS

**[homebrew](https://brew.sh/)** [casks allow bin submissions](https://github.com/Homebrew/homebrew-cask/blob/master/doc/cask_language_reference/readme.md)

Auto-update Linux

Snap [snapcraft](https://snapcraft.io/first-snap#go), flathub is not for terminal apps.

#### Other managers that require sponsorship

Ubuntu/Mint APT .DEB (http://packaging.ubuntu.com/html/packaging-new-software.html)
See 4.5 and 4.6

Fedora, CentOS/ RHEL
https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/index.html
https://rpm-packaging-guide.github.io/
