package cmd

import (
	"github.com/isaquesb/meli-url-shortener/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Meli URL Shortener Worker",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		worker.Consume(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
