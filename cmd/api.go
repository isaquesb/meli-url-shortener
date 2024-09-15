package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/api"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "URL Shortener API",
	Long:  `Start the URL Shortener API`,
	Run: func(cmd *cobra.Command, args []string) {
		app := config.NewApp()
		app.Ctx = cmd.Context()
		dispatcher := app.Api.Dispatcher.Get()
		defer dispatcher.Close()
		instrumentation := app.Instrumentation()
		router := app.Api.Router(instrumentation)
		server := app.Api.Server(http.Options{
			Port: app.Api.Port,
		})
		if _, ok := dispatcher.(output.Listen); ok {
			go dispatcher.(output.Listen).Listen(app.Ctx)
		}

		apiServer := &api.Api{
			Ctx:    app.Ctx,
			Server: server,
			Router: router,
			Instr:  instrumentation,
		}

		apiServer.Start()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
