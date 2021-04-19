package create

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/prompt"
	"retrotxt.com/retrotxt/lib/str"
)

var ErrPort = errors.New("tried and failed to serve using these ports")

// Port checks if the TCP port is available on the local machine.
func Port(port uint) bool {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// Serve data over the internal HTTP server.
func (args *Args) Serve(b *[]byte) {
	args.override()
	port := args.Port
	if port == 0 || !prompt.PortValid(port) {
		port = uint(viper.GetInt("serve"))
	}
	max := 10
	for tries := 0; tries <= max; tries++ {
		switch {
		case !Port(port):
			port++
			continue
		case tries >= max:
			logs.Fatal("http ports", fmt.Sprintf("%d-%d", args.Port, port), ErrPort)
			continue
		default:
			args.Port = port
		}
		break
	}
	if err := args.createDir(b); err != nil {
		logs.Fatal("serve create directory", "", err)
	}
	if err := args.serveDir(); err != nil {
		logs.Fatal("serve directory", "", err)
	}
}

// Override the user flag values which are not yet implemented.
func (args *Args) override() {
	const embed = false
	var s = []string{}
	if args.layout != Standard {
		s = append(s, fmt.Sprintf("%s HTML layout", Standard))
		args.layout = Standard
	}
	if args.FontEmbed {
		s = append(s, "not embedding the font")
		args.FontEmbed = embed
	}
	l := len(s)
	if l == 0 {
		return
	}
	if l == 1 {
		fmt.Printf("Using the %s\n", s[0])
		return
	}
	fmt.Printf("Using %s\n", strings.Join(s, " and "))
}

// CreateDir creates a temporary save directory destination.
func (args *Args) createDir(b *[]byte) (err error) {
	args.Save.AsFiles, args.Save.OW = true, true
	args.Save.Destination, err = ioutil.TempDir(os.TempDir(), "*-serve")
	if err != nil {
		return fmt.Errorf("failed to make a temporary serve directory: %w", err)
	}
	args.Create(b)
	return nil
}

// ServeDir hosts the named souce directory on the internal HTTP server.
func (args *Args) serveDir() (err error) {
	http.Handle("/", http.FileServer(http.Dir(args.Save.Destination)))
	const timeout, localHost = 15, "127.0.0.1"
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%v", localHost, args.Port),
		WriteTimeout: timeout * time.Second,
		ReadTimeout:  timeout * time.Second,
	}
	fmt.Printf("\nWeb server is available at %s\n",
		str.Cp(fmt.Sprintf("http:/%v", srv.Addr)))
	fmt.Println(str.Cinf("Press Ctrl+C to stop\n"))
	args.watch()
	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("tcp listen and serve failed: %w", err)
	}
	return nil
}

// Watch intercepts Ctrl-C keypress combinations and exits to the operating system.
func (args *Args) watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nServer shutdown and directory removal of: %s\n", args.Save.Destination)
		tmp, err := path.Match(fmt.Sprintf("%s%s*",
			string(os.PathSeparator), os.TempDir()), args.Save.Destination)
		if err != nil {
			logs.Fatal("path match pattern failed", "", err)
		}
		if tmp {
			if err := os.RemoveAll(args.Save.Destination); err != nil {
				logs.Fatal("could not clean the temporary directory: %q: %s", args.Save.Destination, err)
			}
		}
		os.Exit(0)
	}()
}
