### Contributions are always welcome.

## Code

_Retrotxt is written entirely in native Go as a platform-agnostic command-line only tool._

* Retrotxt uses [Go modules](https://blog.golang.org/using-go-modules).

* The target Go version is in [`go.mod`](https://github.com/bengarrett/retrotxtgo/blob/master/go.mod).

* Please keep your code to pure-Go, and this includes any 3rd party libraries or dependencies. Avoid including anything that relies on non-Go libraries or [CGO](https://golang.org/cmd/cgo/).

* Where possible follow idiomatic Go practises including formatting, variable naming conventions, etc.

* Keep your code platform agnostic.

* Do not use standard libraries that are platform-specific such as the _golang.org/x/sys/unix_ or _golang.org/x/sys/windows_ libraries.

* Never hardcode directory paths `/usr/local/bin/`, instead use `filepaths.Join("..", "foo")` and relative locations.

* Keep library dependencies to a minimum.
If you only need one or two small functions from a 3rd party library, rewrite those in your code and credit the source with links in the code comments.
**But do not wholesale cut and paste code**.

* The command-line interface uses both [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper), which is in [`/lib/cmd`](/lib/cmd).

Please make sure your code validates with [golangci-lint](https://golangci-lint.run/).
A [configuration file](https://github.com/bengarrett/retrotxtgo/blob/master/.golangci.yml) is in the root of `/retrotxtgo`.

```bash
cd retrotxtgo

# lint the entire application
golangci-lint run
# if nothing returns your code passes
```

Follow the [Go v1.13 static error syntax](https://blog.golang.org/go1.13-errors).

```go
var ErrNoName = errors.New("name cannot be empty")

func read(name string) error {
    if name == "" {
        return fmt.Errorf("read: %w", ErrNoName)
    }
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("read %q: %w", name, err)
	}
    defer f.Close()
    ...
}
```
