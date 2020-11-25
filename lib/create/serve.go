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

// ErrPort port failed.
var ErrPort = errors.New("tried and failed to serve using these ports")

// Port checks the TCP port is available on the local machine.
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

// Serve data over an internal HTTP server.
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

func (args *Args) override() {
	const embed = false
	var s = []string{}
	if args.layout != Standard {
		s = append(s, fmt.Sprintf("layout to %q", Standard))
		args.layout = Standard
	}
	if args.FontEmbed {
		s = append(s, fmt.Sprintf("font-embed to %v", false))
		args.FontEmbed = embed
	}
	l := len(s)
	if l == 0 {
		return
	}
	if l == 1 {
		fmt.Println("Replaced this argument:", s[0])
		return
	}
	fmt.Println("Replaced these arguments:", strings.Join(s, ", "))
}

func (args *Args) createDir(b *[]byte) (err error) {
	args.Output.SaveToFile, args.Output.OW = true, true
	args.Output.Destination, err = ioutil.TempDir(os.TempDir(), "*-serve")
	if err != nil {
		return fmt.Errorf("failed to make a temporary serve directory: %w", err)
	}
	args.Create(b)
	return nil
}

// serveDir creates and serves b over an internal HTTP server.
func (args *Args) serveDir() (err error) {
	http.Handle("/", http.FileServer(http.Dir(args.Output.Destination)))
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

// watch intercepts Ctrl-C exit key combination.
func (args *Args) watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nServer shutdown and directory removal of: %s\n", args.Output.Destination)
		tmp, err := path.Match(fmt.Sprintf("%s%s*",
			string(os.PathSeparator), os.TempDir()), args.Output.Destination)
		if err != nil {
			logs.Fatal("path match pattern failed", "", err)
		}
		if tmp {
			if err := os.RemoveAll(args.Output.Destination); err != nil {
				logs.Fatal("could not clean the temporary directory: %q: %s", args.Output.Destination, err)
			}
		}
		os.Exit(0)
	}()
}
