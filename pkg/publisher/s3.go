package publisher

import (
	"context"
	"errors"

	"github.com/cwaits6/apk-datasource/pkg/generator"
)

// S3Publisher is a placeholder for future S3 publishing support.
type S3Publisher struct{}

// Publish returns an error indicating S3 publishing is not yet implemented.
func (s *S3Publisher) Publish(_ context.Context, _, _ string, _ *generator.RenovatePackage) error {
	return errors.New("S3 publisher is not yet implemented")
}
