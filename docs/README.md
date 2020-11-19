![Build](https://github.com/bengarrett/retrotxtgo/workflows/Go/badge.svg) ![Tests](https://github.com/bengarrett/retrotxtgo/workflows/Go%20tests/badge.svg) ![Lint](https://github.com/bengarrett/retrotxtgo/workflows/golangci-lint/badge.svg)

# RetroTxt on Go

[RetroTxt](https://github.com/bengarrett/retrotxt) for the command line

### _α_ - work-in-progress and feature incomplete

---

## About

#### Text files created in the pre Unicode days often fail to display on modern systems.

```sh
cat samples/ascii-logos.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

#### Use Retrotxt to print legacy encoded text on modern terminals.

```sh
retrotxt view ascii-logos.txt

██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```

#### Or save it to a file.

```sh
retrotxt view ascii-logos.txt > ascii-logos-utf8.txt
```

#### Then turn the text into a static website with accurate fonts and colours.

```sh
retrotxt create --layout=compact ascii-logos.txt
```

```html
<!DOCTYPE html>

<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>RetroTXT | ASCII logos</title>
    <link rel="stylesheet" href="styles.css" />
    <link rel="stylesheet" href="font.css" />
    <script src="scripts.js" defer></script>
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

#### Easily serve it over its own HTTP server.

```sh
retrotxt create -p0 ascii-logos.txt

Web server is available at http://localhost:8080
Press Ctrl+C to stop
```

- insert browser screenshot

#### Save the files as a ready to use webpage.

```sh
retrotxt create --save ascii-logs.txt

saving to /home/ben/scripts.js
saving to /home/ben/styles.css
saving to /home/ben/index.html
saving to /home/ben/ibm-vga8.woff2
saving to /home/ben/font.css
...
```

#### Inline all the assets into a single HTML file for easier sharing.

```sh
retrotxt create --layout=inline --font-embed
```

```html
<!DOCTYPE html>

<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>RetroTXT | ASCII logos</title>
    <style type="text/css">
      body{background-color:#000;display:flex;flex-dir...}
      @font-face{font-family:vga;src:url(data:application/font-woff2;charset=utf-8;base64,d09GMgA...)}
    </style>
    <script defer>
      ...
    </script>
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

- [x] Convert ASCII text to HTML.
- [ ] Convert ANSI text to HTML.
- [ ] Convert BBS text to HTML.
- [x] List or export (json, text, xml) meta details of a text file.
- [x] List or export SAUCE metadata of a file.
- [x] Transform the encoding of a text file. CP437 -> UTF8, UTF8 -> ISO8859-1 ...
- [x] View any legacy encoded text file in a UTF8 terminal by converting on the fly.
- [x] Extensive customisations through command flags or a configuration file with a setup.
- [ ] ANSI compatibility tests, output 16, 88, 256, high and true-color tables.
- [x] Multiplatform support including Windows, macOS, Linux, Raspbian and FreeBSD.
- [x] IO redirection with piping support.

---

## Install

There are [downloads](https://github.com/bengarrett/retrotxtgo/releases/latest/) available for
[Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_windows.zip),
[macOS (intel)](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macos.zip),
[Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.tar.gz),
[FreeBSD](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_freebsd.tar.gz) and for the
[Raspberry Pi](https://github.com/bengarrett/retrotxtgo/releases/latest/).

Otherwise these operating system specific install methods are available.

### Windows

<!-- #### ~~[Chocolatey](https://chocolatey.org/)~~ \*

```ps
choco install retrotxt
retrotxt version
``` -->

#### [Scoop](https://scoop.sh/)

```ps
scoop bucket add retrotxt https://github.com/bengarrett/retrotxtgo.git
scoop install retrotxt
retrotxt version
```

### macOS (intel)

#### [Homebrew](https://brew.sh/)

```sh
brew cask install bengarrett/tap/retrotxt
retrotxt version
```

### Linux

<!-- #### ~~[Linux Snap](https://snapcraft.io/)~~ \*

```sh
snap install retrotxt
retrotxt version
``` -->

#### Raspberry Pi, Linux ARM

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

#### Deb packager

Used by, but not limited to Ubuntu, Mint, Debian

```sh
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.deb
dpkg -i retrotxt_linux.deb
retrotxt version
```

#### RPM packager

Used by, but not limited to Fedora, OpenSUSE, CentOS, RHEL

```sh
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.rpm
rpm -i retrotxt_linux.rpm
retrotxt version
```

\* not implemented

## Build using Go

RetroTxt on Go requires [Go v1.13+](https://github.com/golang/go/wiki/MinimumRequirements).
Assuming [Go](https://golang.org/) and and the relevant build-tools are already installed.

```sh
git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
# recommended, run the package and import tests on your distribution/platform
go test ./...
go build -i
retrotxt version
```

The binary will be installed either at `$GOPATH/bin` or `$HOME/go/bin`

## Compile to other platforms using Go

Go supports a number of operating systems and platforms that can be built using any other supported platform.
Were you on a Linux system and needed to compile a 32-bit version of RetroTxt to target Windows 7.

```sh
# to see a list of supported platforms
go tool dist list

git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
# target 32-bit Windows
env GOOS=windows GOARCH=386 go build
# optional, compress the binary
upx retrotxt.exe
# test the binary
file retrotxt.exe
```

---

### Why Go?

- native [Unicode](https://golang.org/pkg/unicode/) and UTF-8/16 support
- [native legacy text encodings support](golang.org/x/text/encoding/charmap)
- creates a standalone binary with no dependencies
- [wide OS and CPU support](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)
- simple, compact standard library and fast compiling
- it is a language I know 😉
