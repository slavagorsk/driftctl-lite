package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/snyk/driftctl-lite/internal/drift"
)

// isS3Type returns true if the resource type maps to an S3 resource.
func isS3Type(resourceType string) bool {
	return resourceType == "aws_s3_bucket"
}

// s3ClientFromConfig builds a real AWS S3 client from the shared config.
func s3ClientFromConfig(ctx context.Context, region string) (s3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("s3: load aws config: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}

// fetchS3Resource dispatches to the appropriate S3 fetch function.
func fetchS3Resource(ctx context.Context, region, resourceType, id string) (*drift.Resource, error) {
	client, err := s3ClientFromConfig(ctx, region)
	if err != nil {
		return nil, err
	}

	switch resourceType {
	case "aws_s3_bucket":
		return FetchS3Bucket(ctx, client, id)
	default:
		return nil, fmt.Errorf("s3: unsupported resource type %q", resourceType)
	}
}
