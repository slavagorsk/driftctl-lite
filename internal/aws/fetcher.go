package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceFetcher fetches live AWS resource attributes.
type ResourceFetcher struct {
	s3Client *s3.Client
}

// NewResourceFetcher creates a ResourceFetcher using the default AWS config.
func NewResourceFetcher(ctx context.Context) (*ResourceFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &ResourceFetcher{
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}

// FetchS3Bucket returns selected attributes for the given S3 bucket name.
func (f *ResourceFetcher) FetchS3Bucket(ctx context.Context, name string) (map[string]string, error) {
	out, err := f.s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching bucket %q: %w", name, err)
	}
	attrs := map[string]string{
		"bucket": name,
		"region": string(out.LocationConstraint),
	}
	return attrs, nil
}
