package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Meli URL Shortener API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		api.Start(cmd.Context(), &api.Options{
			Port: config.GetIntEnv("API_PORT", "8080"),
		})
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
