package create

import (
	"fmt"
	"net/http"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Serve data over an internal HTTP server.
func (args *Args) Serve(data *[]byte) {
	port := uint(args.Port)
	if port == 0 || !prompt.PortValid(port) {
		port = uint(viper.GetInt("serve"))
	}
	args.Port = port
	if err := args.serveFile(data); err != nil {
		logs.ChkErr(logs.Err{Issue: "server problem", Arg: "HTTP", Msg: err})
	}
}

// serveFile creates and serves the html template over an internal HTTP server.
func (args Args) serveFile(data *[]byte) error {
	t, err := args.newTemplate()
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d, err := args.pagedata(data)
		if err != nil {
			logs.ChkErr(logs.Err{Issue: "servefile encoding", Arg: "http", Msg: err})
		}
		if err = t.Execute(w, d); err != nil {
			logs.ChkErr(logs.Err{Issue: "servefile", Arg: "http", Msg: err})
		}
	})
	fs := http.FileServer(http.Dir("static/"))
	port := args.Port
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("Web server is available at %s\n",
		str.Cp(fmt.Sprintf("http://localhost:%v", port)))
	println(str.Cinf("Press Ctrl+C to stop"))
	if err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		return err
	}
	return nil
}
