package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Meli URL Shortener Worker",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := config.NewApp()
		app.Ctx = cmd.Context()
		dispatcher := app.Worker.GetDispatcher()
		defer dispatcher.Close()
		consumer := app.Worker.GetConsumer()
		worker.Consume(cmd.Context(), consumer)
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
