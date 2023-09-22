// Copyright © 2023 Ben Garrett. All rights reserved.
// Use of this source code is governed by a GNU
// license that can be found in the LICENSE file.

/*
Retrotxt for the terminal.
Read legacy code page and ANSI encoded text files in a modern
Unicode terminal.

Text files and art created before the adoption of Unicode often
fail to display on modern systems.
Use RetroTxt to print legacy text on modern terminals.
Or save it to a Unicode file and use it in other apps.
Otherwise, when using many command prompt or terminal apps,
legacy text is often malformed and even unreadable.

# Features

  - Print legacy code page encoded texts in a modern terminal.
  - Print or export the details of the text files.
  - Print or export the SAUCE metadata of a file.
  - Transform legacy encoded texts and text art into UTF-8 documents for use on the web or with modern systems.
  - Lookup and print code page character tables for dozens of encodings.
  - Support for ISO, PC-DOS/Windows code pages plus IBM EBCDIC, Macintosh, and ShiftJIS.
  - Use io redirection with piping support.

Usage:

	retrotxt [command]

The commands are:

	lang        List the natural languages of legacy code pages
	list        List the legacy code pages that Retrotxt can convert to UTF-8
	table       Display one or more code page tables showing all the characters in use
	tables      Display the characters of every code page table in use
	info        Information on a text file
	view        Print a text file to the terminal using standard output
	example     List the included sample text files available for use with the info and view commands

# Examples

To display information about a text file:

	retrotxt info [filenames]

To display information about a file in JSON format:

	retrotxt info [filenames] --format json

To display a text file to the terminal:

	retrotxt view [filenames]

To display a text file to the terminal by supplying the source encoding:

	retrotxt view [filenames] --input iso-8859-1

To list the sample text files:

	retrotxt example

To list the supported code page encodings with names and aliases:

	retrotxt list

To list the supported code pages and characters as tables:

	retrotxt tables

To list the a code page and characters as a table:

	retrotxt table [code page names or aliases]

To list both the table of Code Page 437 and ISO 8859-1 using aliases:

	retrotxt table cp437 latin1

To list the target natural languages of the supported code pages:

	retrotxt lang
*/
package main
