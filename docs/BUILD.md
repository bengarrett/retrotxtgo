
# RetroTxtGo

### Building RetroTxtGo for other systems

[Go](https://golang.org/doc/install) supports [dozens of architectures and operating systems](https://golang.org/doc/install/source#environment).

Note: The version string returned by `retrotxt -v` will always be `x0.0.0` for user builds.

#### Working on Windows, PowerShell

```powershell
# clone and access repo
git clone https://github.com/bengarrett/retrotxtgo.git; cd retrotxtgo

# build retrotxt for the host system
go test ./...
go build -o retrotxt.exe
.\retrotxt.exe -v

# to see a list of supported platforms
go tool dist list

# to build retrotxt for windows/386, 32-bit Windows 7 and later
$Env:GOOS="windows"; $Env:GOARCH="386"; go build -o retrotxt.exe; .\retrotxt.exe -v

# be sure to remove or revert these two temporary environment vars
$Env:GOOS=""; $Env:GOARCH="";
```

#### Working with other shells and platforms
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
