![Build](https://github.com/bengarrett/retrotxtgo/workflows/Go/badge.svg) ![Tests](https://github.com/bengarrett/retrotxtgo/workflows/Go%20tests/badge.svg)

# RetroTxtGo

### [RetroTxt](https://github.com/bengarrett/retrotxt) for the terminal.

###### version α, a work-in-progress with an incomplete feature set.

## About

#### Text art and files created without Unicode often fail to display on modern systems.

```sh
cat samples/ascii-logos.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

#### Use RetroTxtGo to print legacy text on modern terminals.

```sh
retrotxt view ascii-logos.txt

██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```

#### Or save it to a Unicode file that can be viewed with contemporary software.

```sh
retrotxt view ascii-logos.txt > ascii-logos-utf8.txt
```

#### Also turn the legacy text into a static website with accurate fonts and colors.

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
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝</pre>
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
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝</pre>
    </main>
  </body>
</html>
```

---

# Features

- [x] Convert ASCII text to HTML.
- [ ] Convert ANSI text to HTML.
- [ ] Convert BBS text to HTML.
- [x] List or export (JSON, text, XML) meta details of a text file.
- [x] List or export SAUCE metadata of a file.
- [x] Transform the encoding of a text file. CP437 -> UTF8, UTF8 -> ISO8859-1 ...
- [x] View any legacy encoded text file in a UTF8 terminal by converting on the fly.
- [x] Extensive customizations through command flags or a configuration file with a setup.
- [ ] ANSI compatibility tests, output 16, 88, 256, high and true-color tables.
- [x] Multi-platform support including Windows, macOS, Linux, Raspberry Pi and FreeBSD.
- [x] IO redirection with piping support.

---

## Install

There are [downloads](https://github.com/bengarrett/retrotxtgo/releases/latest/) available for
[Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Windows_Intel.zip),
[macOS Intel](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macOS_Intel.tar.gz),
[macOS M chips](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macOS_M-series.tar.gz),
[Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Linux_Intel.tar.gz),
[FreeBSD](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_FreeBSD_Intel.tar.gz) and the
[Raspberry Pi](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Linux_arm32_.tar.gz).

Otherwise these package manager installations are available.

#### Windows [Scoop](https://scoop.sh/)

```ps
scoop bucket add retrotxt https://github.com/bengarrett/retrotxtgo.git
scoop install bengarrett/retrotxt
retrotxt -v
```

#### macOS [Homebrew](https://brew.sh/)

```sh
brew cask install bengarrett/tap/retrotxt
retrotxt -v
```
#### DEB (Debian package)

```sh
# Intel
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.deb
dpkg -i retrotxt.deb
retrotxt -v
```

```sh
# Raspberry Pi & ARM
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_raspberry_pi.deb
dpkg -i retrotxt_raspberry_pi.deb
retrotxt -v
```


#### RPM (Redhat package)

```sh
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.rpm
rpm -i retrotxt.rpm
retrotxt -v
```

#### APK (Alpine package)
```sh
wget https://github.com/bengarrett/myip/releases/latest/download/retrotxt.apk
apk add retrotxt.apk
retrotxt -v
```

### Building RetroTxtGo for other systems

[Go](https://golang.org/doc/install) supports [dozens of architectures and operating systems](https://golang.org/doc/install/source#environment).


#### Working on Windows, PowerShell

```powershell
# clone and access repo
git clone https://github.com/bengarrett/retrotxtgo.git; cd retrotxtgo

# build retrotxt for the host system
go test ./...
go build -o retrotxt.exe
.\retrotxt.exe

# to see a list of supported platforms
go tool dist list

# to build retrotxt for windows/386, 32-bit Windows 7 and later
$Env:GOOS="windows"; $Env:GOARCH="386"; go build -o .\retrotxt.exe
```

#### Working on other shells and platforms
```bash
# clone and access repo
git clone https://github.com/bengarrett/retrotxtgo.git && cd retrotxtgo

# build retrotxt for the host system
go test ./...
go build -o retrotxt
./retrotxt -v

# to see a list of supported platforms
go tool dist list

# build retrotxt for linux/386, 32-bit Linux
env GOOS=linux GOARCH=386 && go build -o retrotxt && ./retrotxt -v
```

---

### Why create RetroTxt using Go?

- Native [Unicode](https://golang.org/pkg/unicode/), UTF 8/16/32 support.
- [A large selection of native legacy text encodings](golang.org/x/text/encoding/charmap).
- Builds a standalone binary with no dependencies.
- [Wide Operating system and CPU architecture support](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).
- Uses a simple, compact standard library with extremely fast compiling.
- The standard library has helpful and safe web templating such as HTML, JSON, XML.
- It is a language I know. 😉
