package main

import (
	"github.com/spf13/cobra"

	"go.sancus.dev/config/flags"
)

// Command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves example service",
	PreRun: func(cmd *cobra.Command, args []string) {
		flags.GetMapper(cmd.Flags()).Parse()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// prepare server
		r := &Router{
			HashifyAssets: !cfg.Development,
		}

		// compile templates
		if err := r.Compile(); err != nil {
			return err
		}

		return cfg.Server.ListenAndServe(r)
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
