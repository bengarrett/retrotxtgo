package create

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
func (args *Args) Serve(b *[]byte) (issue, arg string, err error) {
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
			return "http ports", fmt.Sprintf("%d-%d", args.Port, port), ErrPort
		default:
			args.Port = port
		}
		break
	}
	if err := args.createDir(b); err != nil {
		return "serve createDir", "", err
	}

	args.serveDir()

	return "", "", nil
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
		if !args.Test {
			fmt.Printf("Using the %s\n", s[0])
		}
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
	err = args.Create(b)
	if err != nil {
		return fmt.Errorf("failed to create files in the temporary directory: %w", err)
	}
	return nil
}

// ServeDir hosts the named souce directory on the internal HTTP server.
func (args *Args) serveDir() {
	const timeout = 5
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", args.Port),
	}
	http.Handle("/", http.FileServer(http.Dir(args.Save.Destination)))

	var ctx context.Context
	var cancel context.CancelFunc
	if args.Test {
		go func() {
			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("tcp listen and serve failed: %s\n", err)
			}
		}()
		ctx, cancel = context.WithCancel(context.Background())
		// cancel the server straight away
		cancel()
	} else {
		// listen for Ctrl+C keyboard interrupt
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			fmt.Printf("\nWeb server is available at %s\n",
				str.Cp(fmt.Sprintf("http://localhost%v", srv.Addr)))
			fmt.Println(str.Cinf("Press Ctrl+C to stop\n"))

			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("tcp listen and serve failed: %s\n", err)
			}
		}()

		<-done
		ctx, cancel = context.WithTimeout(context.Background(), timeout*time.Second)
	}

	defer func() {
		cancel()
		args.cleanup()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown failed: %s\n", err)
	} else if args.Test {
		fmt.Print("Server example was successful")
	}
}

var ErrCleanPath = errors.New("cleanup temporary path match failed")

// Cleanup the temporary files and directories.
func (args *Args) cleanup() {
	if !args.Test {
		fmt.Printf("\n\nServer shutdown and directory removal of: %s\n", args.Save.Destination)
	}
	tmp, err := path.Match(fmt.Sprintf("%s%s*",
		string(os.PathSeparator), os.TempDir()), args.Save.Destination)
	if err != nil {
		logs.ProblemFatal(ErrCleanPath, err)
	}
	if tmp {
		if err := os.RemoveAll(args.Save.Destination); err != nil {
			logs.MarkProblemFatal(args.Save.Destination, logs.ErrTmpClean, err)
		}
	}
}
