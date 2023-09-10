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