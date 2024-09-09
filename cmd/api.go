package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/api"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Meli URL Shortener API",
	Long:  `Start the Meli URL Shortener API`,
	Run: func(cmd *cobra.Command, args []string) {
		app := config.GetApp()
		app.Ctx = cmd.Context()
		dispatcher := app.Api.GetDispatcher()
		defer dispatcher.Close()
		router := app.Api.Router(app.Name, app.Environment)
		server := app.Api.Server(http.Options{
			Port: app.Api.Port,
		})
		api.Start(app.Ctx, server, router)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
