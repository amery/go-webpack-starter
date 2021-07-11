package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "example",
	Short: "just a go-webpack-starter skeleton",
}

func main() {
	// include pid on the logs
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
