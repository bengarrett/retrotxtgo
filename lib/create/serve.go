package create

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Serve data over an internal HTTP server.
func (args *Args) Serve(b *[]byte) {
	port := uint(args.Port)
	if port == 0 || !prompt.PortValid(port) {
		port = uint(viper.GetInt("serve"))
	}
	args.Port = port
	if err := args.createDir(b); err != nil {
		logs.ChkErr(logs.Err{Issue: "create problem", Arg: "HTTP", Msg: err})
	}
	if err := args.serveDir(); err != nil {
		logs.ChkErr(logs.Err{Issue: "server problem", Arg: "HTTP", Msg: err})
	}
}

func (args *Args) createDir(b *[]byte) (err error) {
	// TODO: deep-copy channel support?
	args.SaveToFile, args.OW = true, true
	args.Dest, err = ioutil.TempDir(os.TempDir(), "*-serve")
	if err != nil {
		return err
	}
	args.Create(b)
	return nil
}

// serveDir creates and serves b over an internal HTTP server.
func (args Args) serveDir() (err error) {
	http.Handle("/", http.FileServer(http.Dir(args.Dest)))
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
		return err
	}
	return nil
}

// watch intercepts Ctrl-C exit key combination.
func (args Args) watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nServer shutdown and clean up: %s\n", args.Dest)
		tmp, err := path.Match(fmt.Sprintf("%s%s*",
			string(os.PathSeparator), os.TempDir()), fmt.Sprintf("%s", args.Dest))
		if err != nil {
			log.Fatal(err)
		}
		if tmp {
			os.RemoveAll(args.Dest)
		}
		os.Exit(0)
	}()
}
