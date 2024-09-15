package cmd

import (
	"github.com/isaquesb/url-shortener/config"
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "URL Shortener Worker",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := config.NewApp()
		app.Ctx = cmd.Context()
		dispatcher := app.Worker.Dispatcher.Get()
		defer dispatcher.Close()
		consumer := app.Worker.Consumer.Get()
		if _, ok := dispatcher.(output.Listen); ok {
			go dispatcher.(output.Listen).Listen(app.Ctx)
		}
		worker.Consume(cmd.Context(), consumer)
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
