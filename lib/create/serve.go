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
	"syscall"
	"time"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/prompt"
	"retrotxt.com/retrotxt/lib/str"
)

// Port checks the TCP port is available on the local machine.
func Port(port uint) bool {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}

// Serve data over an internal HTTP server.
func (args *Args) Serve(b *[]byte) {
	port := uint(args.Port)
	if port == 0 || !prompt.PortValid(port) {
		port = uint(viper.GetInt("serve"))
	}
	max := 10
	for tries := 0; tries <= max; tries++ {
		if !Port(port) {
			port++
		} else if tries >= max {
			logs.Fatal("http ports", fmt.Sprintf("%d-%d", args.Port, port),
				errors.New("tried and failed to serve using these ports"))
		} else {
			args.Port = port
			break
		}
	}
	if err := args.createDir(b); err != nil {
		logs.Fatal("serve create directory", "", err)
	}
	if err := args.serveDir(); err != nil {
		logs.Fatal("serve directory", "", err)
	}
}

func (args *Args) createDir(b *[]byte) (err error) {
	args.SaveToFile, args.OW = true, true
	args.Destination, err = ioutil.TempDir(os.TempDir(), "*-serve")
	if err != nil {
		return fmt.Errorf("failed to make a temporary serve directory: %s", err)
	}
	args.Create(b)
	return nil
}

// serveDir creates and serves b over an internal HTTP server.
func (args Args) serveDir() (err error) {
	http.Handle("/", http.FileServer(http.Dir(args.Destination)))
	srv := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%v", args.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Web server is available at %s\n",
		str.Cp(fmt.Sprintf("http:/%v", srv.Addr)))
	fmt.Println(str.Cinf("Press Ctrl+C to stop"))
	args.watch()
	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("tcp listen and serve failed: %s", err)
	}
	return nil
}

// watch intercepts Ctrl-C exit key combination.
func (args Args) watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nServer shutdown and clean up: %s\n", args.Destination)
		tmp, err := path.Match(fmt.Sprintf("%s%s*",
			string(os.PathSeparator), os.TempDir()), fmt.Sprintf("%s", args.Destination))
		if err != nil {
			logs.Fatal("path match pattern failed", "", err)
		}
		if tmp {
			os.RemoveAll(args.Destination)
		}
		os.Exit(0)
	}()
}
