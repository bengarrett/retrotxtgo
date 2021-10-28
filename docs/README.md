# RetroTxtGo

### [RetroTxt](https://github.com/bengarrett/retrotxt) for the terminal.

###### version α, a work-in-progress with an incomplete feature set.

[Developer notes](DEV.md), [Dependencies project](https://github.com/bengarrett/retrotxtgo/projects/2), [TO-DO project](https://github.com/bengarrett/retrotxtgo/projects/1)

## Downloads

There are [downloads](https://github.com/bengarrett/retrotxtgo/releases/latest/) available for
[Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Windows_Intel.zip),
[macOS](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macOS_all.tar.gz),
[Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Linux_Intel.tar.gz),
[FreeBSD](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_FreeBSD_Intel.tar.gz) and the
[Raspberry Pi](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Linux_arm32_.tar.gz).

Otherwise [these package manager installations are available](#installations).

## Quick Usage

#### Text art and files created without Unicode often fail to display on modern systems.

```sh
type ascii-logo.txt # or, cat ascii-logo.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

#### Use RetroTxtGo to print legacy text on modern terminals.

```sh
retrotxt view ascii-logo.txt

██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```

#### Or save it to a Unicode file and view it in contemporary software.

```sh
retrotxt view ascii-logo.txt > ascii-logo-utf8.txt
```

![Windows Notepad viewing ascii-logo-utf8.txt](img/ascii-logo-utf8.txt.png)

#### Also turn the legacy text into a static website with accurate fonts and colors.

```sh
retrotxt create ascii-logo.txt --layout=compact
```

```html
<!DOCTYPE html>

<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>RetroTxtGo</title>
    <link rel="stylesheet" href="styles.css" />
    <link rel="stylesheet" href="font.css" />
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

#### Easily test and serve it over a HTTP server.

```sh
retrotxt create ascii-logo.txt -p0 --title="ASCII logo"

Web server is available at http://localhost:8086
Press Ctrl+C to stop
```

![Chrome browser viewing ascii-logo.txt as a website](img/ascii-logo.png)

#### Save the files as a ready to use webpage.

```sh
retrotxt create ascii-logo.txt --save

saving to styles.css
saving to index.html
saving to ibm-vga8.woff2
saving to font.css
...
```

#### Save the files to a ZIP archive.

```sh
retrotxt create ascii-logo.txt --compress

...
created zip file: retrotxt.zip 22.1 kB
```

#### Or inline all the website assets into a single HTML file for easier sharing.

```sh
retrotxt create ascii-logo.txt --layout=inline --font-embed --title="ASCII logo"
```

```html
<!DOCTYPE html>

<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>ASCII logo</title>
    <style type="text/css">
      body{background-color:#000;display:flex;flex-direction:column;text-shadow:none;color:
      var(--vgagray)}main{display:inline-block;letter-spacing:0;min-width:640px;order:2...}
      @font-face{font-family:vga;src:url(data:application/font-woff2;charset=utf-8;base64,
      d09GMgABAAAAADkcAA0AAAABE3AAADjAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP0ZGVE0cGiYGYACNThEI
      CoPBGILVJwuMHgABNgIkA4wkBCAFigEHsTobd9Z3xN13CsndqlJoJsgwEiFsHIxmBPMTbox8UzbTaPb///+fl
      mwcgRWuBPDebxsDSNVg3JCNeeFKtWkB70aFWWXONQ...)}
    </style>
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

## Features

- Convert ASCII text to HTML.
- [ ] Convert ANSI text to HTML.
- [ ] Convert BBS text to HTML.
- List or export (JSON, text, XML) meta details of a text file.
- List or export SAUCE metadata of a file.
- Transform the encoding of a text file. CP437 -> UTF8, UTF8 -> ISO8859-1 ...
- View any legacy encoded text file in a UTF8 terminal by converting on the fly.
- Extensive customizations through command flags or a configuration file with a setup.
- ANSI compatibility tests, output 16, 88, 256, high and true-color tables.
- Multi-platform support including Windows, macOS, Linux, Raspberry Pi and FreeBSD.
- IO redirection with piping support.

---

## Install

There are [download](https://github.com/bengarrett/retrotxtgo/releases/latest/) releases. Otherwise these package manager installations are available.
<a id="installations"></a>
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


#### [Building RetroTxtGo for other systems](BUILD.md)

---
### Even More Uses

#### Print legacy codepage tables in the terminal.

```
retrotext list table cp437 latin1

 ―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
         IBM Code Page 437 (DOS, OEM-US) - Extended ASCII
     0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
 0 |   | ☺ | ☻ | ♥ | ♦ | ♣ | ♠ | • | ◘ | ○ | ◙ | ♂ | ♀ | ♪ | ♫ | ☼ |
 1 | ► | ◄ | ↕ | ‼ | ¶ | § | ▬ | ↨ | ↑ | ↓ | → | ← | ∟ | ↔ | ▲ | ▼ |
 2 |   | ! | " | # | $ | % | & | ' | ( | ) | * | + | , | - | . | / |
 3 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | : | ; | < | = | > | ? |
 4 | @ | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O |
 5 | P | Q | R | S | T | U | V | W | X | Y | Z | [ | \ | ] | ^ | _ |
 6 | ` | a | b | c | d | e | f | g | h | i | j | k | l | m | n | o |
 7 | p | q | r | s | t | u | v | w | x | y | z | { | | | } | ~ | ⌂ |
 8 | Ç | ü | é | â | ä | à | å | ç | ê | ë | è | ï | î | ì | Ä | Å |
 9 | É | æ | Æ | ô | ö | ò | û | ù |   | Ö | Ü | ¢ | £ | ¥ | ₧ | ƒ |
 A | á | í | ó | ú | ñ | Ñ | ª | º | ¿ | ⌐ | ¬ | ½ | ¼ | ¡ | « | » |
 B | ░ | ▒ | ▓ | │ | ┤ | ╡ | ╢ | ╖ | ╕ | ╣ | ║ | ╗ | ╝ | ╜ | ╛ | ┐ |
 C | └ | ┴ | ┬ | ├ | ─ | ┼ | ╞ | ╟ | ╚ | ╔ | ╩ | ╦ | ╠ | ═ | ╬ | ╧ |
 D | ╨ | ╤ | ╥ | ╙ | ╘ | ╒ | ╓ | ╫ | ╪ | ┘ | ┌ | █ | ▄ | ▌ | ▐ | ▀ |
 E | α | ß | Γ | π | Σ | σ | µ | τ | Φ | Θ | Ω | δ | ∞ | φ | ε | ∩ |
 F | ≡ | ± | ≥ | ≤ | ⌠ | ⌡ | ÷ | ≈ | ° | ∙ | · | √ | ⁿ | ² | ■ |   |

 ―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
          ISO 8859-1 (Western European) - Extended ASCII
     0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
 0 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |
 1 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |
 2 |   | ! | " | # | $ | % | & | ' | ( | ) | * | + | , | - | . | / |
 3 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | : | ; | < | = | > | ? |
 4 | @ | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O |
 5 | P | Q | R | S | T | U | V | W | X | Y | Z | [ | \ | ] | ^ | _ |
 6 | ` | a | b | c | d | e | f | g | h | i | j | k | l | m | n | o |
 7 | p | q | r | s | t | u | v | w | x | y | z | { | | | } | ~ |   |
 8 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |
 9 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |
 A |   | ¡ | ¢ | £ | ¤ | ¥ | ¦ | § | ¨ | © | ª | « | ¬ | ­  | ® | ¯ |
 B | ° | ± | ² | ³ | ´ | µ | ¶ | · | ¸ | ¹ | º | » | ¼ | ½ | ¾ | ¿ |
 C | À | Á | Â | Ã | Ä | Å | Æ | Ç | È | É | Ê | Ë | Ì | Í | Î | Ï |
 D | Ð | Ñ | Ò | Ó | Ô | Õ | Ö | × | Ø | Ù | Ú | Û | Ü | Ý | Þ | ß |
 E | à | á | â | ã | ä | å | æ | ç | è | é | ê | ë | ì | í | î | ï |
 F | ð | ñ | ò | ó | ô | õ | ö | ÷ | ø | ù | ú | û | ü | ý | þ | ÿ |
```

#### Print file information and embedded SAUCE metadata.d

```
retrotxt info retrotxt.ans

────────────────────────────────────────────────────────────────────────────────
                                File information
 slug             retrotxt-ans
 filename         retrotxt.ans
 filetype         Text document with ANSI controls
 UTF-8            ✗
 line break       CRLF (Windows, DOS)
 characters       8,074
 ANSI controls    892
 words            59
 size             8.1 kB
 lines            23
 width            8,065
 modified         15 Aug 2021 23:33
 media mime type  application/octet-stream
 SHA256 checksum  ca1b69fa5ed2c01837b66f03402569f84c43fd308d8399abc85737e2abef2c1f
 CRC64 ECMA       1e8495be4c0edf25
 CRC32            7aeb63ec
 MD5              5adb64b98a10a87ba9bd02435112b049
 ───────────────────────────────
 title            RetroTxt logo
 author           Zeus II
 group            Blocktronics, FUEL
 date             1 Jul 2020
 original size    7.8 kB
 file type        ANSI color text
 data type        text or character stream
 description      ANSI text file with coloring codes and cursor positioning.
 character width  80
 number of lines  31
 comment          The app that lets you view works of ANSI art,  ASCII and NFO
 text in terminal or as HTML. LGPL and available for Windows, Linux, Unix & macOS.
```

```
retrotxt info retrotxt.ans --format=json

{
    "filename": "retrotxt.ans",
    "utf8": false,
    "lineBreak": {
        "string": "CRLF",
        "escape": "\r\n",
        "decimals": [
            13,
            10
        ]
    },
    ...
```


### Why create RetroTxt using Go?

- Native [Unicode](https://golang.org/pkg/unicode/) support.
- [A large selection of native legacy text encodings](golang.org/x/text/encoding/charmap).
- Builds a standalone binary with no dependencies.
- [Wide operating system and CPU architecture support](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).
- Uses a simple, compact standard library with extremely fast compiling.
- The standard library has helpful and safe web templating such as HTML, JSON, XML.
- It is a language I know. 😉
