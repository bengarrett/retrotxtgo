# RetroTxtGo

RetroTxt for the command line 

### *α* work-in-progress, feature incomplete & is not in a usable state

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
retrotxt samples/ascii-logos.txt # creates an index.html file in the home directory
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
      <pre>██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```


---

# Features

- [ ] Convert ASCII text to HTML.
- [ ] Convert ANSI text to HTML.
- [ ] Convert BBS text to HTML.
- [X] List or export (json, text, xml) meta details of a text file.
- [ ] List or export SAUCE metadata of a file.
- [ ] Transform the encoding of a text file. CP437 -> UTF8, UTF8 -> ISO8859-1 ...
- [ ] View any legacy encoded text file in a UTF8 terminal by converting on the fly.
- [X] Extensive customisations through command flags or a configuration file with a setup.
- [ ] ANSI compatibility tests, output 16, 88, 256, high and true-color tables.
- [X] Multiplatform support including Windows, macOS, Linux, Raspbian and FreeBSD.

---

## Install

RetroTxtGo requires [Go v1.13+](https://github.com/golang/go/wiki/MinimumRequirements)

Assuming [Go](https://golang.org/) and and the relevant build-tools are already installed.

```sh
cd ~
git clone https://github.com/bengarrett/retrotxtgo.git
# go install -o retrotxt.exe . # on Windows
go install -o retrotxt .
retrotxt --version
```

The binary will be installed either at `$GOPATH/bin` or `$HOME/go/bin`

## Compile

Assuming Go and and the relevant build-tools are already installed.

```sh
cd ~
git clone https://github.com/bengarrett/retrotxtgo.git
cd retrotxtgo
# optional but recommended, run package and import tests on your distribution/platform
go test ./...
# build with --race to also detect concurrent race conditions on major 64-bit operating systems
go build -o retrotxt --race .
retrotxt version
```

#### Race detection

While in α and to troubleshoot, 64-bit versions of linux (arm/amd64), freebsd, darwin, windows should build with the `--race` flag.
Though this will increase the binary size.

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

### Cobra hints

Print references
```go
    println(xxxCmd.CommandPath()) // retrotxtgo xxx
    println(xxxCmd.Name())        // create
    println(cmd.CalledAs())       // create
    println(xxxCmd.UseLine())     // retrotxt xxx FILE [flags]
```

### Go built-ins

byte arrays

`slice = append(slice, more...)` `slice = append(slice, elm1, elm2)`
`copy(dest, srcSlice) = int` # number of elms copied

`cap(array) = int` `cap(slice) #max length, as oppose to current length`

`len(array, string, slice) int` # length or bytes

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

-----

### TODOs - changes to the existing code

- [ ] **Remove all string references to `retrotxtgo`.**
- [ ] Directory settings, change `.` shortcut to always use the current working directory. 
Include text mentioning this when using `set`.
- [ ] When fetching github release data using the `version` command. 
Use HTTP e-tags cache and save the values to reduce the bandwidth usage.
- [ ] config command should support the global --config flag.
- [ ] config shell should have a `--append/source/or` flag to save shell auto-completion?
- [ ] scan for supported but current shell configuration.
- [ ] when using `create` detect any out of range or unsafe unicode encoding and assume cp437.

---

### Future distribution package managers

### Windows

**[scoop](https://scoop.sh/)** [App Manifests](https://github.com/lukesampson/scoop/wiki/App-Manifests)

**[Chocolatey](https://chocolatey.org/docs/createpackages)**

### macOS

**[homebrew](https://brew.sh/)** [casks allow bin submissions](https://github.com/Homebrew/homebrew-cask/blob/master/doc/cask_language_reference/readme.md)

### Linux (auto update)

**[snapcraft](https://snapcraft.io/)** Snaps (https://snapcraft.io/first-snap#go)

~~**[flathub](https://github.com/flathub/flathub/wiki/App-Submission)**~~ Not intended for CLI apps

### Debian/Ubuntu/Raspbian (manual updates)

**[deb]** `.deb` installed with `dpkg -i`

### Fedora/RHEL/CentOS (manual updates)

**[rpm]** Fedora/RHEL/CentOS installed with `rpm -i`

#### Other managers that require sponsorship

Ubuntu/Mint APT .DEB (http://packaging.ubuntu.com/html/packaging-new-software.html)
See 4.5 and 4.6

Fedora, CentOS/ RHEL
https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/index.html
https://rpm-packaging-guide.github.io/
