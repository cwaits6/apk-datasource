package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cwaits6/apk-datasource/pkg/metrics"
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
	serveMetrics    bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Renovate-compatible JSON over HTTP with periodic refresh",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := log.Logger.WithContext(context.Background())
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		var (
			m              *metrics.Metrics
			metricsHandler http.Handler
		)

		if serveMetrics {
			var err error
			m, metricsHandler, err = metrics.Setup()
			if err != nil {
				return fmt.Errorf("setting up metrics: %w", err)
			}
			log.Info().Msg("metrics enabled on /metrics")
		} else {
			m = metrics.Noop()
		}

		srv := server.New(serveIndexURLs, servePort, refreshInterval, serveSourceURL, serveHomepage, m, metricsHandler)
		return srv.Run(ctx)
	},
}

func init() {
	serveCmd.Flags().StringSliceVar(&serveIndexURLs, "index-url", defaultIndexURLs, "APKINDEX.tar.gz URL(s) to fetch (repeatable)")
	serveCmd.Flags().IntVar(&servePort, "port", 3000, "HTTP port to listen on")
	serveCmd.Flags().DurationVar(&refreshInterval, "refresh-interval", 4*time.Hour, "Interval between index refreshes")
	serveCmd.Flags().StringVar(&serveSourceURL, "source-url", "", "Override source URL for all packages")
	serveCmd.Flags().StringVar(&serveHomepage, "homepage", "", "Override homepage for all packages")
	serveCmd.Flags().BoolVar(&serveMetrics, "metrics", true, "Enable Prometheus metrics on /metrics")
	rootCmd.AddCommand(serveCmd)
}
