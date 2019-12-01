# to-dos

[ ] - either create a flag CREATE --web-serve=true
or incorporate and import CREATE into SERVE.

> nested subcommands?
> CREATE SERVE textfiles/hi.txt --port=8080
> retrotxtgo server textfiles/hi.txt -p 8082

[ ] - built bash and zsh autocomplete (automatic?)

[ ] - global, local and cascading flags

[ ] - cobra add completion

```go
    println(createCmd.CommandPath()) // <-- retrotxtgo create!
    println(createCmd.Name())        // <- create
    println(createCmd.UseLine())     // retrotxt create FILE [flags]
```

HTTP server
https://github.com/julienschmidt/httprouter

Shrink binaries
https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
