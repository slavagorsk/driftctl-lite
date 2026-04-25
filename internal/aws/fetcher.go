package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceFetcher dispatches fetch calls by resource type.
type ResourceFetcher struct {
	s3Client  S3Client
	rdsClient RDSClient
}

// NewResourceFetcher creates a ResourceFetcher using the given AWS region.
func NewResourceFetcher(region string) (*ResourceFetcher, error) {
	if region == "" {
		return nil, fmt.Errorf("fetcher: region must not be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("fetcher: load aws config: %w", err)
	}

	return &ResourceFetcher{
		s3Client:  s3.NewFromConfig(cfg),
		rdsClient: rds.NewFromConfig(cfg),
	}, nil
}

// Fetch retrieves live attributes for a resource identified by type and ID.
func (f *ResourceFetcher) Fetch(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	switch resourceType {
	case "aws_s3_bucket":
		return FetchS3Bucket(ctx, f.s3Client, resourceID)
	case "aws_db_instance":
		return FetchRDSInstance(ctx, f.rdsClient, resourceID)
	default:
		return nil, fmt.Errorf("fetcher: unsupported resource type %q", resourceType)
	}
}
