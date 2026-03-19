package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logLevel  string
	logFormat string
)

var rootCmd = &cobra.Command{
	Use:   "apk-datasource",
	Short: "Generate Renovate-compatible datasource files from APK indexes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			level = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(level)

		if logFormat == "text" {
			log.Logger = zerolog.New(zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			}).With().Timestamp().Logger()
		} else {
			log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "Log format (text, json)")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
