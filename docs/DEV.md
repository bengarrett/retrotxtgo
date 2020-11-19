
### From here onwards the following text are developer notes that will eventually be removed

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

### Future CLI commands

- [ ] add optional argument for destination (dir or file) that overrides the dir configuration.
- - `retrotxt create somefile.txt . # would create index.html in the current directory`
- - `retrotxt create somefile.txt somefile.html # create sometfile.htm`
- [ ] add a flag to --export `create` command results to a tar or zip archive.
- [ ] option for generated HTML naming convention, either use index.html ... index_1.html, index_2.html etc.
      or filename.html, another-file-1.html, etc. When generating multiple HTML files, an index.html proof-sheet
      should be created with hyperlinks to all the other files. Maybe list their file/sauce details and a screenshot.
- [ ] both the `create/view` commands should support walking both directories and file archives.

---

### TODOs - changes to the existing code

- [ ] config shell should have a `--append/source/or` flag to save shell auto-completion?
- [ ] scan for supported but current shell configuration.
- [ ] scan for unique color codes like 24-bit, xterm and aix colors.
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

[Chocolatey](https://chocolatey.org/docs/createpackages)

Snap [snapcraft](https://snapcraft.io/first-snap#go), flathub is not for terminal apps.
