package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/cwaits6/apk-datasource/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	serveIndexURLs  []string
	servePort       int
	refreshInterval time.Duration
	serveSourceURL  string
	serveHomepage   string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Renovate-compatible JSON over HTTP with periodic refresh",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := log.Logger.WithContext(context.Background())
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		srv := server.New(serveIndexURLs, servePort, refreshInterval, serveSourceURL, serveHomepage)
		return srv.Run(ctx)
	},
}

func init() {
	serveCmd.Flags().StringSliceVar(&serveIndexURLs, "index-url", nil, "APKINDEX.tar.gz URL(s) to fetch (repeatable)")
	serveCmd.Flags().IntVar(&servePort, "port", 3000, "HTTP port to listen on")
	serveCmd.Flags().DurationVar(&refreshInterval, "refresh-interval", 4*time.Hour, "Interval between index refreshes")
	serveCmd.Flags().StringVar(&serveSourceURL, "source-url", "", "Override source URL for all packages")
	serveCmd.Flags().StringVar(&serveHomepage, "homepage", "", "Override homepage for all packages")
	_ = serveCmd.MarkFlagRequired("index-url")
	rootCmd.AddCommand(serveCmd)
}
