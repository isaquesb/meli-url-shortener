package cmd

import (
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "shortener",
	Short: "Meli URL Shortener",
	Long:  "Meli URL Shortener",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger.Setup(logger.NewLogger(config.GetEnv("LOG_LEVEL", "info")))

	config.LoadEnv()
}
