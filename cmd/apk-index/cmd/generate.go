package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cwaits6/apk-datasource/pkg/fetcher"
	"github.com/cwaits6/apk-datasource/pkg/generator"
	"github.com/cwaits6/apk-datasource/pkg/publisher"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	indexURLs []string
	outputDir string
	sourceURL string
	homepage  string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Fetch APK indexes and generate Renovate-compatible JSON files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(indexURLs) == 0 {
			return fmt.Errorf("at least one --index-url is required")
		}

		ctx := log.Logger.WithContext(context.Background())
		start := time.Now()

		// Fetch.
		log.Info().Strs("urls", indexURLs).Msg("fetching indexes")
		sources, err := fetcher.FetchAll(ctx, indexURLs)
		if err != nil {
			return fmt.Errorf("fetching indexes: %w", err)
		}
		if len(sources) == 0 {
			return fmt.Errorf("no indexes fetched successfully")
		}

		// Generate Renovate JSON.
		data := generator.Generate(sources, sourceURL, homepage)

		// Publish.
		pub := publisher.NewFilesystemPublisher(outputDir)
		totalPackages := 0
		hasErrors := false

		for arch, pkgs := range data {
			for name, pkg := range pkgs {
				if err := pub.Publish(ctx, arch, name, pkg); err != nil {
					log.Error().Err(err).Str("package", name).Str("arch", arch).Msg("failed to publish")
					hasErrors = true
					continue
				}
				totalPackages++
			}
		}

		elapsed := time.Since(start)
		log.Info().
			Int("totalPackages", totalPackages).
			Int("architectures", len(data)).
			Dur("elapsed", elapsed).
			Msg("generation complete")

		if hasErrors {
			os.Exit(2)
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringSliceVar(&indexURLs, "index-url", nil, "APKINDEX.tar.gz URL(s) to fetch (repeatable)")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "./output", "Output directory for generated JSON files")
	generateCmd.Flags().StringVar(&sourceURL, "source-url", "", "Override source URL for all packages")
	generateCmd.Flags().StringVar(&homepage, "homepage", "", "Override homepage for all packages")
	generateCmd.MarkFlagRequired("index-url")
	rootCmd.AddCommand(generateCmd)
}
