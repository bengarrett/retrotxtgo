# to-dos

```go
    println(createCmd.CommandPath()) // <-- retrotxtgo create!
    println(createCmd.Name())        // <- create
    println(cmd.CalledAs())          // create
    println(createCmd.UseLine())     // retrotxt create FILE [flags]
```

---

Golang built-ins

`slice = append(slice, more...)` `slice = append(slice, elm1, elm2)`
`copy(dest, srcSlice) = int` # number of elms copied

`cap(array) = int` `cap(slice) #max length`

`len(array, string, slice) int` # length or bytes

---

[ ] - Writing Web Applications https://golang.org/doc/articles/wiki/
A simple golang web server with basic logging, tracing, health check, graceful shut-down and zero dependencies
https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7
https://flaviocopes.com/golang-enable-cors/
https://developpaper.com/golang-http-connection-reuse-method/
https://sysadmins.co.za/golang-building-a-basic-web-server-in-go/

[ ] - Add a build subcommand for raw text data outputs

- --version, --build-date, --build-commit
- --sha256 (self-hash)

[ ] - handle Windows colour support

[ ] - create a VIEW command to convert cp437 into UTF8 and display on screen.

[ ] - flesh out the _Long_ commands

[ ] - config command should support the global --config flag

[ ] - test config editor on Windows

[ ] - config shell should have a --append/source/or ? to save auto-completion

- do a scan for the shell file
- save a new file in the save directory or elsewhere when appropriate .retrotxtgo-completions.sh and import to user's shell
- when successful tell user to `source xxx` or reload the terminal/command prompt

[ ] - add a flag to export `create` to an tar/zip archive

---

## SAUCE debug?

`fmt.Sprintf()`

#### Integer

%b base 2 (binary)
%c unicode code point
%d base 10
%o base 8
%0 base 8 with 0o prefix
%x base 16 lower-case %X upper-case
%U unicode format

\+ always print a sign (%+q)
\- pad with spaces to the right rather than the left
\# alternate formats:
add leading 0b for binary (%#b)
add leading 0 for octal (%#o)
add 0x or 0X for hex (%#x)
suppress 0x for %p (%#p)
always print a decimal point for floats (%#f)...
0 pad with leading zeros rather than spaces (%0d)...

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

## Possible create options

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

#### Target package managers

Shrink binaries
https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/

Minimum requirements
https://github.com/golang/go/wiki/MinimumRequirements

Ubuntu/Mint APT .DEB (http://packaging.ubuntu.com/html/packaging-new-software.html)
Maybe just include instructions for local install?

Fedora, CentOS/ RHEL
https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/index.html
https://rpm-packaging-guide.github.io/

Windows
scoop.sh
https://github.com/lukesampson/scoop/wiki/App-Manifests

Chocolatey
https://chocolatey.org/docs/createpackages

Homebrew (macOS/Linux) [requires compiling]
https://docs.brew.sh/Formula-Cookbook

Homebrew confusing terminology
https://docs.brew.sh/Formula-Cookbook#homebrew-terminology

- formula - Package definition
- keg - installation prefix of a formula
- cellar - location of kegs
- tap - git repo of formulae or commands
- bottle - prebuilt keg

Homebrew Cask [allows binaries]
https://github.com/Homebrew/homebrew-cask/blob/master/doc/cask_language_reference/readme.md
https://github.com/Homebrew/homebrew-cask/blob/master/doc/cask_language_reference/all_stanzas.md

Raspbian ?
