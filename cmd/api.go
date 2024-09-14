package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/api"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "URL Shortener API",
	Long:  `Start the URL Shortener API`,
	Run: func(cmd *cobra.Command, args []string) {
		app := config.NewApp()
		app.Ctx = cmd.Context()
		dispatcher := app.Api.GetDispatcher()
		defer dispatcher.Close()
		instrumentation := app.Instrumentation()
		router := app.Api.Router(instrumentation)
		server := app.Api.Server(http.Options{
			Port: app.Api.Port,
		})
		api.Start(app.Ctx, server, router, instrumentation)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
