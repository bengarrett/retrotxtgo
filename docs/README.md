# Retrotxt
[![Go Reference](https://pkg.go.dev/badge/github.com/bengarrett/retrotxtgo.svg)](https://pkg.go.dev/github.com/bengarrett/retrotxtgo) [![GoReleaser](https://github.com/bengarrett/retrotxtgo/actions/workflows/release.yml/badge.svg)](https://github.com/bengarrett/retrotxtgo/actions/workflows/release.yml)

### _[RetroTxt](https://github.com/bengarrett/retrotxt) for the terminal_

Read legacy code pages and ANSI-encoded text files in modern Unicode terminals.

## Downloads

[Numerous downloads are available](https://github.com/bengarrett/retrotxtgo/releases/latest/) for:
- [Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_windows.zip)
- [macOS](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_apple_silicon.gz)
- [Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.gz)

Windows users can use File Explorer to decompress it.

The Linux and macOS downloads are standalone terminal applications in a gzip compressed binary. 

```
# replace 'foo' with the remainder of the filename
$ gzip -d retrotxt_foo.gz

# after decompression, to confirm the download and version
$ retrotxt -v
```

macOS users will need to delete the 'quarantine' extended attribute that is applied to all 
program downloads that are not notarized by Apple for a fee.

```
$ xattr -d com.apple.quarantine retrotxt
```

#### Linux packages:

- [Debian/Ubuntu (.deb)](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt.deb)
- [Fedora (.rpm)](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt.rpm)
- [Arch Linux (.zst)](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt.pkg.tar.zst)
- [Alpine Linux (.apk)](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt.apk)

#### Homebrew

macOS and Linux users can install via Homebrew:

```bash
brew tap bengarrett/retrotxt https://github.com/bengarrett/retrotxtgo
brew install bengarrett/retrotxt/retrotxt
```

Update to the latest version with:

```bash
brew upgrade bengarrett/retrotxt/retrotxt
```

## Quick Usage

Text files and art created before Unicode was widely adopted often fail to display correctly on modern systems.

#### Use RetroTxt to display legacy text in modern terminals.

```sh
$ retrotxt view ascii-logo.txt

██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```

#### Or save it to a Unicode file and use it in other apps.

```sh
$ retrotxt view ascii-logo.txt > ascii-logo-utf8.txt
```

![Windows Notepad viewing ascii-logo-utf8.txt](img/ascii-logo-utf8.txt.png)

Otherwise, legacy text is often malformed and unreadable when using most terminal apps.

```sh
$ type ascii-logo.txt # or, cat ascii-logo.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

---

## In the real world

### Swedish text

```
$ cat 14670.txt

�mnet d�r, fortfor den andre, r�r till en del de antika gallerna; men
f�r det mesta sker h�ndelserna i andev�rlden.

--I andev�rlden? Ja just. Jag har f�rnummit, redan i Songes, att s�
skall vara. Vad tycker du, Eleonora? Frans, Aurora ... i andev�rlden,
mina barn!

Farv�l, farv�l tills vi r�kas!

End of the Project Gutenberg EBook of Det g�r an, by Carl Jonas Love Almqvist
```

```
$ retrotxt 14670.txt --input latin1

Ämnet där, fortfor den andre, rör till en del de antika gallerna; men
för det mesta sker händelserna i andevärlden.

--I andevärlden? Ja just. Jag har förnummit, redan i Songes, att så
skall vara. Vad tycker du, Eleonora? Frans, Aurora ... i andevärlden,
mina barn!

Farväl, farväl tills vi råkas!

End of the Project Gutenberg EBook of Det går an, by Carl Jonas Love Almqvist
```

### Japanese text

```sh
$ cat rshmn10j.txt
# no text is displayed as the text isn't ASCII compatible
```

```
$ retrotxt rshmn10j.txt --input shiftjis

　暫、死んだように倒れていた老婆が、屍骸の中から、その裸の体を起こしたのは、そ
れから間もなくの事である。老婆は、つぶやくような、うめくような声を立てながら、
まだ燃えている火の光をたよりに、梯子の口まで、這って行った。そうして、そこから、
短い白髪を倒（さかさま）にして、門の下を覗きこんだ。外には、唯、黒洞々（こくと
うとう）たる夜があるばかりである。
　下人は、既に、雨を冒して、京都の町へ強盗を働きに急いでいた。
```

### Cyrillic text

```
$ cat olavg10.txt

"������ � ��� ���������, ����� ���,
���� ������� ����, �������
�� ������� ���� � ��������� �����,
����� � � ������� ����������.
��� ����� ����� �� ���� �� ������
���-�������, ���-�������� - ����
�� ��������� �� ������� ���������,
����� �� �� ������� � ��� ���."

*** END OF THE PROJECT GUTENBERG EBOOK, OLAF VAN GELDERN ***
```

```
$ retrotxt olavg10.txt --input cp1251

"Честит е тоз избранник, чийто дух,
като ковчега Ноев, пренесе
от прежний свят в послешний онова,
което е в промени непроменно.
Той подир смърт от себе си оставя
най-чистото, най-хубавото - жица
от царството на сенките безплътни,
която ще го свързва с тоя мир."

*** END OF THE PROJECT GUTENBERG EBOOK, OLAF VAN GELDERN ***
```


---

## Features

- Display legacy code page encoded texts in modern terminals.
- Print or export detailed information about text files.
- Print or export the [SAUCE metadata](https://www.acid.org/info/sauce/sauce.htm) of files.
- Transform legacy encoded texts and text art into UTF-8 documents for use on the web or with modern systems.
- Look up code page and character tables for dozens of encodings and print the results.
- Support for ISO, PC-DOS/Windows code pages, IBM EBCDIC, Macintosh, and ShiftJIS encodings.
- Use I/O redirection with piping support.

---

### Known code pages and text encodings

```
$ retrotxt list

┌──────────────────────────────────────────────────────────────────────────────┐
│ Formal name                  Named value     Numeric value    Alias value    │
│ IBM Code Page 037            cp037           37               ibm037         │
│ IBM Code Page 437            cp437           437              msdos          │
│ IBM Code Page 850            cp850           850              latinI         │
│ IBM Code Page 852            cp852           852              latinII        │
│ IBM Code Page 855            cp855           855              ibm855         │
│ Windows Code Page 858        cp858           858              ibm00858       │
│ IBM Code Page 860            cp860           860              ibm860         │
│ IBM Code Page 862            cp862           862                             │
│ IBM Code Page 863            cp863           863              ibm863         │
│ IBM Code Page 865            cp865           865              ibm865         │
│ IBM Code Page 866            ibm866          866                             │
│ IBM Code Page 1047           cp1047          1047             ibm1047        │
│ IBM Code Page 1140           cp1140          1140             ibm01140       │
│ ISO 8859-1                   iso-8859-1      1                latin1         │
│ ISO 8859-2                   iso-8859-2      2                latin2         │
│ ISO 8859-3                   iso-8859-3      3                latin3         │
│ ISO 8859-4                   iso-8859-4      4                latin4         │
│ ISO 8859-5                   iso-8859-5      5                cyrillic       │
│ ISO 8859-6                   iso-8859-6      6                arabic         │
│ ISO-8859-6E                  iso-8859-6-e                     iso88596e      │
│ ISO-8859-6I                  iso-8859-6-i                     iso88596i      │
│ ISO 8859-7                   iso-8859-7      7                greek          │
│ ISO 8859-8                   iso-8859-8      8                hebrew         │
│ ISO-8859-8E                  iso-8859-8-e                     iso88598e      │
│ ISO-8859-8I                  iso-8859-8-i                     iso88598i      │
│ ISO 8859-9                   iso-8859-9      9                latin5         │
│ ISO 8859-10                  iso-8859-10     10               latin6         │
│ ISO-8859-11                  iso-8859-11     11               iso885911      │
│ ISO 8859-13                  iso-8859-13     13               iso885913      │
│ ISO 8859-14                  iso-8859-14     14               iso885914      │
│ ISO 8859-15                  iso-8859-15     15               iso885915      │
│ ISO 8859-16                  iso-8859-16     16               iso885916      │
│ KOI8-R                       koi8-r                           koi8r          │
│ KOI8-U                       koi8-u                           koi8u          │
│ Macintosh                    macintosh                        mac            │
│ Windows 874                  cp874           874              windows-874    │
│ Windows 1250                 cp1250          1250             windows-1250   │
│ Windows 1251                 cp1251          1251             windows-1251   │
│ Windows 1252                 cp1252          1252             windows-1252   │
│ Windows 1253                 cp1253          1253             windows-1253   │
│ Windows 1254                 cp1254          1254             windows-1254   │
│ Windows 1255                 cp1255          1255             windows-1255   │
│ Windows 1256                 cp1256          1256             windows-1256   │
│ Windows 1257                 cp1257          1257             windows-1257   │
│ Windows 1258                 cp1258          1258             windows-1258   │
│ Shift JIS                    shift_jis                        shiftjis       │
│ Big5                         big5                             big-5          │
│ UTF-8                        utf-8                            utf8           │
│ UTF-16BE (Use BOM)           utf-16                           utf16          │
│ UTF-16BE (Ignore BOM)        utf-16be                         utf16be        │
│ UTF-16LE (Ignore BOM)        utf-16le                         utf16le        │
│ UTF-32BE (Use BOM)           utf-32                           utf32          │
│ UTF-32BE (Ignore BOM)        utf-32be                         utf32be        │
│ UTF-32LE (Ignore BOM)        utf-32le                         utf32le        │
│ ASA X3.4 1963                ascii-63        1963                            │
│ ASA X3.4 1965                ascii-65        1965                            │
│ ANSI X3.4 1967/77/86         ascii-67        1967             ansi           │
└──────────────────────────────────────────────────────────────────────────────┘

 Yellow text indicates EBCDIC encodings found on some IBM mainframes.
 Darker text indicates encodings not usable with the table command.
 Purple text indicates encodings only usable with the table command.
 You can use the "table ascii" command to list all three X3.4 tables.

Named, numeric, or alias values are all valid code page arguments.
These values all match ISO 8859-1:
  retrotxt table iso-8859-1  # named
  retrotxt table 1           # numeric
  retrotxt table latin1      # alias

 Today, most applications and the web use Unicode UTF-8.
 As a subset, UTF-8 is backwards compatible with ANSI X3.4 (US-ASCII).
 IBM Code Page 437 is commonly used on MS-DOS and for ANSI art.
 ISO 8859-1 is found on historic Unix, Amiga, and the early Internet.
 Windows 1252 is found on Windows in the 1980s and 1990s.
 Macintosh is found on Mac OS 9 and earlier systems.
 EBCDIC is incompatible with ANSI X3.4, most computers, and the web.
```

### Even More Uses

#### Display legacy code page tables in the terminal.

```
$ retrotxt table cp437 latin1

┌─────────────────────────────────────────────────────────────────────┐
│ IBM Code Page 437 (DOS, OEM-US) - Extended ASCII                    │
│   . 0 . 1 . 2 . 3 . 4 . 5 . 6 . 7 . 8 . 9 . A . B . C . D . E . F . │
│ 0 | ␀ | ☺ | ☻ | ♥ | ♦ | ♣ | ♠ | • | ◘ | ○ | ◙ | ♂ | ♀ | ♪ | ♫ | ☼ | │
│ 1 | ► | ◄ | ↕ | ‼ | ¶ | § | ▬ | ↨ | ↑ | ↓ | → | ← | ∟ | ↔ | ▲ | ▼ | │
│ 2 |   | ! | " | # | $ | % | & | ' | ( | ) | * | + | , | - | . | / | │
│ 3 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | : | ; | < | = | > | ? | │
│ 4 | @ | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | │
│ 5 | P | Q | R | S | T | U | V | W | X | Y | Z | [ | \ | ] | ^ | _ | │
│ 6 | ` | a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | │
│ 7 | p | q | r | s | t | u | v | w | x | y | z | { | | | } | ~ | ⌂ | │
│ 8 | Ç | ü | é | â | ä | à | å | ç | ê | ë | è | ï | î | ì | Ä | Å | │
│ 9 | É | æ | Æ | ô | ö | ò | û | ù | ÿ | Ö | Ü | ¢ | £ | ¥ | ₧ | ƒ | │
│ A | á | í | ó | ú | ñ | Ñ | ª | º | ¿ | ⌐ | ¬ | ½ | ¼ | ¡ | « | » | │
│ B | ░ | ▒ | ▓ | │ | ┤ | ╡ | ╢ | ╖ | ╕ | ╣ | ║ | ╗ | ╝ | ╜ | ╛ | ┐ | │
│ C | └ | ┴ | ┬ | ├ | ─ | ┼ | ╞ | ╟ | ╚ | ╔ | ╩ | ╦ | ╠ | ═ | ╬ | ╧ | │
│ D | ╨ | ╤ | ╥ | ╙ | ╘ | ╒ | ╓ | ╫ | ╪ | ┘ | ┌ | █ | ▄ | ▌ | ▐ | ▀ | │
│ E | α | ß | Γ | π | Σ | σ | µ | τ | Φ | Θ | Ω | δ | ∞ | φ | ε | ∩ | │
│ F | ≡ | ± | ≥ | ≤ | ⌠ | ⌡ | ÷ | ≈ | ° | ∙ | · | √ | ⁿ | ² | ■ |   | │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│ ISO 8859-1 (Western European) - Extended ASCII                      │
│   . 0 . 1 . 2 . 3 . 4 . 5 . 6 . 7 . 8 . 9 . A . B . C . D . E . F . │
│ 0 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   | │
│ 1 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   | │
│ 2 |   | ! | " | # | $ | % | & | ' | ( | ) | * | + | , | - | . | / | │
│ 3 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | : | ; | < | = | > | ? | │
│ 4 | @ | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | │
│ 5 | P | Q | R | S | T | U | V | W | X | Y | Z | [ | \ | ] | ^ | _ | │
│ 6 | ` | a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | │
│ 7 | p | q | r | s | t | u | v | w | x | y | z | { | | | } | ~ |   | │
│ 8 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   | │
│ 9 |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   |   | │
│ A |   | ¡ | ¢ | £ | ¤ | ¥ | ¦ | § | ¨ | © | ª | « | ¬ | ­  | ® | ¯ | │
│ B | ° | ± | ² | ³ | ´ | µ | ¶ | · | ¸ | ¹ | º | » | ¼ | ½ | ¾ | ¿ | │
│ C | À | Á | Â | Ã | Ä | Å | Æ | Ç | È | É | Ê | Ë | Ì | Í | Î | Ï | │
│ D | Ð | Ñ | Ò | Ó | Ô | Õ | Ö | × | Ø | Ù | Ú | Û | Ü | Ý | Þ | ß | │
│ E | à | á | â | ã | ä | å | æ | ç | è | é | ê | ë | ì | í | î | ï | │
│ F | ð | ñ | ò | ó | ô | õ | ö | ÷ | ø | ù | ú | û | ü | ý | þ | ÿ | │
└─────────────────────────────────────────────────────────────────────┘
```

#### Display file information and embedded SAUCE metadata.

```
$ retrotxt info retrotxt.ans

┌──────────────────┐
│ File Information │
└──────────────────┘

retrotxt.ans
├── Basic Information
│   ├── slug: retrotxt-ans
│   ├── filename: retrotxt.ans
│   ├── filetype: Text document with ANSI controls
│   ├── Unicode: no
│   └── line break: CRLF (Windows, DOS)
├── Content Statistics
│   ├── characters: 8,074
│   ├── ANSI controls: 892
│   ├── words: 59
│   ├── size: 8.1 kB
│   ├── lines: 31
│   └── width: 8,101
├── File Metadata
│   ├── modified: 3 Oct 2023
│   └── media mime type: application/octet-stream
├── Checksums & Integrity
│   └── SHA256 checksum: ca1b69fa5ed2c01837b66f03402569f84c43fd308d8399abc85737e2abef2c1f
└── SAUCE Metadata
    ├── title: RetroTxt logo
    ├── author: Zeus II
    ├── group: Blocktronics, FUEL
    ├── date: 1 Jul 2020
    ├── original size: 7.8 kB
    ├── file type: ANSI color text
    ├── data type: text or character stream
    ├── description: ANSI text file with coloring codes and cursor positioning.
    ├── character width: 80
    ├── number of lines: 31
    └── interpretation: blink mode, invalid value
        └── Comments
            The app that lets you view works of ANSI art, ASCII and NFO text
             in terminal or as HTML. LGPL and available for Windows, Linux,
            Unix & macOS.
```


```sh
$ retrotxt info retrotxt.ans --format json
```

```json
{
    "filename": "retrotxt.ans",
    "unicode": "no",
    "lineBreak": {
        "string": "CRLF",
        "escape": "\r\n",
        "decimal": [
            13,
            10
        ]
    },
    "counts": {
        "characters": 8074,
        "ansiControls": 892,
        "words": 59
    },
    "size": {
        "bytes": 8101,
        "decimal": "8.1 kB",
        "binary": "7.9 KiB"
    },
    "lines": 31,
    "width": 8101,
    "modified": {
        "iso": "2023-09-10T06:25:49.00179453Z",
        "epoch": 1694327149
    },
    "checksums": {
        "crc32": "7aeb63ec",
        "crc64": "1e8495be4c0edf25",
        "md5": "5adb64b98a10a87ba9bd02435112b049",
        "sha256": "ca1b69fa5ed2c01837b66f03402569f84c43fd308d8399abc85737e2abef2c1f"
    },
    "mime": {
        "media": "application",
        "subMedia": "octet-stream",
        "comment": "Text document with ANSI controls"
    },
    "slug": "retrotxt-ansi",
    "sauce": {
        "id": "SAUCE",
        "version": "00",
        "title": "RetroTxt logo",
        "author": "Zeus II",
        "group": "Blocktronics, FUEL",
        "date": {
            "value": "20200701",
            "iso": "2020-07-01T00:00:00Z",
            "epoch": 1593561600
        },
        "filesize": {
            "bytes": 7775,
            "decimal": "7.8 kB",
            "binary": "7.6 KiB"
        },
        "dataType": {
            "type": 1,
            "name": "text or character stream"
        },
        "fileType": {
            "type": 1,
            "name": "ANSI color text"
        },
        "typeInfo": {
            "1": {
                "value": 80,
                "info": "character width"
            },
            "2": {
                "value": 31,
                "info": "number of lines"
            },
            "3": {
                "value": 0,
                "info": ""
            },
            "flags": {
                "decimal": 3,
                "binary": "00011",
                "nonBlinkMode": {
                    "flag": "0",
                    "interpretation": "blink mode"
                },
                "letterSpacing": {
                    "flag": "00",
                    "interpretation": "no preference"
                },
                "aspectRatio": {
                    "flag": "11",
                    "interpretation": "invalid value"
                }
            },
            "fontName": "IBM VGA"
        },
        "comments": {
            "id": "COMNT",
            "count": 3,
            "lines": [
                "The app that lets you view works of ANSI art, ASCII and NFO text in terminal or as HTML. LGPL and available for Windows, Linux, Unix \u0026 macOS.                                                   "
            ]
        }
    },
    "zipComment": "",
    "UTF8": false
}
```

#### Dump the hexadecimal bytes of a file.

```sh
$ retrotxt dump sauce-embed.txt
```

```
00000000  44 6f 63 75 6d 65 6e 74  20 65 6e 63 6f 64 69 6e  |Document encodin|
00000010  67 3a 20 57 69 6e 64 6f  77 73 2d 31 32 35 32 0d  |g: Windows-1252.|
00000020  0a 27 54 72 61 6e 73 63  6f 64 65 20 74 65 78 74  |.'Transcode text|
00000030  27 20 73 65 74 74 69 6e  67 3a 20 41 75 74 6f 6d  |' setting: Autom|
00000040  61 74 69 63 0d 0a 46 69  6c 65 20 6e 61 6d 65 3a  |atic..File name:|
00000050  20 73 61 75 63 65 2e 74  78 74 0d 0a 0d 0a 52 65  | sauce.txt....Re|
00000060  74 72 6f 54 78 74 20 53  41 55 43 45 20 74 65 73  |troTxt SAUCE tes|
00000070  74 0d 0a 0d 0a 54 68 65  20 74 6f 70 20 6f 66 20  |t....The top of |
00000080  74 68 65 20 70 61 67 65  20 73 68 6f 75 6c 64 20  |the page should |
00000090  64 69 73 70 6c 61 79 20  74 68 65 20 66 6f 6c 6c  |display the foll|
000000a0  6f 77 69 6e 67 20 74 65  78 74 0d 0a 0d 0a 27 53  |owing text....'S|
000000b0  61 75 63 65 20 74 69 74  6c 65 27 20 62 79 20 27  |auce title' by '|
000000c0  53 61 75 63 65 20 61 75  74 68 6f 72 27 20 6f 66  |Sauce author' of|
000000d0  20 53 61 75 63 65 20 67  72 6f 75 70 2c 20 64 61  | Sauce group, da|
000000e0  74 65 64 20 32 30 31 36  20 4e 6f 76 65 6d 62 65  |ted 2016 Novembe|
000000f0  72 20 32 36 0d 0a 41 6e  79 20 63 6f 6d 6d 65 6e  |r 26..Any commen|
00000100  74 73 20 67 6f 20 68 65  72 65 2e 0d 0a 0d 0a 43  |ts go here.....C|
00000110  72 61 73 20 73 69 74 20  61 6d 65 74 20 70 75 72  |ras sit amet pur|
00000120  75 73 20 75 72 6e 61 2e  20 50 68 61 73 65 6c 6c  |us urna. Phasell|
00000130  75 73 20 69 6e 20 64 61  70 69 62 75 73 20 65 78  |us in dapibus ex|
00000140  2e 20 50 72 6f 69 6e 20  70 72 65 74 69 75 6d 20  |. Proin pretium |
00000150  65 67 65 74 20 6c 65 6f  20 75 74 20 67 72 61 76  |eget leo ut grav|
00000160  69 64 61 2e 20 50 72 61  65 73 65 6e 74 20 65 67  |ida. Praesent eg|
00000170  65 73 74 61 73 20 75 72  6e 61 20 61 74 20 74 69  |estas urna at ti|
00000180  6e 63 69 64 75 6e 74 20  6d 6f 6c 6c 69 73 2e 20  |ncidunt mollis. |
00000190  56 69 76 61 6d 75 73 20  6e 65 63 20 75 72 6e 61  |Vivamus nec urna|
000001a0  20 6c 6f 72 65 6d 2e 20  56 65 73 74 69 62 75 6c  | lorem. Vestibul|
000001b0  75 6d 20 6d 6f 6c 65 73  74 69 65 20 61 63 63 75  |um molestie accu|
000001c0  6d 73 61 6e 20 6c 65 63  74 75 73 2c 20 69 6e 20  |msan lectus, in |
000001d0  65 67 65 73 74 61 73 20  6d 65 74 75 73 20 66 61  |egestas metus fa|
000001e0  63 69 6c 69 73 69 73 20  65 67 65 74 2e 20 4e 61  |cilisis eget. Na|
000001f0  6d 20 63 6f 6e 73 65 63  74 65 74 75 72 2c 20 6d  |m consectetur, m|
00000200  65 74 75 73 20 65 74 20  73 6f 64 61 6c 65 73 20  |etus et sodales |
00000210  61 6c 69 71 75 61 6d 2c  20 6d 69 20 64 75 69 20  |aliquam, mi dui |
00000220  64 61 70 69 62 75 73 20  6d 65 74 75 73 2c 20 6e  |dapibus metus, n|
00000230  6f 6e 20 63 75 72 73 75  73 20 6c 69 62 65 72 6f  |on cursus libero|
00000240  20 66 65 6c 69 73 20 61  63 20 6e 75 6e 63 2e 20  | felis ac nunc. |
00000250  4e 75 6c 6c 61 20 65 75  69 73 6d 6f 64 2c 20 74  |Nulla euismod, t|
00000260  75 72 70 69 73 20 73 65  64 20 6d 6f 6c 6c 69 73  |urpis sed mollis|
00000270  20 66 61 75 63 69 62 75  73 2c 20 6f 72 63 69 20  | faucibus, orci |
00000280  65 6c 69 74 20 64 61 70  69 62 75 73 20 6c 65 6f  |elit dapibus leo|
00000290  2c 20 76 69 74 61 65 20  70 6c 61 63 65 72 61 74  |, vitae placerat|
000002a0  20 64 69 61 6d 20 65 72  6f 73 20 73 65 64 20 76  | diam eros sed v|
000002b0  65 6c 69 74 2e 20 46 75  73 63 65 20 63 6f 6e 76  |elit. Fusce conv|
000002c0  61 6c 6c 69 73 2c 20 6c  6f 72 65 6d 20 75 74 20  |allis, lorem ut |
000002d0  76 75 6c 70 75 74 61 74  65 20 73 75 73 63 69 70  |vulputate suscip|
000002e0  69 74 2c 20 74 6f 72 74  6f 72 20 72 69 73 75 73  |it, tortor risus|
000002f0  20 72 68 6f 6e 63 75 73  20 6d 61 75 72 69 73 2c  | rhoncus mauris,|
00000300  20 61 20 6d 61 74 74 69  73 20 6d 65 74 75 73 20  | a mattis metus |
00000310  6d 61 67 6e 61 20 61 74  20 6c 6f 72 65 6d 2e 20  |magna at lorem. |
00000320  53 65 64 20 6d 6f 6c 65  73 74 69 65 20 76 65 6c  |Sed molestie vel|
00000330  69 74 20 69 70 73 75 6d  2c 20 69 6e 20 76 75 6c  |it ipsum, in vul|
00000340  70 75 74 61 74 65 20 6d  65 74 75 73 20 63 6f 6e  |putate metus con|
00000350  73 65 71 75 61 74 20 65  67 65 74 2e 20 46 75 73  |sequat eget. Fus|
00000360  63 65 20 71 75 69 73 20  64 75 69 20 6c 61 63 69  |ce quis dui laci|
00000370  6e 69 61 2c 20 6c 61 6f  72 65 65 74 20 6c 65 63  |nia, laoreet lec|
00000380  74 75 73 20 65 74 2c 20  6c 75 63 74 75 73 20 76  |tus et, luctus v|
00000390  65 6c 69 74 2e 20 50 65  6c 6c 65 6e 74 65 73 71  |elit. Pellentesq|
000003a0  75 65 20 75 74 20 6e 69  73 69 20 71 75 69 73 20  |ue ut nisi quis |
000003b0  6f 72 63 69 20 70 75 6c  76 69 6e 61 72 20 70 6c  |orci pulvinar pl|
000003c0  61 63 65 72 61 74 20 76  65 6c 20 61 63 20 6c 6f  |acerat vel ac lo|
000003d0  72 65 6d 2e 20 4d 61 65  63 65 6e 61 73 20 66 69  |rem. Maecenas fi|
000003e0  6e 69 62 75 73 20 66 65  72 6d 65 6e 74 75 6d 20  |nibus fermentum |
000003f0  65 72 61 74 2c 20 61 20  70 75 6c 76 69 6e 61 72  |erat, a pulvinar|
00000400  20 61 75 67 75 65 20 64  69 63 74 75 6d 20 6d 61  | augue dictum ma|
00000410  74 74 69 73 2e 20 41 65  6e 65 61 6e 20 76 75 6c  |ttis. Aenean vul|
00000420  70 75 74 61 74 65 20 63  6f 6e 73 65 63 74 65 74  |putate consectet|
00000430  75 72 20 76 65 6c 69 74  20 61 74 20 64 69 63 74  |ur velit at dict|
00000440  75 6d 2e 20 44 6f 6e 65  63 20 76 65 68 69 63 75  |um. Donec vehicu|
00000450  6c 61 20 61 6e 74 65 20  71 75 69 73 20 61 6e 74  |la ante quis ant|
00000460  65 20 76 65 6e 65 6e 61  74 69 73 2c 20 65 75 20  |e venenatis, eu |
00000470  75 6c 74 72 69 63 65 73  20 6c 65 63 74 75 73 20  |ultrices lectus |
00000480  65 67 65 73 74 61 73 2e  20 56 65 73 74 69 62 75  |egestas. Vestibu|
00000490  6c 75 6d 20 61 6e 74 65  20 69 70 73 75 6d 20 70  |lum ante ipsum p|
000004a0  72 69 6d 69 73 20 69 6e  20 66 61 75 63 69 62 75  |rimis in faucibu|
000004b0  73 20 6f 72 63 69 20 6c  75 63 74 75 73 20 65 74  |s orci luctus et|
000004c0  20 75 6c 74 72 69 63 65  73 20 70 6f 73 75 65 72  | ultrices posuer|
000004d0  65 20 63 75 62 69 6c 69  61 20 43 75 72 61 65 3b  |e cubilia Curae;|
000004e0  0d 0a 1a 43 4f 4d 4e 54  41 6e 79 20 63 6f 6d 6d  |...COMNTAny comm|
000004f0  65 6e 74 73 20 67 6f 20  68 65 72 65 2e 20 20 20  |ents go here.   |
00000500  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000510  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000520  20 20 20 20 20 20 20 20  53 41 55 43 45 30 30 53  |        SAUCE00S|
00000530  61 75 63 65 20 74 69 74  6c 65 20 20 20 20 20 20  |auce title      |
00000540  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000550  20 20 53 61 75 63 65 20  61 75 74 68 6f 72 20 20  |  Sauce author  |
00000560  20 20 20 20 20 20 53 61  75 63 65 20 67 72 6f 75  |      Sauce grou|
00000570  70 20 20 20 20 20 20 20  20 20 32 30 31 36 31 31  |p         201611|
00000580  32 36 81 0e 00 00 01 00  d1 03 09 00 00 00 00 00  |26..............|
00000590  01 13 49 42 4d 20 56 47  41 00 00 00 00 00 00 00  |..IBM VGA.......|
000005a0  00 00 00 00 00 00 00 00                           |........|
```

---

> [!TIP]
> Also available is [RetroTxt, the browser Extension](https://github.com/bengarrett/RetroTxt) that turns ANSI, ASCII, and NFO text into web documents.
