package main

import (
	"net/http"

	"github.com/spf13/cobra"

	"go.sancus.dev/config/flags"

	"github.com/amery/go-webpack-starter/web"
)

type Router struct {
	web.Router
	http.Handler
}

// Command
var serveCmd = &cobra.Command{
	Use: "serve",
	Short: "serves example service",
	PreRun: func(cmd *cobra.Command, args []string) {
		flags.GetMapper(cmd.Flags()).Parse()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// compile templates
		html, err := web.CompileHtml(!cfg.Development)
		if err != nil {
			return err
		}

		// prepare server
		r := web.Router{
			HashifyAssets: !cfg.Development,
		}

		h := &Router{
			Router:  r,
			Handler: r.Handler(html),
		}

		return cfg.Server.ListenAndServe(h)
	},
}

// Flags
func init() {
	flags.NewCobraMapper(serveCmd.Flags()).
		BoolVarP(&cfg.Development, "dev", "d", false, "Don't hashify static files").
		UintVar16P(&cfg.Server.Port, "port", "p", 0, "HTTP port to listen").
		StringVarP(&cfg.Server.PIDFile, "pid", "f", "", "Path to PID file").
		DurationVarP(&cfg.Server.GracefulTimeout, "graceful", "t", 0, "Maximum duration to wait for in-flight requests")

	rootCmd.AddCommand(serveCmd)
}
