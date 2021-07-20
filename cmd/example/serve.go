package main

import (
	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/flags/cobra"
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
	cobra.NewMapper(serveCmd.Flags()).
		VarP(&cfg.Development, "dev", 'd', "Don't hashify static files").
		VarP(&cfg.Server.Port, "port", 'p', "HTTP port to listen").
		VarP(&cfg.Server.PIDFile, "pid", 'f', "Path to PID file").
		VarP(&cfg.Server.GracefulTimeout, "graceful", 't', "Maximum duration to wait for in-flight requests")

	rootCmd.AddCommand(serveCmd)
}
