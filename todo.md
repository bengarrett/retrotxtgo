# to-dos

```go
    println(createCmd.CommandPath()) // <-- retrotxtgo create!
    println(createCmd.Name())        // <- create
    println(createCmd.UseLine())     // retrotxt create FILE [flags]
```

Shrink binaries
https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/

---

[ ] - Add (C) and licences + credits

[ ] - check for and handle cobra specific errors with flags etc.

- execute cobra flag needs an argument: 'f' in -f # go run . version -f
- execute cobra flag needs an argument: --format # version --format
- go run . info --format=color > no feedback
- create -p # does nothing, expect -n / --name

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
