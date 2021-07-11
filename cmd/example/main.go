package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pborman/getopt/v2"

	"github.com/amery/go-webpack-starter/web"
)

var config = NewConfig()

func init() {
	getopt.FlagLong(&config.Development, "dev", 'd', "Don't hashify static files")
	getopt.FlagLong(&config.Server.Port, "port", 'p', "HTTP port to listen")
	getopt.FlagLong(&config.Server.PIDFile, "pid", 'f', "Path to PID file")
	getopt.FlagLong(&config.Server.GracefulTimeout, "graceful", 't', "Maximum duration to wait for in-flight requests")

	getopt.Parse()

	// TODO: validate config
}

type Router struct {
	web.Router
	http.Handler
}

func main() {
	// include pid on the logs
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	// compile templates
	html, err := web.CompileHtml(!config.Development)
	if err != nil {
		log.Fatal(err)
	}

	// prepare server
	r := web.Router{
		HashifyAssets: !config.Development,
	}

	h := &Router{
		Router:  r,
		Handler: r.Handler(html),
	}

	if err := config.Server.ListenAndServe(h); err != nil {
		log.Fatal(err)
	}
}
