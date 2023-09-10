# Retrotxt

### _[Retrotxt](https://github.com/bengarrett/retrotxt) for the terminal_.

Read legacy codepage and ANSI encoded textfiles in a modern Unicode terminal.

## Downloads

There are [numerous downloads](https://github.com/bengarrett/retrotxtgo/releases/latest/) available for
[Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Windows_Intel.zip),
[macOS](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_macOS_all.tar.gz),
[Linux](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_Linux_Intel.tar.gz) and more.

Otherwise [these installation options are available](INSTALL.md).

## Quick Usage

#### Text art and files created without Unicode often fail to display on modern systems.

#### Use RetroTxt to print legacy text on modern terminals.

```sh
retrotxt view ascii-logo.txt

██████╗ ███████╗████████╗██████╗  ██████╗ ████████╗██╗  ██╗████████╗
██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝╚██╗██╔╝╚══██╔══╝
██████╔╝█████╗     ██║   ██████╔╝██║   ██║   ██║    ╚███╔╝    ██║
██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║   ██║    ██╔██╗    ██║
██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝   ██║   ██╔╝ ██╗   ██║
╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝   ╚═╝
```

#### Or save it to a Unicode file and use it in other apps.

```sh
retrotxt view ascii-logo.txt > ascii-logo-utf8.txt
```

![Windows Notepad viewing ascii-logo-utf8.txt](img/ascii-logo-utf8.txt.png)

Otherwise, using the common shell programs, the legacy text is often malformed and even unreadable.

```sh
type ascii-logo.txt # or, cat ascii-logo.txt

�����ۻ ������ۻ�������ۻ�����ۻ  �����ۻ �������ۻ�ۻ  �ۻ�������ۻ
������ۻ������ͼ�������ͼ������ۻ�������ۻ�������ͼ��ۻ��ɼ�������ͼ
������ɼ����ۻ     �ۺ   ������ɼ�ۺ   �ۺ   �ۺ    ����ɼ    �ۺ
������ۻ����ͼ     �ۺ   ������ۻ�ۺ   �ۺ   �ۺ    ����ۻ    �ۺ
�ۺ  �ۺ������ۻ   �ۺ   �ۺ  �ۺ�������ɼ   �ۺ   ��ɼ �ۻ   �ۺ
�ͼ  �ͼ������ͼ   �ͼ   �ͼ  �ͼ �����ͼ    �ͼ   �ͼ  �ͼ   �ͼ
```

---

## Features

- Print legacy codepage encoded texts in a modern terminal.
- Print or export the legacy details of the textfiles.
- Print or export the SAUCE metadata of a file.
- Transform legacy encoded texts and text art into UTF-8 documents for use on the web or with modern systems.
- Lookup and print codepage character tables for dozens of encodings.
- IO redirection with piping support.

---

### Known code pages and text encodings

```
retrotxt list codepages

─────────────────────────────────────────────────────────────────────────
 Formal name              | Named value   | Numeric value  | Alias value   |
 * IBM Code Page 037      | cp037         | 37             | ibm037        |
 IBM Code Page 437        | cp437         | 437            | msdos         |
 IBM Code Page 850        | cp850         | 850            | latinI        |
 IBM Code Page 852        | cp852         | 852            | latinII       |
 IBM Code Page 855        | cp855         | 855            | ibm855        |
 Windows Code Page 858    | cp858         | 858            | ibm00858      |
 IBM Code Page 860        | cp860         | 860            | ibm860        |
 IBM Code Page 862        | cp862         | 862            |               |
 IBM Code Page 863        | cp863         | 863            | ibm863        |
 IBM Code Page 865        | cp865         | 865            | ibm865        |
 IBM Code Page 866        | ibm866        | 866            |               |
 * IBM Code Page 1047     | cp1047        | 1047           | ibm1047       |
 * IBM Code Page 1140     | cp1140        | 1140           | ibm01140      |
 ISO 8859-1               | iso-8859-1    | 1              | latin1        |
 ISO 8859-2               | iso-8859-2    | 2              | latin2        |
 ISO 8859-3               | iso-8859-3    | 3              | latin3        |
 ISO 8859-4               | iso-8859-4    | 4              | latin4        |
 ISO 8859-5               | iso-8859-5    | 5              | cyrillic      |
 ISO 8859-6               | iso-8859-6    | 6              | arabic        |
 ISO-8859-6E              | iso-8859-6-e  |                | iso88596e     |
 ISO-8859-6I              | iso-8859-6-i  |                | iso88596i     |
 ISO 8859-7               | iso-8859-7    | 7              | greek         |
 ISO 8859-8               | iso-8859-8    | 8              | hebrew        |
 ISO-8859-8E              | iso-8859-8-e  |                | iso88598e     |
 ISO-8859-8I              | iso-8859-8-i  |                | iso88598i     |
 ISO 8859-9               | iso-8859-9    | 9              | latin5        |
 ISO 8859-10              | iso-8859-10   | 10             | latin6        |
 ISO 8895-11              | iso-8895-11   | 11             | iso889511     |
 ISO 8859-13              | iso-8859-13   | 13             | iso885913     |
 ISO 8859-14              | iso-8859-14   | 14             | iso885914     |
 ISO 8859-15              | iso-8859-15   | 15             | iso885915     |
 ISO 8859-16              | iso-8859-16   | 16             | iso885916     |
 KOI8-R                   | koi8-r        |                | koi8r         |
 KOI8-U                   | koi8-u        |                | koi8u         |
 Macintosh                | macintosh     |                | mac           |
 Windows 874              | cp874         | 874            | windows-874   |
 Windows 1250             | cp1250        | 1250           | windows-1250  |
 Windows 1251             | cp1251        | 1251           | windows-1251  |
 Windows 1252             | cp1252        | 1252           | windows-1252  |
 Windows 1253             | cp1253        | 1253           | windows-1253  |
 Windows 1254             | cp1254        | 1254           | windows-1254  |
 Windows 1255             | cp1255        | 1255           | windows-1255  |
 Windows 1256             | cp1256        | 1256           | windows-1256  |
 Windows 1257             | cp1257        | 1257           | windows-1257  |
 Windows 1258             | cp1258        | 1258           | windows-1258  |
 Shift JIS                | shift_jis     |                | shiftjis      |
 UTF-8                    | utf-8         |                | utf8          |
 † UTF-16BE (Use BOM)     | utf-16        |                | utf16         |
 † UTF-16BE (Ignore BOM)  | utf-16be      |                | utf16be       |
 † UTF-16LE (Ignore BOM)  | utf-16le      |                | utf16le       |
 † UTF-32BE (Use BOM)     | utf-32        |                | utf32         |
 † UTF-32BE (Ignore BOM)  | utf-32be      |                | utf32be       |
 † UTF-32LE (Ignore BOM)  | utf-32le      |                | utf32le       |
 ⁑ ASA X3.4 1963          | ascii-63      |                |               |
 ⁑ ASA X3.4 1965          | ascii-65      |                |               |
 ⁑ ANSI X3.4 1967/77/86   | ascii-67      |                |               |

 * A EBCDIC encoding in use on IBM mainframes but is not ASCII compatible.
 † UTF-16/32 encodings are NOT usable with the list table command.
 ⁑ ANSI X3.4 encodings are only usable with the list table command.
   You can use the list table ascii command to list all three X3.4 tables.
```

### Even More Uses

#### Print legacy codepage tables in the terminal.

```
retrotxt list table cp437 latin1

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

#### Print file information and embedded SAUCE metadata.

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

```sh
retrotxt info retrotxt.ans --format=json
```

```json
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